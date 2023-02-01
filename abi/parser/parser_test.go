package parser

import (
	"fmt"
	"os"
	"testing"
)

var METHOD = `




<internal name="transfer">
        <input>
            transfer#5fcc3d14 query_id:uint64 new_owner:MsgAddress response_destination:MsgAddress custom_payload:(Maybe ^Cell) forward_amount:(VarUInteger 16) forward_payload:(Either Cell ^Cell) = InternalMsgBody;
        </input>
        <ouput>
            ownership_assigned#05138d91 query_id:uint64 prev_owner:MsgAddress forward_payload:(Either Cell ^Cell) = InternalMsgBody;
        </ouput>
        <ouput>
            excesses#d53276db query_id:uint64 = InternalMsgBody;
        </ouput>
    </internal>
`

func TestParseMethod(t *testing.T) {
	i, err := ParseMethod([]byte(METHOD))
	fmt.Println(i, err)
}

func TestParseInterface(t *testing.T) {
	b, err := os.ReadFile("../known.xml")
	if err != nil {
		t.Fatal(err)
	}
	i, err := ParseInterface(b)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(i)
}
