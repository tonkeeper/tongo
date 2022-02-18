package tvm

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"tongo/boc"
)

type EntryType int

const (
	Int EntryType = iota
	Null
	Cell
	Tuple
)

type TvmStackEntry struct {
	Type     EntryType
	IntVal   big.Int
	CellVal  *boc.Cell
	TupleVal []*TvmStackEntry
}

func NewIntStackEntry(val big.Int) TvmStackEntry {
	return TvmStackEntry{
		Type:   Int,
		IntVal: val,
	}
}

func NewNullStackEntry() TvmStackEntry {
	return TvmStackEntry{
		Type: Null,
	}
}

func (e *TvmStackEntry) Int() big.Int {
	return e.IntVal
}

func (e *TvmStackEntry) Cell() *boc.Cell {
	return e.CellVal
}

func (e *TvmStackEntry) IsNull() bool {
	return e.Type == Null
}

func (e *TvmStackEntry) IsInt() bool {
	return e.Type == Int
}

func (e *TvmStackEntry) IsCell() bool {
	return e.Type == Cell
}

func (e *TvmStackEntry) IsTuple() bool {
	return e.Type == Cell
}

func (e *TvmStackEntry) UnmarshalJSON(data []byte) error {
	var m map[string]json.RawMessage

	//m := map[string]string{}

	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	var entryType string
	err = json.Unmarshal(m["type"], &entryType)
	if err != nil {
		return err
	}

	if entryType == "int" {
		var intEntry struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		}
		err = json.Unmarshal(data, &intEntry)
		if err != nil {
			return err
		}

		e.Type = Int
		e.IntVal.SetString(intEntry.Value, 10)
	} else if entryType == "null" {
		e.Type = Null
	} else if entryType == "tuple" {
		var tupleEntry struct {
			Type  string           `json:"type"`
			Value []*TvmStackEntry `json:"value"`
		}
		err = json.Unmarshal(data, &tupleEntry)
		if err != nil {
			return err
		}
		e.Type = Tuple
		e.TupleVal = tupleEntry.Value
	} else if entryType == "cell" {
		var cellEntry struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		}
		err = json.Unmarshal(data, &cellEntry)
		if err != nil {
			return err
		}

		e.Type = Cell
		cellData, err := base64.StdEncoding.DecodeString(cellEntry.Value)
		if err != nil {
			return err
		}
		parsedBoc, err := boc.DeserializeBoc(cellData)
		if err != nil {
			return err
		}
		e.CellVal = parsedBoc[0]
	} else {
		return errors.New("unknown stack entry type")
	}
	return nil
}
