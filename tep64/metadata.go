package tep64

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

type ContentLayout int

const (
	Undefined ContentLayout = iota
	OffChain
	OnChain
	SemiChain
)

type Metadata struct {
	// Uri points to JSON document with metadata. Used by SemiChain layout. ASCII string.
	Uri         string `json:"uri,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Image       string `json:"image,omitempty"`
	ImageData   []byte `json:"image_data,omitempty"`
	Symbol      string `json:"symbol,omitempty"`
	Decimals    string `json:"decimals,omitempty"`
	RenderType  string `json:"render_type,omitempty"`
	AmountStyle string `json:"amount_style,omitempty"`
}

func (m *Metadata) Merge(other *Metadata) {
	if other == nil {
		return
	}
	if other.Uri != "" {
		m.Uri = other.Uri
	}
	if other.Name != "" {
		m.Name = other.Name
	}
	if other.Description != "" {
		m.Description = other.Description
	}
	if other.Image != "" {
		m.Image = other.Image
	}
	if other.ImageData != nil {
		m.ImageData = other.ImageData
	}
	if other.Symbol != "" {
		m.Symbol = other.Symbol
	}
	if other.Decimals != "" {
		m.Decimals = other.Decimals
	}
	if other.RenderType != "" {
		m.RenderType = other.RenderType
	}
	if other.AmountStyle != "" {
		m.AmountStyle = other.AmountStyle
	}
}

// FullContent is either a link to metadata or metadata itself depending on the layout.
type FullContent struct {
	Layout ContentLayout
	Data   []byte
	// OnchainMetadata contains a decoded metadata when the layout is onchain.
	OnchainMetadata *Metadata
	// OffchainURL contains a link to JSON when the layout is offchain.
	OffchainURL string
}

var (
	ErrUnsupportedContentType = errors.New("unsupported content type")
)

// TEP-64 Token Data Standard
// https://github.com/ton-blockchain/TEPs/blob/master/text/0064-token-data-standard.md
func ConvertOnchainData(content tlb.FullContent) (Metadata, error) {
	if content.SumType != "Onchain" {
		return Metadata{}, fmt.Errorf("not Onchain content")
	}
	var m Metadata
	for i, v := range content.Onchain.Data.Values() {
		keyS := hex.EncodeToString(content.Onchain.Data.Keys()[i][:])
		switch keyS {
		case "70e5d7b6a29b392f85076fe15ca2f2053c56c2338728c4e33c9e8ddb1ee827cc": // sha256(uri)
			b, err := v.Value.Bytes()
			if err != nil {
				return Metadata{}, err
			}
			m.Uri = string(b)
		case "82a3537ff0dbce7eec35d69edc3a189ee6f17d82f353a553f9aa96cb0be3ce89": // sha256(name)
			b, err := v.Value.Bytes()
			if err != nil {
				return Metadata{}, err
			}
			m.Name = string(b)
		case "c9046f7a37ad0ea7cee73355984fa5428982f8b37c8f7bcec91f7ac71a7cd104": // sha256(description)
			b, err := v.Value.Bytes()
			if err != nil {
				return Metadata{}, err
			}
			m.Description = string(b)
		case "6105d6cc76af400325e94d588ce511be5bfdbb73b437dc51eca43917d7a43e3d": // sha256(image)
			b, err := v.Value.Bytes()
			if err != nil {
				return Metadata{}, err
			}
			m.Image = string(b)
		case "d9a88ccec79eef59c84b671136a20ece4cd00caaad5bc47e2c208829154ee9e4": // sha256(image_data)
			b, err := v.Value.Bytes()
			if err != nil {
				return Metadata{}, err
			}
			m.ImageData = b
		case "b76a7ca153c24671658335bbd08946350ffc621fa1c516e7123095d4ffd5c581": // sha256(symbol)
			b, err := v.Value.Bytes()
			if err != nil {
				return Metadata{}, err
			}
			m.Symbol = string(b)
		case "ee80fd2f1e03480e2282363596ee752d7bb27f50776b95086a0279189675923e": // sha256(decimals)
			b, err := v.Value.Bytes()
			if err != nil {
				return Metadata{}, err
			}
			m.Decimals = string(b)
		case "d33ae06043036d0d1c3be27201ac15ee4c73da8cdb7c8f3462ce308026095ac0": // sha256(render_type)
			b, err := v.Value.Bytes()
			if err != nil {
				return Metadata{}, err
			}
			m.RenderType = string(b)
		case "8b10e058ce46c44bc1ba139bc9761721e49170e2c0a176129250a70af053b700": // sha256(amount_style)
			b, err := v.Value.Bytes()
			if err != nil {
				return Metadata{}, err
			}
			m.AmountStyle = string(b)
		}
	}
	return m, nil
}

func DecodeFullContentFromCell(cell *boc.Cell) (FullContent, error) {
	var content tlb.FullContent
	err := tlb.Unmarshal(cell, &content)
	if err != nil {
		return FullContent{}, fmt.Errorf("%v content decoding: %v", content.SumType, err)
	}
	return DecodeFullContent(content)
}

func DecodeFullContent(content tlb.FullContent) (FullContent, error) {
	switch content.SumType {
	case "Onchain":
		meta, err := ConvertOnchainData(content)
		if err != nil {
			return FullContent{}, err
		}
		layout := OnChain
		offchainURL := ""
		if len(meta.Uri) > 0 {
			// according to https://github.com/ton-blockchain/TEPs/blob/master/text/0064-token-data-standard.md
			// if uri is present, then the layout is SemiChain
			layout = SemiChain
		}
		result, err := json.Marshal(meta)
		if err != nil {
			return FullContent{}, err
		}
		return FullContent{
			Layout:          layout,
			Data:            result,
			OnchainMetadata: &meta,
			OffchainURL:     offchainURL,
		}, nil

	case "Offchain":
		bs := boc.BitString(content.Offchain.Uri)
		if bs.BitsAvailableForRead()%8 != 0 {
			return FullContent{}, fmt.Errorf("text data is not multiple of 8 bits")
		}
		result, err := bs.GetTopUppedArray()
		if err != nil {
			return FullContent{}, err
		}
		return FullContent{
			Layout:      OffChain,
			Data:        result,
			OffchainURL: string(result),
		}, nil
	}
	return FullContent{}, ErrUnsupportedContentType
}
