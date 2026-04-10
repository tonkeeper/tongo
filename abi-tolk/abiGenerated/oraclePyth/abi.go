// Code generated - DO NOT EDIT.

package abiGenerated

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

const PrefixUpdateGuardianSetMessage = 0x00000001

type UpdateGuardianSetMessage struct {
	WormholeMessage boc.Cell `json:"wormholeMessage"`
}

func (v *UpdateGuardianSetMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != uint64(0x00000001) {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if v.WormholeMessage, err = c.NextRef(); err != nil {
		return err
	}
	return nil
}

const PrefixUpdatePriceFeedsMessage = 0x00000002

type UpdatePriceFeedsMessage struct {
	UpdateData boc.Cell `json:"updateData"`
}

func (v *UpdatePriceFeedsMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != uint64(0x00000002) {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if v.UpdateData, err = c.NextRef(); err != nil {
		return err
	}
	return nil
}

const PrefixExecuteGovernanceActionMessage = 0x00000003

type ExecuteGovernanceActionMessage struct {
	GovernanceVm boc.Cell `json:"governanceVm"`
}

func (v *ExecuteGovernanceActionMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != uint64(0x00000003) {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if v.GovernanceVm, err = c.NextRef(); err != nil {
		return err
	}
	return nil
}

const PrefixUpgradeContractMessage = 0x00000004

type UpgradeContractMessage struct {
	NewCode boc.Cell `json:"newCode"`
}

func (v *UpgradeContractMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != uint64(0x00000004) {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if v.NewCode, err = c.NextRef(); err != nil {
		return err
	}
	return nil
}

const PrefixParsePriceFeedUpdatesMessage = 0x00000005

type ParsePriceFeedUpdatesMessage struct {
	UpdateData     boc.Cell       `json:"updateData"`
	PriceIds       boc.Cell       `json:"priceIds"`
	MinPublishTime tlb.Uint64     `json:"minPublishTime"`
	MaxPublishTime tlb.Uint64     `json:"maxPublishTime"`
	TargetAddress  tlb.MsgAddress `json:"targetAddress"`
	CustomPayload  boc.Cell       `json:"customPayload"`
}

func (v *ParsePriceFeedUpdatesMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != uint64(0x00000005) {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if v.UpdateData, err = c.NextRef(); err != nil {
		return err
	}
	if v.PriceIds, err = c.NextRef(); err != nil {
		return err
	}
	if v.MinPublishTime, err = tlb.UnmarshalT[tlb.Uint64](c, decoder); err != nil {
		return err
	}
	if v.MaxPublishTime, err = tlb.UnmarshalT[tlb.Uint64](c, decoder); err != nil {
		return err
	}
	if v.TargetAddress, err = tlb.UnmarshalT[tlb.MsgAddress](c, decoder); err != nil {
		return err
	}
	if v.CustomPayload, err = c.NextRef(); err != nil {
		return err
	}
	return nil
}

const PrefixParseUniquePriceFeedUpdatesMessage = 0x00000006

type ParseUniquePriceFeedUpdatesMessage struct {
	UpdateData    boc.Cell       `json:"updateData"`
	PriceIds      boc.Cell       `json:"priceIds"`
	PublishTime   tlb.Uint64     `json:"publishTime"`
	MaxStaleness  tlb.Uint64     `json:"maxStaleness"`
	TargetAddress tlb.MsgAddress `json:"targetAddress"`
	CustomPayload boc.Cell       `json:"customPayload"`
}

func (v *ParseUniquePriceFeedUpdatesMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != uint64(0x00000006) {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if v.UpdateData, err = c.NextRef(); err != nil {
		return err
	}
	if v.PriceIds, err = c.NextRef(); err != nil {
		return err
	}
	if v.PublishTime, err = tlb.UnmarshalT[tlb.Uint64](c, decoder); err != nil {
		return err
	}
	if v.MaxStaleness, err = tlb.UnmarshalT[tlb.Uint64](c, decoder); err != nil {
		return err
	}
	if v.TargetAddress, err = tlb.UnmarshalT[tlb.MsgAddress](c, decoder); err != nil {
		return err
	}
	if v.CustomPayload, err = c.NextRef(); err != nil {
		return err
	}
	return nil
}

type PricePoint struct {
	Price       tlb.Int64  `json:"price"`
	Conf        tlb.Uint64 `json:"conf"`
	Expo        tlb.Int32  `json:"expo"`
	PublishTime tlb.Uint64 `json:"publishTime"`
}

func (v *PricePoint) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if v.Price, err = tlb.UnmarshalT[tlb.Int64](c, decoder); err != nil {
		return err
	}
	if v.Conf, err = tlb.UnmarshalT[tlb.Uint64](c, decoder); err != nil {
		return err
	}
	if v.Expo, err = tlb.UnmarshalT[tlb.Int32](c, decoder); err != nil {
		return err
	}
	if v.PublishTime, err = tlb.UnmarshalT[tlb.Uint64](c, decoder); err != nil {
		return err
	}
	return nil
}

type StoredPriceFeed struct {
	Price    *boc.Cell `json:"price"`
	EmaPrice *boc.Cell `json:"emaPrice"`
}

func (v *StoredPriceFeed) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if v.Price, err = loadCellOf(); err != nil {
		return err
	}
	if v.EmaPrice, err = loadCellOf(); err != nil {
		return err
	}
	return nil
}

type PriceFeedsSection struct {
	LatestPriceFeeds tlb.Hashmap[tlb.Uint256, *boc.Cell] `json:"latestPriceFeeds"`
	SingleUpdateFee  tlb.Uint256                         `json:"singleUpdateFee"`
}

func (v *PriceFeedsSection) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if v.LatestPriceFeeds, err = tlb.UnmarshalT[tlb.Hashmap[tlb.Uint256, *boc.Cell]](c, decoder); err != nil {
		return err
	}
	if v.SingleUpdateFee, err = tlb.UnmarshalT[tlb.Uint256](c, decoder); err != nil {
		return err
	}
	return nil
}

type DataSourcesSection struct {
	IsValidDataSource tlb.Hashmap[tlb.Uint256, bool] `json:"isValidDataSource"`
}

func (v *DataSourcesSection) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if v.IsValidDataSource, err = tlb.UnmarshalT[tlb.Hashmap[tlb.Uint256, bool]](c, decoder); err != nil {
		return err
	}
	return nil
}

type GuardianSetRecord struct {
	ExpirationTime tlb.Uint64                          `json:"expirationTime"`
	GuardianKeys   tlb.Hashmap[tlb.Uint8, tlb.Bits160] `json:"guardianKeys"`
}

func (v *GuardianSetRecord) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if v.ExpirationTime, err = tlb.UnmarshalT[tlb.Uint64](c, decoder); err != nil {
		return err
	}
	if v.GuardianKeys, err = tlb.UnmarshalT[tlb.Hashmap[tlb.Uint8, tlb.Bits160]](c, decoder); err != nil {
		return err
	}
	return nil
}

type GuardianSetsSection struct {
	CurrentGuardianSetIndex tlb.Uint32                         `json:"currentGuardianSetIndex"`
	GuardianSets            tlb.Hashmap[tlb.Uint32, *boc.Cell] `json:"guardianSets"`
}

func (v *GuardianSetsSection) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if v.CurrentGuardianSetIndex, err = tlb.UnmarshalT[tlb.Uint32](c, decoder); err != nil {
		return err
	}
	if v.GuardianSets, err = tlb.UnmarshalT[tlb.Hashmap[tlb.Uint32, *boc.Cell]](c, decoder); err != nil {
		return err
	}
	return nil
}

