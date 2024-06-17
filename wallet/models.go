package wallet

import (
	"context"
	"fmt"
	"time"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type Version int

const (
	V1R1 Version = iota
	V1R2
	V1R3
	V2R1
	V2R2
	V3R1
	V3R2
	V3R2Lockup
	V4R1
	V4R2
	V5Beta
	HighLoadV1R1
	HighLoadV1R2
	HighLoadV2
	HighLoadV2R1
	HighLoadV2R2
	// TODO: maybe add lockup wallet
)

const (
	// DefaultSubWallet is a recommended default value of subWalletID according to
	// https://docs.ton.org/develop/smart-contracts/tutorials/wallet#subwallet-ids.
	DefaultSubWallet       = 698983191
	DefaultMessageLifetime = time.Minute * 3
	DefaultMessageMode     = 3
)

var (
	ErrAccountIsFrozen         = fmt.Errorf("account is frozen")
	ErrAccountIsNotInitialized = fmt.Errorf("account is not initialized")
)

var codes = map[Version]string{
	V1R1:         "te6cckEBAQEARAAAhP8AIN2k8mCBAgDXGCDXCx/tRNDTH9P/0VESuvKhIvkBVBBE+RDyovgAAdMfMSDXSpbTB9QC+wDe0aTIyx/L/8ntVEH98Ik=",
	V1R2:         "te6cckEBAQEAUwAAov8AIN0gggFMl7qXMO1E0NcLH+Ck8mCBAgDXGCDXCx/tRNDTH9P/0VESuvKhIvkBVBBE+RDyovgAAdMfMSDXSpbTB9QC+wDe0aTIyx/L/8ntVNDieG8=",
	V1R3:         "te6cckEBAQEAXwAAuv8AIN0gggFMl7ohggEznLqxnHGw7UTQ0x/XC//jBOCk8mCBAgDXGCDXCx/tRNDTH9P/0VESuvKhIvkBVBBE+RDyovgAAdMfMSDXSpbTB9QC+wDe0aTIyx/L/8ntVLW4bkI=",
	V2R1:         "te6cckEBAQEAVwAAqv8AIN0gggFMl7qXMO1E0NcLH+Ck8mCDCNcYINMf0x8B+CO78mPtRNDTH9P/0VExuvKhA/kBVBBC+RDyovgAApMg10qW0wfUAvsA6NGkyMsfy//J7VShNwu2",
	V2R2:         "te6cckEBAQEAYwAAwv8AIN0gggFMl7ohggEznLqxnHGw7UTQ0x/XC//jBOCk8mCDCNcYINMf0x8B+CO78mPtRNDTH9P/0VExuvKhA/kBVBBC+RDyovgAApMg10qW0wfUAvsA6NGkyMsfy//J7VQETNeh",
	V3R1:         "te6cckEBAQEAYgAAwP8AIN0gggFMl7qXMO1E0NcLH+Ck8mCDCNcYINMf0x/TH/gjE7vyY+1E0NMf0x/T/9FRMrryoVFEuvKiBPkBVBBV+RDyo/gAkyDXSpbTB9QC+wDo0QGkyMsfyx/L/8ntVD++buA=",
	V3R2:         "te6cckEBAQEAcQAA3v8AIN0gggFMl7ohggEznLqxn3Gw7UTQ0x/THzHXC//jBOCk8mCDCNcYINMf0x/TH/gjE7vyY+1E0NMf0x/T/9FRMrryoVFEuvKiBPkBVBBV+RDyo/gAkyDXSpbTB9QC+wDo0QGkyMsfyx/L/8ntVBC9ba0=",
	V3R2Lockup:   "te6ccgECHgEAAmEAART/APSkE/S88sgLAQIBIAIDAgFIBAUB8vKDCNcYINMf0x/TH4AkA/gjuxPy8vADgCJRqboa8vSAI1G3uhvy9IAfC/kBVBDF+RAa8vT4AFBX+CPwBlCY+CPwBiBxKJMg10qOi9MHMdRRG9s8ErAB6DCSKaDfcvsCBpMg10qW0wfUAvsA6NEDpEdoFBVDMPAE7VQdAgLNBgcCASATFAIBIAgJAgEgDxACASAKCwAtXtRNDTH9Mf0//T//QE+gD0BPoA9ATRgD9wB0NMDAXGwkl8D4PpAMCHHAJJfA+AB0x8hwQKSXwTg8ANRtPABghCC6vnEUrC9sJJfDOCAKIIQgur5xBu6GvL0gCErghA7msoAvvL0B4MI1xiAICH5AVQQNvkQEvL00x+AKYIQNzqp9BO6EvL00wDTHzAB4w8QSBA3XjKAMDQ4AEwh10n0qG+lbDGAADBA5SArwBQAWEDdBCvAFCBBXUFYAEBAkQwDwBO1UAgEgERIARUjh4igCD0lm+lIJMwI7uRMeIgmDX6ANEToUATkmwh4rPmMIADUCMjKHxfKHxXL/xPL//QAAfoC9AAB+gL0AMmAAQxRIqBTE4Ag9A5voZb6ANEToAKRMOLIUAP6AkATgCD0QwGACASAVFgAVven3gBiCQvhHgAwCASAXGAIBSBscAC21GH4AbYiGioJgngDGIH4Axj8E7eILMAIBWBkaABetznaiaGmfmOuF/8AAF6x49qJoaY+Y64WPwAARsyX7UTQ1wsfgABex0b4I4IBCMPtQ9iAAKAHQ0wMBeLCSW3/g+kAx+kAwAfAB",
	V4R1:         "te6cckECFQEAAvUAART/APSkE/S88sgLAQIBIAIDAgFIBAUE+PKDCNcYINMf0x/THwL4I7vyY+1E0NMf0x/T//QE0VFDuvKhUVG68qIF+QFUEGT5EPKj+AAkpMjLH1JAyx9SMMv/UhD0AMntVPgPAdMHIcAAn2xRkyDXSpbTB9QC+wDoMOAhwAHjACHAAuMAAcADkTDjDQOkyMsfEssfy/8REhMUA+7QAdDTAwFxsJFb4CHXScEgkVvgAdMfIYIQcGx1Z70ighBibG5jvbAighBkc3RyvbCSXwPgAvpAMCD6RAHIygfL/8nQ7UTQgQFA1yH0BDBcgQEI9ApvoTGzkl8F4ATTP8glghBwbHVnupEx4w0kghBibG5juuMABAYHCAIBIAkKAFAB+gD0BDCCEHBsdWeDHrFwgBhQBcsFJ88WUAP6AvQAEstpyx9SEMs/AFL4J28ighBibG5jgx6xcIAYUAXLBSfPFiT6AhTLahPLH1Iwyz8B+gL0AACSghBkc3Ryuo41BIEBCPRZMO1E0IEBQNcgyAHPFvQAye1UghBkc3Rygx6xcIAYUATLBVjPFiL6AhLLassfyz+UEDRfBOLJgED7AAIBIAsMAFm9JCtvaiaECAoGuQ+gIYRw1AgIR6STfSmRDOaQPp/5g3gSgBt4EBSJhxWfMYQCAVgNDgARuMl+1E0NcLH4AD2ynftRNCBAUDXIfQEMALIygfL/8nQAYEBCPQKb6ExgAgEgDxAAGa3OdqJoQCBrkOuF/8AAGa8d9qJoQBBrkOuFj8AAbtIH+gDU1CL5AAXIygcVy//J0Hd0gBjIywXLAiLPFlAF+gIUy2sSzMzJcfsAyEAUgQEI9FHypwIAbIEBCNcYyFQgJYEBCPRR8qeCEG5vdGVwdIAYyMsFywJQBM8WghAF9eEA+gITy2oSyx/JcfsAAgBygQEI1xgwUgKBAQj0WfKn+CWCEGRzdHJwdIAYyMsFywJQBc8WghAF9eEA+gIUy2oTyx8Syz/Jc/sAAAr0AMntVEap808=",
	V4R2:         "te6cckECFAEAAtQAART/APSkE/S88sgLAQIBIAIDAgFIBAUE+PKDCNcYINMf0x/THwL4I7vyZO1E0NMf0x/T//QE0VFDuvKhUVG68qIF+QFUEGT5EPKj+AAkpMjLH1JAyx9SMMv/UhD0AMntVPgPAdMHIcAAn2xRkyDXSpbTB9QC+wDoMOAhwAHjACHAAuMAAcADkTDjDQOkyMsfEssfy/8QERITAubQAdDTAyFxsJJfBOAi10nBIJJfBOAC0x8hghBwbHVnvSKCEGRzdHK9sJJfBeAD+kAwIPpEAcjKB8v/ydDtRNCBAUDXIfQEMFyBAQj0Cm+hMbOSXwfgBdM/yCWCEHBsdWe6kjgw4w0DghBkc3RyupJfBuMNBgcCASAICQB4AfoA9AQw+CdvIjBQCqEhvvLgUIIQcGx1Z4MesXCAGFAEywUmzxZY+gIZ9ADLaRfLH1Jgyz8gyYBA+wAGAIpQBIEBCPRZMO1E0IEBQNcgyAHPFvQAye1UAXKwjiOCEGRzdHKDHrFwgBhQBcsFUAPPFiP6AhPLassfyz/JgED7AJJfA+ICASAKCwBZvSQrb2omhAgKBrkPoCGEcNQICEekk30pkQzmkD6f+YN4EoAbeBAUiYcVnzGEAgFYDA0AEbjJftRNDXCx+AA9sp37UTQgQFA1yH0BDACyMoHy//J0AGBAQj0Cm+hMYAIBIA4PABmtznaiaEAga5Drhf/AABmvHfaiaEAQa5DrhY/AAG7SB/oA1NQi+QAFyMoHFcv/ydB3dIAYyMsFywIizxZQBfoCFMtrEszMyXP7AMhAFIEBCPRR8qcCAHCBAQjXGPoA0z/IVCBHgQEI9FHyp4IQbm90ZXB0gBjIywXLAlAGzxZQBPoCFMtqEssfyz/Jc/sAAgBsgQEI1xj6ANM/MFIkgQEI9Fnyp4IQZHN0cnB0gBjIywXLAlAFzxZQA/oCE8tqyx8Syz/Jc/sAAAr0AMntVGliJeU=",
	V5Beta:       "te6ccgEBAQEAIwAIQgLkzzsvTG1qYeoPK1RH0mZ4WyavNjfbLe7mvNGqgm80Eg==",
	HighLoadV1R1: "te6ccgEBBgEAhgABFP8A9KQT9KDyyAsBAgEgAgMCAUgEBQC88oMI1xgg0x/TH9Mf+CMTu/Jj7UTQ0x/TH9P/0VEyuvKhUUS68qIE+QFUEFX5EPKj9ATR+AB/jhghgBD0eG+hb6EgmALTB9QwAfsAkTLiAbPmWwGkyMsfyx/L/8ntVAAE0DAAEaCZL9qJoa4WPw==",
	HighLoadV1R2: "te6ccgEBCAEAmQABFP8A9KQT9LzyyAsBAgEgAgMCAUgEBQC88oMI1xgg0x/TH9Mf+CMTu/Jj7UTQ0x/TH9P/0VEyuvKhUUS68qIE+QFUEFX5EPKj9ATR+AB/jhghgBD0eG+hb6EgmALTB9QwAfsAkTLiAbPmWwGkyMsfyx/L/8ntVAAE0DACAUgGBwAXuznO1E0NM/MdcL/4ABG4yX7UTQ1wsfg=",
	HighLoadV2:   "te6ccgEBCQEA5QABFP8A9KQT9LzyyAsBAgEgAgcCAUgDBAAE0DACASAFBgAXvZznaiaGmvmOuF/8AEG+X5dqJoaY+Y6Z/p/5j6AmipEEAgegc30JjJLb/JXdHxQB6vKDCNcYINMf0z/4I6ofUyC58mPtRNDTH9M/0//0BNFTYIBA9A5voTHyYFFzuvKiB/kBVBCH+RDyowL0BNH4AH+OFiGAEPR4b6UgmALTB9QwAfsAkTLiAbPmW4MlochANIBA9EOK5jEByMsfE8s/y//0AMntVAgANCCAQPSWb6VsEiCUMFMDud4gkzM2AZJsIeKz",
	HighLoadV2R1: "te6ccgEBBwEA1gABFP8A9KQT9KDyyAsBAgEgAgMCAUgEBQHu8oMI1xgg0x/TP/gjqh9TILnyY+1E0NMf0z/T//QE0VNggED0Dm+hMfJgUXO68qIH+QFUEIf5EPKjAvQE0fgAf44YIYAQ9HhvoW+hIJgC0wfUMAH7AJEy4gGz5luDJaHIQDSAQPRDiuYxyBLLHxPLP8v/9ADJ7VQGAATQMABBoZfl2omhpj5jpn+n/mPoCaKkQQCB6BzfQmMktv8ld0fFADgggED0lm+hb6EyURCUMFMDud4gkzM2AZIyMOKz",
	HighLoadV2R2: "te6ccgEBCQEA6QABFP8A9KQT9LzyyAsBAgEgAgMCAUgEBQHu8oMI1xgg0x/TP/gjqh9TILnyY+1E0NMf0z/T//QE0VNggED0Dm+hMfJgUXO68qIH+QFUEIf5EPKjAvQE0fgAf44YIYAQ9HhvoW+hIJgC0wfUMAH7AJEy4gGz5luDJaHIQDSAQPRDiuYxyBLLHxPLP8v/9ADJ7VQIAATQMAIBIAYHABe9nOdqJoaa+Y64X/wAQb5fl2omhpj5jpn+n/mPoCaKkQQCB6BzfQmMktv8ld0fFAA4IIBA9JZvoW+hMlEQlDBTA7neIJMzNgGSMjDisw==",
}

// codeHashToVersion maps code's hash to a wallet version.
var codeHashToVersion = map[tlb.Bits256]Version{}

func init() {
	for ver := range codes {
		codeHashToVersion[GetCodeHashByVer(ver)] = ver
	}
}

// GetWalletVersion returns a wallet version by the given state of an account and an incoming message to the account.
// An incoming message is needed in case when a wallet has not been initialized yet.
// In this case, we take its code from the message's StateInit.
func GetWalletVersion(state tlb.ShardAccount, msg tlb.Message) (Version, bool, error) {
	if state.Account.SumType == "AccountNone" || state.Account.Account.Storage.State.SumType == "AccountUninit" {
		if !msg.Init.Exists {
			return 0, false, ErrAccountIsNotInitialized
		}
		if !msg.Init.Value.Value.Code.Exists {
			return 0, false, ErrAccountIsNotInitialized
		}
		code := msg.Init.Value.Value.Code.Value.Value
		hash, err := code.Hash256()
		if err != nil {
			return 0, false, ErrAccountIsNotInitialized
		}
		ver, ok := GetVerByCodeHash(hash)
		return ver, ok, nil
	}
	if state.Account.Account.Storage.State.SumType == "AccountFrozen" {
		return 0, false, ErrAccountIsFrozen
	}
	code := state.Account.Account.Storage.State.AccountActive.StateInit.Code
	if code.Exists {
		hash, err := code.Value.Value.Hash256()
		if err != nil {
			return 0, false, err
		}
		ver, ok := GetVerByCodeHash(hash)
		return ver, ok, nil
	}
	return 0, false, ErrAccountIsNotInitialized
}

type blockchain interface {
	GetSeqno(ctx context.Context, account ton.AccountID) (uint32, error)
	SendMessage(ctx context.Context, payload []byte) (uint32, error)
	GetAccountState(ctx context.Context, accountID ton.AccountID) (tlb.ShardAccount, error)
}

func GetCodeByVer(ver Version) *boc.Cell {
	c, err := boc.DeserializeBocBase64(codes[ver])
	if err != nil {
		panic("invalid wallet hardcoded code")
	}
	if len(c) != 1 {
		panic("code must have one root cell")
	}
	return c[0]
}

func GetCodeHashByVer(ver Version) tlb.Bits256 {
	code := GetCodeByVer(ver)
	h, err := code.Hash()
	if err != nil {
		panic("can not calc hash for hardcoded code")
	}
	var hash tlb.Bits256
	copy(hash[:], h[:])
	return hash
}

// GetVerByCodeHash returns (Version, true) if there is code with the given hash.
// Otherwise, it returns (0, false).
func GetVerByCodeHash(hash tlb.Bits256) (Version, bool) {
	if ver, ok := codeHashToVersion[hash]; ok {
		return ver, true
	}
	return 0, false
}

func (v Version) ToString() string {
	names := []string{"v1R1", "v1R2", "v1R3", "v2R1", "v2R2", "v3R1", "v3R2", "v4R1", "v4R2", "v5R1", "highload_v1R1", "highload_v1R2", "highload_v2", "highload_v2R1", "highload_v2R2"}
	if int(v) > len(names) {
		panic("to string conversion for this ver not supported")
	}
	return names[v]
}

type Sendable interface {
	ToInternal() (tlb.Message, uint8, error)
}

type SimpleTransfer struct {
	Amount     tlb.Grams
	Address    ton.AccountID
	Comment    string
	Bounceable bool
}

func (m SimpleTransfer) ToInternal() (message tlb.Message, mode uint8, err error) {
	info := tlb.CommonMsgInfo{
		SumType: "IntMsgInfo",
	}

	info.IntMsgInfo = &struct {
		IhrDisabled bool
		Bounce      bool
		Bounced     bool
		Src         tlb.MsgAddress
		Dest        tlb.MsgAddress
		Value       tlb.CurrencyCollection
		IhrFee      tlb.Grams
		FwdFee      tlb.Grams
		CreatedLt   uint64
		CreatedAt   uint32
	}{
		IhrDisabled: true,
		Bounce:      m.Bounceable,
		Src:         (*ton.AccountID)(nil).ToMsgAddress(),
		Dest:        m.Address.ToMsgAddress(),
	}
	info.IntMsgInfo.Value.Grams = m.Amount

	intMsg := tlb.Message{
		Info: info,
	}

	if m.Comment != "" {
		body := boc.NewCell()
		err := tlb.Marshal(body, TextComment(m.Comment))
		if err != nil {
			return tlb.Message{}, 0, err
		}
		intMsg.Body.IsRight = true //todo: check length and
		intMsg.Body.Value = tlb.Any(*body)
	}
	return intMsg, DefaultMessageMode, nil
}

type Message struct {
	Amount  tlb.Grams
	Address ton.AccountID
	Body    *boc.Cell
	Code    *boc.Cell
	Data    *boc.Cell
	Bounce  bool
	Mode    uint8
}

func (m Message) ToInternal() (message tlb.Message, mode uint8, err error) {
	info := tlb.CommonMsgInfo{
		SumType: "IntMsgInfo",
	}

	info.IntMsgInfo = &struct {
		IhrDisabled bool
		Bounce      bool
		Bounced     bool
		Src         tlb.MsgAddress
		Dest        tlb.MsgAddress
		Value       tlb.CurrencyCollection
		IhrFee      tlb.Grams
		FwdFee      tlb.Grams
		CreatedLt   uint64
		CreatedAt   uint32
	}{
		IhrDisabled: true,
		Bounce:      m.Bounce,
		Src:         (*ton.AccountID)(nil).ToMsgAddress(),
		Dest:        m.Address.ToMsgAddress(),
	}
	info.IntMsgInfo.Value.Grams = m.Amount

	intMsg := tlb.Message{
		Info: info,
	}

	if m.Body != nil {
		intMsg.Body.IsRight = true //todo: check length and
		intMsg.Body.Value = tlb.Any(*m.Body)
	}
	if m.Code != nil && m.Data != nil {
		intMsg.Init.Exists = true
		intMsg.Init.Value.IsRight = true
		intMsg.Init.Value.Value.Code.Exists = true
		intMsg.Init.Value.Value.Data.Exists = true
		intMsg.Init.Value.Value.Code.Value.Value = *m.Code
		intMsg.Init.Value.Value.Data.Value.Value = *m.Data
	}

	return intMsg, m.Mode, nil
}

type TextComment string

func (t TextComment) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error { // TODO: implement for binary comment
	err := c.WriteUint(0, 32) // text comment tag
	if err != nil {
		return err
	}
	return tlb.Marshal(c, tlb.Text(t))
}

func (t *TextComment) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error { // TODO: implement for binary comment
	val, err := c.ReadUint(32) // text comment tag
	if err != nil {
		return err
	}
	if val != 0 {
		return fmt.Errorf("not a text comment")
	}
	var text tlb.Text
	err = tlb.Unmarshal(c, &text)
	if err != nil {
		return err
	}
	*t = TextComment(text)
	return nil
}
