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
	if v.WormholeMessage, err = (func() (boc.Cell, error) {
		cref, err := c.NextRef()
		if err != nil {
			return boc.Cell{}, err
		}
		return *cref, nil
	}()); err != nil {
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
	if v.UpdateData, err = (func() (boc.Cell, error) {
		cref, err := c.NextRef()
		if err != nil {
			return boc.Cell{}, err
		}
		return *cref, nil
	}()); err != nil {
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
	if v.GovernanceVm, err = (func() (boc.Cell, error) {
		cref, err := c.NextRef()
		if err != nil {
			return boc.Cell{}, err
		}
		return *cref, nil
	}()); err != nil {
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
	if v.NewCode, err = (func() (boc.Cell, error) {
		cref, err := c.NextRef()
		if err != nil {
			return boc.Cell{}, err
		}
		return *cref, nil
	}()); err != nil {
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
	if v.UpdateData, err = (func() (boc.Cell, error) {
		cref, err := c.NextRef()
		if err != nil {
			return boc.Cell{}, err
		}
		return *cref, nil
	}()); err != nil {
		return err
	}
	if v.PriceIds, err = (func() (boc.Cell, error) {
		cref, err := c.NextRef()
		if err != nil {
			return boc.Cell{}, err
		}
		return *cref, nil
	}()); err != nil {
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
	if v.CustomPayload, err = (func() (boc.Cell, error) {
		cref, err := c.NextRef()
		if err != nil {
			return boc.Cell{}, err
		}
		return *cref, nil
	}()); err != nil {
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
	if v.UpdateData, err = (func() (boc.Cell, error) {
		cref, err := c.NextRef()
		if err != nil {
			return boc.Cell{}, err
		}
		return *cref, nil
	}()); err != nil {
		return err
	}
	if v.PriceIds, err = (func() (boc.Cell, error) {
		cref, err := c.NextRef()
		if err != nil {
			return boc.Cell{}, err
		}
		return *cref, nil
	}()); err != nil {
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
	if v.CustomPayload, err = (func() (boc.Cell, error) {
		cref, err := c.NextRef()
		if err != nil {
			return boc.Cell{}, err
		}
		return *cref, nil
	}()); err != nil {
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

type StoredPriceFeed struct {
	Price    tlb.RefT[*PricePoint] `json:"price"`
	EmaPrice tlb.RefT[*PricePoint] `json:"emaPrice"`
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

type PriceFeedsSection struct {
	LatestPriceFeeds tlb.Hashmap[tlb.Uint256, tlb.RefT[*StoredPriceFeed]] `json:"latestPriceFeeds"`
	SingleUpdateFee  tlb.Uint256                                          `json:"singleUpdateFee"`
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

type DataSourcesSection struct {
	IsValidDataSource tlb.Hashmap[tlb.Uint256, bool] `json:"isValidDataSource"`
}

func (v *DataSourcesSection) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.IsValidDataSource.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}

type GuardianSetRecord struct {
	ExpirationTime tlb.Uint64                          `json:"expirationTime"`
	GuardianKeys   tlb.Hashmap[tlb.Uint8, tlb.Bits160] `json:"guardianKeys"`
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

type GuardianSetsSection struct {
	CurrentGuardianSetIndex tlb.Uint32                                            `json:"currentGuardianSetIndex"`
	GuardianSets            tlb.Hashmap[tlb.Uint32, tlb.RefT[*GuardianSetRecord]] `json:"guardianSets"`
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

type DataSource struct {
	EmitterChainId tlb.Uint16  `json:"emitterChainId"`
	EmitterAddress tlb.Uint256 `json:"emitterAddress"`
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

type GovernanceSection struct {
	ChainId                        tlb.Uint16                     `json:"chainId"`
	GovernanceChainId              tlb.Uint16                     `json:"governanceChainId"`
	GovernanceContract             tlb.Uint256                    `json:"governanceContract"`
	ConsumedGovernanceActions      tlb.Hashmap[tlb.Uint256, bool] `json:"consumedGovernanceActions"`
	GovernanceDataSource           tlb.RefT[*DataSource]          `json:"governanceDataSource"`
	LastExecutedGovernanceSequence tlb.Uint64                     `json:"lastExecutedGovernanceSequence"`
	GovernanceDataSourceIndex      tlb.Uint32                     `json:"governanceDataSourceIndex"`
	UpgradeCodeHash                tlb.Uint256                    `json:"upgradeCodeHash"`
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

type MainStorage struct {
	PriceFeeds   tlb.RefT[*PriceFeedsSection]   `json:"priceFeeds"`
	DataSources  tlb.RefT[*DataSourcesSection]  `json:"dataSources"`
	GuardianSets tlb.RefT[*GuardianSetsSection] `json:"guardianSets"`
	Governance   tlb.RefT[*GovernanceSection]   `json:"governance"`
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

type PriceFeedResponseEntry struct {
	PriceId   tlb.Uint256                                  `json:"priceId"`
	PriceFeed tlb.RefT[*StoredPriceFeed]                   `json:"priceFeed"`
	Next      tlb.Maybe[tlb.RefT[*PriceFeedResponseEntry]] `json:"next"`
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

type PriceFeedUpdateResponse struct {
	Op             tlb.Uint32                        `json:"op"`
	PriceFeedCount tlb.Uint8                         `json:"priceFeedCount"`
	PriceFeeds     tlb.RefT[*PriceFeedResponseEntry] `json:"priceFeeds"`
	OriginalSender tlb.MsgAddress                    `json:"originalSender"`
	CustomPayload  boc.Cell                          `json:"customPayload"`
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
	if v.CustomPayload, err = (func() (boc.Cell, error) {
		cref, err := c.NextRef()
		if err != nil {
			return boc.Cell{}, err
		}
		return *cref, nil
	}()); err != nil {
		return err
	}
	return nil
}

const PrefixErrorResponse = 0x10002

type ErrorResponse struct {
	ErrorCode     tlb.Uint32 `json:"errorCode"`
	Operation     tlb.Uint32 `json:"operation"`
	CustomPayload boc.Cell   `json:"customPayload"`
}

func (v *ErrorResponse) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(20); err != nil {
		return err
	} else if prefix != uint64(0x10002) {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.ErrorCode.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.Operation.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if v.CustomPayload, err = (func() (boc.Cell, error) {
		cref, err := c.NextRef()
		if err != nil {
			return boc.Cell{}, err
		}
		return *cref, nil
	}()); err != nil {
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
	if v.Result, err = (func() (boc.Cell, error) {
		cref, err := c.NextRef()
		if err != nil {
			return boc.Cell{}, err
		}
		return *cref, nil
	}()); err != nil {
		return err
	}
	if v.CustomPayload, err = (func() (boc.Cell, error) {
		cref, err := c.NextRef()
		if err != nil {
			return boc.Cell{}, err
		}
		return *cref, nil
	}()); err != nil {
		return err
	}
	return nil
}