type DataSource struct {
	EmitterChainId tlb.Uint16  `json:"emitterChainId"`
	EmitterAddress tlb.Uint256 `json:"emitterAddress"`
}

func (v *DataSource) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if v.EmitterChainId, err = tlb.UnmarshalT[tlb.Uint16](c, decoder); err != nil {
		return err
	}
	if v.EmitterAddress, err = tlb.UnmarshalT[tlb.Uint256](c, decoder); err != nil {
		return err
	}
	return nil
}

type GovernanceSection struct {
	ChainId                        tlb.Uint16                     `json:"chainId"`
	GovernanceChainId              tlb.Uint16                     `json:"governanceChainId"`
	GovernanceContract             tlb.Uint256                    `json:"governanceContract"`
	ConsumedGovernanceActions      tlb.Hashmap[tlb.Uint256, bool] `json:"consumedGovernanceActions"`
	GovernanceDataSource           *boc.Cell                      `json:"governanceDataSource"`
	LastExecutedGovernanceSequence tlb.Uint64                     `json:"lastExecutedGovernanceSequence"`
	GovernanceDataSourceIndex      tlb.Uint32                     `json:"governanceDataSourceIndex"`
	UpgradeCodeHash                tlb.Uint256                    `json:"upgradeCodeHash"`
}

