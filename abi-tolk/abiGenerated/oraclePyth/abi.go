// Code generated - DO NOT EDIT.

package abiGenerated

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

const PrefixUpdateGuardianSetMessage uint64 = 0x00000001

type UpdateGuardianSetMessage struct {
	WormholeMessage boc.Cell
}

const PrefixUpdatePriceFeedsMessage uint64 = 0x00000002

type UpdatePriceFeedsMessage struct {
	UpdateData boc.Cell
}

const PrefixExecuteGovernanceActionMessage uint64 = 0x00000003

type ExecuteGovernanceActionMessage struct {
	GovernanceVm boc.Cell
}

const PrefixUpgradeContractMessage uint64 = 0x00000004

type UpgradeContractMessage struct {
	NewCode boc.Cell
}

const PrefixParsePriceFeedUpdatesMessage uint64 = 0x00000005

type ParsePriceFeedUpdatesMessage struct {
	UpdateData     boc.Cell
	PriceIds       boc.Cell
	MinPublishTime tlb.Uint64
	MaxPublishTime tlb.Uint64
	TargetAddress  tlb.MsgAddress
	CustomPayload  boc.Cell
}

const PrefixParseUniquePriceFeedUpdatesMessage uint64 = 0x00000006

type ParseUniquePriceFeedUpdatesMessage struct {
	UpdateData    boc.Cell
	PriceIds      boc.Cell
	PublishTime   tlb.Uint64
	MaxStaleness  tlb.Uint64
	TargetAddress tlb.MsgAddress
	CustomPayload boc.Cell
}
type PricePoint struct {
	Price       tlb.Int64
	Conf        tlb.Uint64
	Expo        tlb.Int32
	PublishTime tlb.Uint64
}
type StoredPriceFeed struct {
	Price    tlb.RefT[*PricePoint]
	EmaPrice tlb.RefT[*PricePoint]
}
type PriceFeedsSection struct {
	LatestPriceFeeds tlb.Hashmap[tlb.Uint256, tlb.RefT[*StoredPriceFeed]]
	SingleUpdateFee  tlb.Uint256
}
type DataSourcesSection struct {
	IsValidDataSource tlb.Hashmap[tlb.Uint256, bool]
}
type GuardianSetRecord struct {
	ExpirationTime tlb.Uint64
	GuardianKeys   tlb.Hashmap[tlb.Uint8, tlb.Bits160]
}
type GuardianSetsSection struct {
	CurrentGuardianSetIndex tlb.Uint32
	GuardianSets            tlb.Hashmap[tlb.Uint32, tlb.RefT[*GuardianSetRecord]]
}
type DataSource struct {
	EmitterChainId tlb.Uint16
	EmitterAddress tlb.Uint256
}
type GovernanceSection struct {
	ChainId                        tlb.Uint16
	GovernanceChainId              tlb.Uint16
	GovernanceContract             tlb.Uint256
	ConsumedGovernanceActions      tlb.Hashmap[tlb.Uint256, bool]
	GovernanceDataSource           tlb.RefT[*DataSource]
	LastExecutedGovernanceSequence tlb.Uint64
	GovernanceDataSourceIndex      tlb.Uint32
	UpgradeCodeHash                tlb.Uint256
}
type MainStorage struct {
	PriceFeeds   tlb.RefT[*PriceFeedsSection]
	DataSources  tlb.RefT[*DataSourcesSection]
	GuardianSets tlb.RefT[*GuardianSetsSection]
	Governance   tlb.RefT[*GovernanceSection]
}
type PriceFeedResponseEntry struct {
	PriceId   tlb.Uint256
	PriceFeed tlb.RefT[*StoredPriceFeed]
	Next      tlb.Maybe[tlb.RefT[*PriceFeedResponseEntry]]
}
type PriceFeedUpdateResponse struct {
	Op             tlb.Uint32
	PriceFeedCount tlb.Uint8
	PriceFeeds     tlb.RefT[*PriceFeedResponseEntry]
	OriginalSender tlb.MsgAddress
	CustomPayload  boc.Cell
}

