package parser

import (
	"fmt"
	"github.com/startfellows/tongo/tlb"
	tlbParser "github.com/startfellows/tongo/tlb/parser"
	"github.com/startfellows/tongo/utils"
	"go/format"
	"sort"
	"strings"
)

var defaultKnownTypes = map[string]string{
	"accountid": "tongo.AccountID",
	"cell":      "boc.Cell",
	"int8":      "int8",
	"int64":     "int64",
	"bool":      "bool",
	"uint16":    "uint16",
	"uint64":    "uint64",
	"int256":    "tlb.Int256",
	"int257":    "tlb.Int257",
	"any":       "tlb.Any",
	"[]byte":    "[]byte",
	"big.int":   "big.Int",
}

var (
	msgDecoderReturnErr = "if err != nil {return \"\", nil, err}\n"
)

type TLBMsgBody struct {
	TypePrefix string
	TypeName   string
	Tag        uint64
	Code       string
}

type Generator struct {
	knownTypes        map[string]string
	loadedTlbTypes    []string
	loadedTlbMsgTypes map[uint32]TLBMsgBody
	typeName          string
}

func NewGenerator(knownTypes map[string]string, typeName string) *Generator {
	if knownTypes == nil {
		knownTypes = defaultKnownTypes
	}
	return &Generator{
		knownTypes:        knownTypes,
		typeName:          typeName,
		loadedTlbMsgTypes: make(map[uint32]TLBMsgBody),
	}
}

func (g *Generator) GetMethods(interfaces []Interface) (string, error) {
	var builder strings.Builder
	for _, i := range interfaces {
		for _, m := range i.Methods {
			s, err := g.GetMethod(m)
			if err != nil {
				return "", err
			}
			builder.WriteString(s)
			builder.WriteRune('\n')
		}
	}
	return builder.String(), nil
}