func (v *GovernanceSection) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if v.ChainId, err = tlb.UnmarshalT[tlb.Uint16](c, decoder); err != nil {
		return err
	}
	if v.GovernanceChainId, err = tlb.UnmarshalT[tlb.Uint16](c, decoder); err != nil {
		return err
	}
	if v.GovernanceContract, err = tlb.UnmarshalT[tlb.Uint256](c, decoder); err != nil {
		return err
	}
	if v.ConsumedGovernanceActions, err = tlb.UnmarshalT[tlb.Hashmap[tlb.Uint256, bool]](c, decoder); err != nil {
		return err
	}
	if v.GovernanceDataSource, err = loadCellOf(); err != nil {
		return err
	}
	if v.LastExecutedGovernanceSequence, err = tlb.UnmarshalT[tlb.Uint64](c, decoder); err != nil {
		return err
	}
	if v.GovernanceDataSourceIndex, err = tlb.UnmarshalT[tlb.Uint32](c, decoder); err != nil {
		return err
	}
	if v.UpgradeCodeHash, err = tlb.UnmarshalT[tlb.Uint256](c, decoder); err != nil {
		return err
	}
	return nil
}

type MainStorage struct {
	PriceFeeds   *boc.Cell `json:"priceFeeds"`
	DataSources  *boc.Cell `json:"dataSources"`
	GuardianSets *boc.Cell `json:"guardianSets"`
	Governance   *boc.Cell `json:"governance"`
}

func (v *MainStorage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if v.PriceFeeds, err = loadCellOf(); err != nil {
		return err
	}
	if v.DataSources, err = loadCellOf(); err != nil {
		return err
	}
	if v.GuardianSets, err = loadCellOf(); err != nil {
		return err
	}
	if v.Governance, err = loadCellOf(); err != nil {
		return err
	}
	return nil
}

type PriceFeedResponseEntry struct {
	PriceId   tlb.Uint256 `json:"priceId"`
	PriceFeed *boc.Cell   `json:"priceFeed"`
	Next      **boc.Cell  `json:"next"`
}

func (v *PriceFeedResponseEntry) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if v.PriceId, err = tlb.UnmarshalT[tlb.Uint256](c, decoder); err != nil {
		return err
	}
	if v.PriceFeed, err = loadCellOf(); err != nil {
		return err
	}
	if v.Next, err = (func() (value *boc.Cell, err error) {
		if isNotNull, err := c.ReadBit(); err != nil {
			return value, err
		} else if isNotNull {
			return loadCellOf(), nil
		}
		return
	})(); err != nil {
		return err
	}
	return nil
}

type PriceFeedUpdateResponse struct {
	Op             tlb.Uint32     `json:"op"`
	PriceFeedCount tlb.Uint8      `json:"priceFeedCount"`
	PriceFeeds     *boc.Cell      `json:"priceFeeds"`
	OriginalSender tlb.MsgAddress `json:"originalSender"`
	CustomPayload  *boc.Cell      `json:"customPayload"`
}

func (v *PriceFeedUpdateResponse) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if v.Op, err = tlb.UnmarshalT[tlb.Uint32](c, decoder); err != nil {
		return err
	}
	if v.PriceFeedCount, err = tlb.UnmarshalT[tlb.Uint8](c, decoder); err != nil {
		return err
	}
	if v.PriceFeeds, err = loadCellOf(); err != nil {
		return err
	}
	if v.OriginalSender, err = tlb.UnmarshalT[tlb.MsgAddress](c, decoder); err != nil {
		return err
	}
	if v.CustomPayload, err = loadCellOf(); err != nil {
		return err
	}
	return nil
}

const PrefixErrorResponse = 0x10002

type ErrorResponse struct {
	ErrorCode     tlb.Uint32 `json:"errorCode"`
	Operation     tlb.Uint32 `json:"operation"`
	CustomPayload *boc.Cell  `json:"customPayload"`
}

func (v *ErrorResponse) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(20); err != nil {
		return err
	} else if prefix != uint64(0x10002) {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if v.ErrorCode, err = tlb.UnmarshalT[tlb.Uint32](c, decoder); err != nil {
		return err
	}
	if v.Operation, err = tlb.UnmarshalT[tlb.Uint32](c, decoder); err != nil {
		return err
	}
	if v.CustomPayload, err = loadCellOf(); err != nil {
		return err
	}
	return nil
}

const PrefixSuccessResponse = 0x10001

type SuccessResponse struct {
	Result        boc.Cell `json:"result"`
	CustomPayload boc.Cell `json:"customPayload"`
}

func (v *SuccessResponse) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(20); err != nil {
		return err
	} else if prefix != uint64(0x10001) {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if v.Result, err = c.NextRef(); err != nil {
		return err
	}
	if v.CustomPayload, err = c.NextRef(); err != nil {
		return err
	}
	return nil
}
