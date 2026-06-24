package abi

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	codePkg "github.com/tonkeeper/tongo/code"
	"github.com/tonkeeper/tongo/liteclient"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/utils"
)

type MethodInvocation struct {
	Result   any
	TypeHint string
}

type ContractDescription struct {
	// Interfaces is a list of interfaces implemented by a contract.
	ContractInterfaces []ContractInterface
	GetMethods         []MethodInvocation
	MethodsInspected   int
}

type libResolver interface {
	GetLibraries(ctx context.Context, libraryList []ton.Bits256) (map[ton.Bits256]*boc.Cell, error)
}

type contractInspector struct {
	knownMethods    []MethodDescription
	knownInterfaces []InterfaceDescription
	scanAllMethods  bool
	libResolver     libResolver
}

type InspectorOptions struct {
	additionalMethods []MethodDescription
	knownInterfaces   []InterfaceDescription
	scanAllMethods    bool
	libResolver       libResolver
}

type ContractInterface uint32

func (c ContractInterface) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *ContractInterface) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*c = ContractInterfaceFromString(s)
	return nil
}

func (c ContractInterface) Implements(other ContractInterface) bool {
	if c == other {
		return true
	}
	return c.recursiveImplements(other)
}

type InvokeFn func(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (string, any, error)

// MethodDescription describes a particular method and provides a function to execute it.
type MethodDescription struct {
	Name string
	// InvokeFn executes this method on a contract and returns parsed execution results.
	InvokeFn InvokeFn
}

type knownContractDescription struct {
	contractInterfaces []ContractInterface
	getMethods         []InvokeFn
}

type InterfaceDescription struct {
	Name    ContractInterface
	Results []string
}

type InspectorOption func(o *InspectorOptions)

func InspectWithAdditionalMethods(list []MethodDescription) InspectorOption {
	return func(o *InspectorOptions) {
		o.additionalMethods = list
	}
}

func InspectWithAdditionalInterfaces(list []InterfaceDescription) InspectorOption {
	return func(o *InspectorOptions) {
		o.knownInterfaces = list
	}
}
func InspectWithAllMethods() InspectorOption {
	return func(o *InspectorOptions) {
		o.scanAllMethods = true
	}
}

func InspectWithLibraryResolver(resolver libResolver) InspectorOption {
	return func(o *InspectorOptions) {
		o.libResolver = resolver
	}
}

func NewContractInspector(opts ...InspectorOption) *contractInspector {
	options := &InspectorOptions{}
	for _, o := range opts {
		o(options)
	}
	return &contractInspector{
		knownMethods:    append(append(methodInvocationOrder, tolkMethods...), options.additionalMethods...),
		knownInterfaces: append(append(contractInterfacesOrder, tolkInterfaceOrder...), options.knownInterfaces...),
		scanAllMethods:  options.scanAllMethods,
		libResolver:     options.libResolver,
	}
}

// MethodStatus describes the outcome of attempting a single get method during inspection.
type MethodStatus int

const (
	// MethodTransportError means the call to the executor failed (network / lite-server / emulator
	// infrastructure error). The result of such a method cannot be trusted. Use
	// liteclient.IsClientError on Err to tell a lite-server error apart from other failures.
	MethodTransportError MethodStatus = iota
	// MethodReverted means the method ran but returned a non-success VM exit code
	// (anything other than 0 or 1). Typically the method does not exist on the contract
	// or it intentionally threw. See ExitCode.
	MethodReverted
	// MethodDecodeFailed means the method ran successfully and returned a stack, but none of the
	// known decoders could parse that stack. This is the signal the old InspectContract swallowed:
	// the contract HAS the method but our ABI did not match its output. Stack holds the raw,
	// undecoded result so the caller can debug or re-decode it.
	MethodDecodeFailed
	// MethodDecoded means the method ran and its output was decoded. Result/TypeHint are populated.
	MethodDecoded
)

func (s MethodStatus) String() string {
	switch s {
	case MethodTransportError:
		return "transport_error"
	case MethodReverted:
		return "reverted"
	case MethodDecodeFailed:
		return "decode_failed"
	case MethodDecoded:
		return "decoded"
	default:
		return "unknown"
	}
}

// MethodResult is the full outcome of attempting one get method, with every failure mode surfaced.
type MethodResult struct {
	// Name is the get-method name. It is empty for methods of a known contract, which are
	// identified by MethodID instead.
	Name     string
	MethodID int
	Status   MethodStatus
	// ExitCode is the VM exit code. Meaningful when the method actually ran
	// (Status is MethodReverted, MethodDecodeFailed or MethodDecoded).
	ExitCode uint32
	// Stack is the raw stack returned by the VM. It is preserved on MethodDecodeFailed so the
	// caller has all the information needed to inspect an unrecognized result.
	Stack tlb.VmStack
	// TypeHint and Result are populated only when Status is MethodDecoded.
	TypeHint string
	Result   any
	// Err carries the underlying error for every non-decoded status.
	Err error
}

// InterfaceMatch reports whether a contract interface was detected and, if not, what was missing.
type InterfaceMatch struct {
	Name ContractInterface
	// Matched is true when every required result was produced and decoded.
	Matched bool
	// Required lists the result type hints this interface needs.
	Required []string
	// Missing lists the required result type hints that were not produced (method absent, reverted,
	// or its output failed to decode). When Missing is non-empty cross-reference Inspection.Methods /
	// Problems() to learn whether a required getter decode-failed rather than being truly absent.
	Missing []string
}

// Inspection is the complete, loss-free result of inspecting a contract. Unlike ContractDescription
// it records every method that was attempted together with how it failed, and explains interface
// detection. The caller decides what to make of partial or failed results.
type Inspection struct {
	CodeHash ton.Bits256
	// Methods holds one entry per method that was actually attempted (methods known to be absent
	// from the code are not attempted and are not listed).
	Methods    []MethodResult
	Interfaces []InterfaceMatch
	// MethodsError is set when the contract's get-methods dictionary could not be read from its
	// code. When non-nil, no methods were attempted because the inspector could not learn which
	// methods exist (typically an unresolved library), as opposed to the contract genuinely
	// having no get methods. This is the failure the lossy InspectContract reports as "no
	// interfaces, no methods, nil error".
	MethodsError error
}

// Decoded returns the methods that ran and decoded successfully.
func (i *Inspection) Decoded() []MethodResult {
	var out []MethodResult
	for _, m := range i.Methods {
		if m.Status == MethodDecoded {
			out = append(out, m)
		}
	}
	return out
}

// Problems returns the methods that were attempted but did not decode (transport error, reverted,
// or decode failure).
func (i *Inspection) Problems() []MethodResult {
	var out []MethodResult
	for _, m := range i.Methods {
		if m.Status != MethodDecoded {
			out = append(out, m)
		}
	}
	return out
}

// TransportErrors returns the methods whose execution failed at the transport level.
func (i *Inspection) TransportErrors() []MethodResult {
	var out []MethodResult
	for _, m := range i.Methods {
		if m.Status == MethodTransportError {
			out = append(out, m)
		}
	}
	return out
}

// ImplementedInterfaces returns the names of the matched interfaces, preserving inspection order.
func (i *Inspection) ImplementedInterfaces() []ContractInterface {
	var out []ContractInterface
	for _, iface := range i.Interfaces {
		if iface.Matched {
			out = append(out, iface.Name)
		}
	}
	return out
}

// decodedHints returns the set of result type hints produced by successfully decoded methods.
func (i *Inspection) decodedHints() map[string]struct{} {
	hints := make(map[string]struct{})
	for _, m := range i.Methods {
		if m.Status == MethodDecoded {
			hints[m.TypeHint] = struct{}{}
		}
	}
	return hints
}

// recordingExecutor wraps an Executor and remembers the outcome of the most recent call. Because
// every generated InvokeFn calls RunSmcMethodByID exactly once, this lets the inspector recover the
// exit code, raw stack and transport error that the InvokeFn's (string, any, error) signature hides.
type recordingExecutor struct {
	inner Executor
	last  *invocationRecord
}

type invocationRecord struct {
	methodID int
	exitCode uint32
	stack    tlb.VmStack
	err      error
}

func (r *recordingExecutor) RunSmcMethodByID(ctx context.Context, accountID ton.AccountID, methodID int, params tlb.VmStack) (uint32, tlb.VmStack, error) {
	exitCode, stack, err := r.inner.RunSmcMethodByID(ctx, accountID, methodID, params)
	r.last = &invocationRecord{methodID: methodID, exitCode: exitCode, stack: stack, err: err}
	return exitCode, stack, err
}

func (r *recordingExecutor) reset() { r.last = nil }

// classify turns the outcome of a single InvokeFn call into a MethodResult, using the recorded
// executor call to distinguish a transport error from a revert from a decode failure.
func classify(name string, rec *recordingExecutor, typeHint string, result any, invErr error) MethodResult {
	mr := MethodResult{Name: name}
	if rec.last == nil {
		// The InvokeFn returned before reaching the executor (e.g. it failed to build its input
		// stack). There is no VM outcome to report; treat it as a decode/setup failure.
		mr.Status = MethodDecodeFailed
		mr.Err = invErr
		return mr
	}
	mr.MethodID = rec.last.methodID
	mr.ExitCode = rec.last.exitCode
	mr.Stack = rec.last.stack
	switch {
	case rec.last.err != nil:
		mr.Status = MethodTransportError
		mr.Err = rec.last.err
	case rec.last.exitCode != 0 && rec.last.exitCode != 1:
		mr.Status = MethodReverted
		mr.Err = invErr
	case invErr != nil:
		mr.Status = MethodDecodeFailed
		mr.Err = invErr
	default:
		mr.Status = MethodDecoded
		mr.TypeHint = typeHint
		mr.Result = result
	}
	return mr
}

// InspectContract2 executes the known get methods on a contract and returns a complete, loss-free
// Inspection: every attempted method is reported with its outcome (decoded, reverted, decode-failed
// or transport error) and the raw stack is preserved on decode failures. The returned top-level
// error is non-nil only for failures that prevent inspection from starting at all (a malformed code
// cell or a library-resolution failure); per-method failures — including transport errors — are
// recorded in the Inspection and left for the caller to act on.
//
// The contract is not provided directly; the executor is expected to have been created with
// knowledge of the contract's code and data. The executor must be ready to run multiple different
// methods and must not rely on a particular order of execution.
func (ci contractInspector) InspectContract2(ctx context.Context, code []byte, executor Executor, reqAccountID ton.AccountID) (*Inspection, error) {
	if len(code) == 0 {
		return &Inspection{}, nil
	}
	info, err := GetCodeInfo(ctx, code, ci.libResolver)
	if err != nil {
		return nil, err
	}
	insp := &Inspection{CodeHash: info.Hash, MethodsError: info.MethodsError}
	rec := &recordingExecutor{inner: executor}

	// For known contracts the interfaces are fixed; we only run their get methods to collect results.
	if contract, ok := knownContracts[info.Hash]; ok {
		for _, method := range contract.getMethods {
			rec.reset()
			typeHint, result, invErr := method(ctx, rec, reqAccountID)
			insp.Methods = append(insp.Methods, classify("", rec, typeHint, result, invErr))
		}
		for _, name := range contract.contractInterfaces {
			insp.Interfaces = append(insp.Interfaces, InterfaceMatch{Name: name, Matched: true})
		}
		return insp, nil
	}

	for _, method := range ci.knownMethods {
		// let's avoid running get methods that we know don't exist
		if !ci.scanAllMethods && !info.isMethodOkToTry(method.Name) {
			continue
		}
		rec.reset()
		typeHint, result, invErr := method.InvokeFn(ctx, rec, reqAccountID)
		insp.Methods = append(insp.Methods, classify(method.Name, rec, typeHint, result, invErr))
	}

	decoded := insp.decodedHints()
	for _, iface := range ci.knownInterfaces {
		match := InterfaceMatch{Name: iface.Name, Required: iface.Results}
		for _, r := range iface.Results {
			if _, ok := decoded[r]; !ok {
				match.Missing = append(match.Missing, r)
			}
		}
		match.Matched = len(match.Missing) == 0
		insp.Interfaces = append(insp.Interfaces, match)
	}

	return insp, nil
}

// InspectContract tries to execute all known get method on a contract and returns a summary.
// The contract is not provided directly,
// instead, it is expected that "executor" has been created with knowledge of the contract's code and data.
// Executor must be ready to execute multiple different methods and must not rely on a particular order of execution.
//
// InspectContract is a lossy, backwards-compatible projection of InspectContract2: it discards
// methods that failed to execute or decode and only reports successfully decoded results and matched
// interfaces. Prefer InspectContract2 when you need to know why a method or interface was not detected.
func (ci contractInspector) InspectContract(ctx context.Context, code []byte, executor Executor, reqAccountID ton.AccountID) (*ContractDescription, error) {
	insp, err := ci.InspectContract2(ctx, code, executor, reqAccountID)
	if err != nil {
		return nil, err
	}
	// Preserve the original behaviour of aborting on lite-server client errors.
	for _, m := range insp.Methods {
		if m.Status == MethodTransportError && liteclient.IsClientError(m.Err) {
			return nil, m.Err
		}
	}
	desc := &ContractDescription{
		ContractInterfaces: insp.ImplementedInterfaces(),
		MethodsInspected:   len(insp.Methods),
	}
	for _, m := range insp.Decoded() {
		desc.GetMethods = append(desc.GetMethods, MethodInvocation{
			Result:   m.Result,
			TypeHint: m.TypeHint,
		})
	}
	return desc, nil
}

type CodeInfo struct {
	Hash    ton.Bits256
	Methods map[int64]struct{}
	// MethodsError is set when the get-methods dictionary could not be read from the code
	// (the reference to it was missing or its decoding failed, often because a referenced
	// library could not be resolved). In that case Methods is nil but GetCodeInfo still
	// returns a nil error, because a contract is allowed to expose no get methods. Inspect
	// this to tell "no get methods" apart from "we failed to read the get methods".
	MethodsError error
}

func (i CodeInfo) isMethodOkToTry(name string) bool {
	if i.Methods == nil {
		return false
	}
	methodID := utils.MethodIdFromName(name)
	_, ok := i.Methods[int64(methodID)]
	return ok
}

func GetCodeInfo(ctx context.Context, code []byte, resolver libResolver) (*CodeInfo, error) {
	cells, err := boc.DeserializeBoc(code)
	if err != nil {
		return nil, err
	}
	if len(cells) == 0 {
		return nil, fmt.Errorf("failed to find a root cell")
	}
	root := cells[0]
	libHashes, err := codePkg.FindLibraries(root)
	if err != nil {
		return nil, fmt.Errorf("failed while looking for libraries inside cell: %w", err)
	}
	var libs map[ton.Bits256]*boc.Cell
	if len(libHashes) > 0 {
		if resolver == nil {
			return nil, fmt.Errorf("found libraries in cell, but no resolver provided")
		}
		libs, err = resolver.GetLibraries(ctx, libHashes)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch libraries: %w", err)
		}
	}
	h, err := root.Hash256()
	if err != nil {
		return nil, err
	}
	root.ResetCounters()
	if root.IsLibrary() {
		hash, err := root.GetLibraryHash()
		if err != nil {
			return nil, err
		}
		cell, ok := libs[ton.Bits256(hash)]
		if !ok {
			return nil, fmt.Errorf("library not found")
		}
		root = cell
	}
	c, err := root.NextRef()
	if err != nil {
		// we are OK, if there is no information about get methods, but record why so callers
		// can distinguish a genuine absence from a failure to read the get-methods dictionary.
		return &CodeInfo{Hash: h, MethodsError: fmt.Errorf("failed to reach get-methods cell: %w", err)}, nil
	}
	if c.IsLibrary() {
		hash, err := c.GetLibraryHash()
		if err != nil {
			return nil, err
		}
		cell, ok := libs[ton.Bits256(hash)]
		if !ok {
			return nil, fmt.Errorf("library not found")
		}
		c = cell
	}
	type GetMethods struct {
		Hashmap tlb.Hashmap[tlb.Uint19, boc.Cell]
	}
	var getMethods GetMethods
	decoder := tlb.NewDecoder().WithLibraryResolver(func(hash tlb.Bits256) (*boc.Cell, error) {
		if resolver == nil {
			return nil, fmt.Errorf("failed to fetch library: no resolver provided")
		}
		cell, ok := libs[ton.Bits256(hash)]
		if ok {
			return cell, nil
		}
		localLibs, err := resolver.GetLibraries(ctx, []ton.Bits256{ton.Bits256(hash)})
		if err != nil {
			return nil, err
		}
		if len(localLibs) == 0 {
			return nil, fmt.Errorf("library not found")
		}
		return localLibs[ton.Bits256(hash)], nil
	})

	err = decoder.Unmarshal(c, &getMethods)
	if err != nil {
		// we are OK, if there is no information about get methods, but record why so callers
		// can distinguish a genuine absence from a failure to decode the get-methods dictionary.
		return &CodeInfo{Hash: h, MethodsError: fmt.Errorf("failed to decode get-methods dictionary: %w", err)}, nil
	}
	keys := getMethods.Hashmap.Keys()
	methods := make(map[int64]struct{}, len(keys))
	for _, key := range keys {
		methods[int64(key)] = struct{}{}
	}
	return &CodeInfo{Hash: h, Methods: methods}, nil
}