const PrefixErrorResponse uint64 = 0x10002

type ErrorResponse struct {
	ErrorCode     tlb.Uint32
	Operation     tlb.Uint32
	CustomPayload boc.Cell
}

const PrefixSuccessResponse uint64 = 0x10001

type SuccessResponse struct {
	Result        boc.Cell
	CustomPayload boc.Cell
}

func (v *UpdateGuardianSetMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixUpdateGuardianSetMessage {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if v.WormholeMessage, err = c.NextRefV(); err != nil {
		return err
	}
	return nil
}
func (v *UpdatePriceFeedsMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixUpdatePriceFeedsMessage {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if v.UpdateData, err = c.NextRefV(); err != nil {
		return err
	}
	return nil
}
func (v *ExecuteGovernanceActionMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixExecuteGovernanceActionMessage {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if v.GovernanceVm, err = c.NextRefV(); err != nil {
		return err
	}
	return nil
}
func (v *UpgradeContractMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixUpgradeContractMessage {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if v.NewCode, err = c.NextRefV(); err != nil {
		return err
	}
	return nil
}
func (v *ParsePriceFeedUpdatesMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixParsePriceFeedUpdatesMessage {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if v.UpdateData, err = c.NextRefV(); err != nil {
		return err
	}
	if v.PriceIds, err = c.NextRefV(); err != nil {
		return err
	}
	if err = v.MinPublishTime.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.MaxPublishTime.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.TargetAddress.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if v.CustomPayload, err = c.NextRefV(); err != nil {
		return err
	}
	return nil
}
func (v *ParseUniquePriceFeedUpdatesMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixParseUniquePriceFeedUpdatesMessage {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if v.UpdateData, err = c.NextRefV(); err != nil {
		return err
	}
	if v.PriceIds, err = c.NextRefV(); err != nil {
		return err
	}
	if err = v.PublishTime.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.MaxStaleness.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.TargetAddress.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if v.CustomPayload, err = c.NextRefV(); err != nil {
		return err
	}
	return nil
}
func (v *PricePoint) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Price.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.Conf.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.Expo.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.PublishTime.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *StoredPriceFeed) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Price.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.EmaPrice.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *PriceFeedsSection) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.LatestPriceFeeds.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.SingleUpdateFee.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *DataSourcesSection) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.IsValidDataSource.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *GuardianSetRecord) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.ExpirationTime.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.GuardianKeys.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *GuardianSetsSection) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.CurrentGuardianSetIndex.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.GuardianSets.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *DataSource) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.EmitterChainId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.EmitterAddress.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *GovernanceSection) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.ChainId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.GovernanceChainId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.GovernanceContract.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.ConsumedGovernanceActions.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.GovernanceDataSource.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.LastExecutedGovernanceSequence.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.GovernanceDataSourceIndex.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.UpgradeCodeHash.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *MainStorage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.PriceFeeds.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.DataSources.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.GuardianSets.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.Governance.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *PriceFeedResponseEntry) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.PriceId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.PriceFeed.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.Next.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *PriceFeedUpdateResponse) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Op.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.PriceFeedCount.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.PriceFeeds.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.OriginalSender.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if v.CustomPayload, err = c.NextRefV(); err != nil {
		return err
	}
	return nil
}
func (v *ErrorResponse) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(20); err != nil {
		return err
	} else if prefix != PrefixErrorResponse {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.ErrorCode.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.Operation.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if v.CustomPayload, err = c.NextRefV(); err != nil {
		return err
	}
	return nil
}
func (v *SuccessResponse) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(20); err != nil {
		return err
	} else if prefix != PrefixSuccessResponse {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if v.Result, err = c.NextRefV(); err != nil {
		return err
	}
	if v.CustomPayload, err = c.NextRefV(); err != nil {
		return err
	}
	return nil
}
