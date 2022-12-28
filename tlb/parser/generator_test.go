package parser

import (
	"fmt"
	"testing"
)

var SOURCE = `
transfer#5fcc3d14 query_id:uint64 new_owner:MsgAddress response_destination:MsgAddress
custom_payload:(Maybe ^Cell) forward_amount:(VarUInteger 16) forward_payload:(Either Cell ^Cell) =
InternalMsgBody;
                
ownership_assigned#05138d91 query_id:uint64 prev_owner:MsgAddress forward_payload:(Either Cell ^Cell) = InternalMsgBody;
                
excesses#d53276db query_id:uint64 = InternalMsgBody;
				
get_static_data#2fcb26a2 query_id:uint64 = InternalMsgBody;

report_static_data#8b771735 query_id:uint64 index:uint256 collection:MsgAddress = InternalMsgBody;

text#_ {n:#} data:(SnakeData ~n) = Text;
snake#00 data:(SnakeData ~n) = ContentData;
chunks#01 data:ChunkedData = ContentData;
onchain#00 data:(HashmapE 256 ^ContentData) = FullContent;
offchain#01 uri:Text = FullContent;

get_royalty_params#693d3950 query_id:uint64 = InternalMsgBody;

report_royalty_params#a8cb00ad query_id:uint64 numerator:uint16 denominator:uint16 destination:MsgAddress = InternalMsgBody;

transfer#0f8a7ea5 query_id:uint64 amount:(VarUInteger 16) destination:MsgAddress
response_destination:MsgAddress custom_payload:(Maybe ^Cell)
forward_ton_amount:(VarUInteger 16) forward_payload:(Either Cell ^Cell)
= InternalMsgBody;

transfer_notification#7362d09c query_id:uint64 amount:(VarUInteger 16)
sender:MsgAddress forward_payload:(Either Cell ^Cell)
= InternalMsgBody;

burn#595f07bc query_id:uint64 amount:(VarUInteger 16)
response_destination:MsgAddress custom_payload:(Maybe ^Cell)
= InternalMsgBody;

dns_text#1eda _:Text = DNSRecord;
dns_next_resolver#ba93 resolver:MsgAddressInt = DNSRecord;  // usually in record #-1
dns_adnl_address#ad01 adnl_addr:bits256 flags:(## 8)
proto_list:flags . 0?ProtoList = DNSRecord;  // often in record #2
dns_smc_address#9fd3 smc_addr:MsgAddressInt flags:(## 8)
cap_list:flags . 0?SmcCapList = DNSRecord;   // often in record #1

prove_ownership#04ded148 query_id:uint64 dest:MsgAddress
forward_payload:^Cell with_content:Bool = InternalMsgBody;

ownership_proof#0524c7ae query_id:uint64 item_id:uint256 owner:MsgAddress
data:^Cell revoked_at:uint64 content:(Maybe ^Cell) = InternalMsgBody;

request_owner#d0c3bfea query_id:uint64 dest:MsgAddress
forward_payload:^Cell with_content:Bool = InternalMsgBody;

owner_info#0dd607e3 query_id:uint64 item_id:uint256 initiator:MsgAddress owner:MsgAddress
data:^Cell revoked_at:uint64 content:(Maybe ^Cell) = InternalMsgBody;
`

func TestGenerateGolangTypes(t *testing.T) {
	parsed, err := Parse(SOURCE)
	if err != nil {
		panic(err)
	}
	g := NewGenerator(nil, "LightClient")

	s, err := g.LoadTypes(parsed.Declarations)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", s)
	//_ = s
}

func TestGenerateVarUintTypes(t *testing.T) {
	fmt.Println(GenerateConstantInts(17))
}
