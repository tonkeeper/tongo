package boc

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestCell_MarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		createCell func() (*Cell, error)
	}{
		{
			name: "one cell",
			createCell: func() (*Cell, error) {
				c := NewCell()
				if err := c.WriteBytes([]byte("hello")); err != nil {
					return nil, err
				}
				if err := c.WriteInt(34, 64); err != nil {
					return nil, err
				}
				return c, nil
			},
		},
		{
			name: "cell with refs",
			createCell: func() (*Cell, error) {
				c := NewCell()
				if err := c.WriteBytes([]byte("hello")); err != nil {
					return nil, err
				}
				if err := c.WriteInt(99, 64); err != nil {
					return nil, err
				}
				for i := 0; i <= 3; i++ {
					ref := NewCell()
					if err := c.WriteBytes([]byte(fmt.Sprintf("ref %v", i))); err != nil {
						return nil, err
					}
					if err := c.AddRef(ref); err != nil {
						return nil, err
					}
				}
				return c, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell, err := tt.createCell()
			if err != nil {
				t.Fatalf("failed to create a cell: %v", err)
				return
			}
			type testType struct {
				C *Cell
			}
			data := testType{C: cell}
			bytes, err := json.Marshal(data)
			if err != nil {
				t.Fatalf("json.Marshal() failed: %v", err)
				return
			}
			var unmarshalled testType
			if err := json.Unmarshal(bytes, &unmarshalled); err != nil {
				t.Fatalf("json.Unmarshal() failed: %v", err)
				return
			}
			if unmarshalled.C.ToString() != data.C.ToString() {
				t.Fatalf("want: %v, got: %v", data.C.ToString(), unmarshalled.C.ToString())
			}
		})
	}
}
