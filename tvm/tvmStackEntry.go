package tvm

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/startfellows/tongo/boc"
	"math/big"
)

type EntryType int

const (
	Int EntryType = iota
	Null
	Cell
	Tuple
	CellSlice
)

type TvmStackEntry struct {
	Type         EntryType
	intVal       big.Int
	cellVal      *boc.Cell
	cellSliceVal *boc.Cell
	tupleVal     []TvmStackEntry
}

type basicEntry struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type nullEntry struct {
	Type string `json:"type"`
}

type tupleEntry struct {
	Type  string          `json:"type"`
	Value []TvmStackEntry `json:"value"`
}

func NewBigIntStackEntry(val big.Int) TvmStackEntry {
	return TvmStackEntry{
		Type:   Int,
		intVal: val,
	}
}

func NewIntStackEntry(val int) TvmStackEntry {
	return TvmStackEntry{
		Type:   Int,
		intVal: *big.NewInt(int64(val)),
	}
}

func NewNullStackEntry() TvmStackEntry {
	return TvmStackEntry{
		Type: Null,
	}
}

func NewTupleStackEntry(val []TvmStackEntry) TvmStackEntry {
	return TvmStackEntry{
		Type:     Tuple,
		tupleVal: val,
	}
}

func (e *TvmStackEntry) Int() big.Int {
	return e.intVal
}

func (e *TvmStackEntry) Int64() int64 {
	return e.intVal.Int64()
}

func (e *TvmStackEntry) Uint64() uint64 {
	return e.intVal.Uint64()
}

func (e *TvmStackEntry) Cell() *boc.Cell {
	return e.cellVal
}

func (e *TvmStackEntry) Tuple() []TvmStackEntry {
	return e.tupleVal
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
	return e.Type == Tuple
}

func (e *TvmStackEntry) IsCellSlice() bool {
	return e.Type == CellSlice
}

func (e *TvmStackEntry) UnmarshalJSON(data []byte) error {
	var m map[string]json.RawMessage

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
		var intEntry basicEntry
		err = json.Unmarshal(data, &intEntry)
		if err != nil {
			return err
		}

		e.Type = Int
		e.intVal.SetString(intEntry.Value, 10)
	} else if entryType == "null" {
		e.Type = Null
	} else if entryType == "tuple" {
		var tupleEntry tupleEntry
		err = json.Unmarshal(data, &tupleEntry)
		if err != nil {
			return err
		}
		e.Type = Tuple
		e.tupleVal = tupleEntry.Value
	} else if entryType == "cell" {
		var cellEntry basicEntry
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
		e.cellVal = parsedBoc[0]
	} else if entryType == "cell_slice" {
		var cellEntry basicEntry
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
		e.cellSliceVal = parsedBoc[0]
	} else {
		return errors.New("unknown stack entry type")
	}

	return nil
}

func (e TvmStackEntry) MarshalJSON() ([]byte, error) {
	if e.Type == Int {
		return json.Marshal(&basicEntry{Type: "int", Value: e.intVal.String()})
	} else if e.Type == Cell {
		bocStr, err := e.cellVal.ToBocBase64Custom(false, true, false, 0)
		if err != nil {
			return nil, err
		}
		return json.Marshal(&basicEntry{Type: "cell", Value: bocStr})
	} else if e.Type == Tuple {
		return json.Marshal(&tupleEntry{Type: "tuple", Value: e.tupleVal})
	} else if e.Type == Null {
		return json.Marshal(&nullEntry{Type: "null"})
	}
	return nil, errors.New("unable to serialize tvm stack entry")
}