func (g *Generator) RegisterTypes(interfaces []Interface) error {
	for _, i := range interfaces {
		if i.Types != "" {
			err := g.registerType(i.Types)
			if err != nil {
				return err
			}
		}

		for _, internal := range i.Internals {
			err := g.registerMsgType(i.Name, internal.Input)
			if err != nil {
				return err
			}
			for _, out := range internal.Outputs {
				err := g.registerMsgType(i.Name, out)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (g *Generator) registerType(s string) error {
	tlbData, err := tlbParser.Parse(s)
	if err != nil {
		return fmt.Errorf("can't decode %v error %w", s, err)
	}
	if len(tlbData.Declarations) == 0 {
		return fmt.Errorf("can't parse type %v", s)
	}

	gen := tlbParser.NewGenerator(nil, "")
	types, err := gen.LoadTypes(tlbData.Declarations, "", false)
	if err != nil {
		return fmt.Errorf("load types error: %v", err)
	}

	g.loadedTlbTypes = append(g.loadedTlbTypes, types)

	return nil
}

func (g *Generator) registerMsgType(name, s string) error {
	parsed, err := tlbParser.Parse(s)
	if err != nil {
		return fmt.Errorf("can't decode %v error %w", s, err)
	}
	if len(parsed.Declarations) != 1 {
		return fmt.Errorf("must be only one declaration for MsgBody %v", s)
	}

	gen := tlbParser.NewGenerator(nil, name)

	tag, err := tlb.ParseTag(parsed.Declarations[0].Constructor.Prefix)
	if err != nil {
		return fmt.Errorf("can't decode tag error %w", err)
	}
	if tag.Len != 32 {
		return fmt.Errorf("message %s body tag must be 32 bit lenght", parsed.Declarations[0].Constructor.Name)
	}

	var typePrefix string
	_, ok := g.loadedTlbMsgTypes[uint32(tag.Val)]
	if ok {
		typePrefix = utils.ToCamelCase(parsed.Declarations[0].Constructor.Name)
	} else {
		typePrefix = fmt.Sprintf("%s%s", utils.ToCamelCase(name), utils.ToCamelCase(parsed.Declarations[0].Constructor.Name))
	}

	t, err := gen.LoadTypes(parsed.Declarations, typePrefix, true)
	if err != nil {
		return fmt.Errorf("can't decode %v error %w", s, err)
	}

	g.loadedTlbMsgTypes[uint32(tag.Val)] = TLBMsgBody{
		TypePrefix: typePrefix,
		TypeName:   utils.ToCamelCase(parsed.Declarations[0].Combinator.Name),
		Tag:        tag.Val,
		Code:       t,
	}

	return nil
}

func (g *Generator) checkType(s string) (string, error) {
	if typeName, prs := g.knownTypes[strings.ToLower(s)]; prs {
		return typeName, nil
	}
	tlb, err := tlbParser.Parse(s)
	if err != nil {
		return "", fmt.Errorf("can't decoder %v error %w", s, err)
	}
	if len(tlb.Declarations) == 0 {
		return "", fmt.Errorf("can't parse type %v", s)
	}
	g.loadedTlbTypes = append(g.loadedTlbTypes, s)
	return tlb.Declarations[len(tlb.Declarations)-1].Combinator.Name, nil
}

func (g *Generator) GetMethod(m GetMethod) (string, error) {
	var builder strings.Builder
	var args []string

	builder.WriteString(fmt.Sprintf("func (c %v) %v(", g.typeName, utils.ToCamelCase(m.Name)))

	for _, s := range m.Input.StackValues {
		t, err := g.checkType(s.Type)
		if err != nil {
			return "", err
		}
		args = append(args, fmt.Sprintf("%v %v", utils.ToCamelCasePrivate(s.Name), t))
	}
	builder.WriteString(strings.Join(args, ", "))
	builder.WriteString(") (")
	//for _, o := range m.Output {
	//	for _, c := range o.Stack {
	//		t, err := g.checkType(c.Type)
	//		if err != nil {
	//			return "", err
	//		}
	//		result = append(result, t)
	//	}
	//}
	//
	//result = append(result, "error")
	//builder.WriteString(strings.Join(result, ", "))
	builder.WriteString(") {\n}\n")
	//b, err := format.Source([]byte(builder.String()))
	//if err != nil {
	//	return "", err
	//}
	//return string(b), nil
	return builder.String(), nil
}

func (g *Generator) sortedMessages() []TLBMsgBody {
	var keys []int
	for k := range g.loadedTlbMsgTypes {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	var res []TLBMsgBody
	for _, k := range keys {
		res = append(res, g.loadedTlbMsgTypes[uint32(k)])
	}
	return res
}

func (g *Generator) CollectedTypes() string {
	var builder strings.Builder
	builder.WriteString(strings.Join(g.loadedTlbTypes, "\n"))

	for _, m := range g.sortedMessages() {
		builder.WriteString(m.Code)
		builder.WriteRune('\n')
	}
	builder.WriteRune('\n')
	b, err := format.Source([]byte(builder.String()))
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (g *Generator) GenerateMsgDecoder() string {
	var builder strings.Builder

	builder.WriteString("func MessageDecoder(cell *boc.Cell) (string, any, error) {\n")

	builder.WriteString("tag, err := cell.ReadUint(32)\n")
	builder.WriteString(msgDecoderReturnErr)

	builder.WriteString("switch tag {\n")

	for _, m := range g.sortedMessages() {
		builder.WriteString(fmt.Sprintf("case 0x%x:\n", m.Tag))
		builder.WriteString(fmt.Sprintf("var res %s%s\n", m.TypePrefix, m.TypeName))
		builder.WriteString("err = tlb.Unmarshal(cell, &res)\n")
		builder.WriteString(fmt.Sprintf("return \"%s\", res, err\n", m.TypePrefix))
	}

	builder.WriteString("}\n")
	builder.WriteString("return \"\", nil, fmt.Errorf(\"invalid message tag\")\n")
	builder.WriteString("}\n")
	builder.WriteRune('\n')

	b, err := format.Source([]byte(builder.String()))
	if err != nil {
		panic(err)
	}
	return string(b)
}
