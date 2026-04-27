// Code generated - DO NOT EDIT.

package abiPythOracle

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *UpdateGuardianSetMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixUpdateGuardianSetMessage); err != nil {
		return err
	}
	if v.WormholeMessage, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .WormholeMessage: %v", err)
	}
	return nil
}
func (v UpdateGuardianSetMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixUpdateGuardianSetMessage, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = c.AddRef(&v.WormholeMessage); err != nil {
		return fmt.Errorf("failed to .WormholeMessage: %v", err)
	}
	return nil
}
func (v UpdateGuardianSetMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *UpdatePriceFeedsMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixUpdatePriceFeedsMessage); err != nil {
		return err
	}
	if v.UpdateData, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .UpdateData: %v", err)
	}
	return nil
}
func (v UpdatePriceFeedsMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixUpdatePriceFeedsMessage, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = c.AddRef(&v.UpdateData); err != nil {
		return fmt.Errorf("failed to .UpdateData: %v", err)
	}
	return nil
}
func (v UpdatePriceFeedsMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ExecuteGovernanceActionMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixExecuteGovernanceActionMessage); err != nil {
		return err
	}
	if v.GovernanceVm, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .GovernanceVm: %v", err)
	}
	return nil
}
func (v ExecuteGovernanceActionMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixExecuteGovernanceActionMessage, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = c.AddRef(&v.GovernanceVm); err != nil {
		return fmt.Errorf("failed to .GovernanceVm: %v", err)
	}
	return nil
}
func (v ExecuteGovernanceActionMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *UpgradeContractMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixUpgradeContractMessage); err != nil {
		return err
	}
	if v.NewCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .NewCode: %v", err)
	}
	return nil
}
func (v UpgradeContractMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixUpgradeContractMessage, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = c.AddRef(&v.NewCode); err != nil {
		return fmt.Errorf("failed to .NewCode: %v", err)
	}
	return nil
}
func (v UpgradeContractMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *PriceFeedMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.MessageType.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MessageType: %v", err)
	}
	if err = v.PriceId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PriceId: %v", err)
	}
	if err = v.Price.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Price: %v", err)
	}
	if err = v.Conf.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Conf: %v", err)
	}
	if err = v.Expo.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Expo: %v", err)
	}
	if err = v.PublishTime.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PublishTime: %v", err)
	}
	if err = v.PrevPublishTime.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PrevPublishTime: %v", err)
	}
	if err = v.EmaPrice.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .EmaPrice: %v", err)
	}
	if err = v.EmaConf.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .EmaConf: %v", err)
	}
	return nil
}
func (v PriceFeedMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.MessageType.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MessageType: %v", err)
	}
	if err = v.PriceId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PriceId: %v", err)
	}
	if err = v.Price.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Price: %v", err)
	}
	if err = v.Conf.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Conf: %v", err)
	}
	if err = v.Expo.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Expo: %v", err)
	}
	if err = v.PublishTime.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PublishTime: %v", err)
	}
	if err = v.PrevPublishTime.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PrevPublishTime: %v", err)
	}
	if err = v.EmaPrice.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .EmaPrice: %v", err)
	}
	if err = v.EmaConf.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .EmaConf: %v", err)
	}
	return nil
}
func (v PriceFeedMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *PriceFeedUpdateData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Magic.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Magic: %v", err)
	}
	if err = v.MajorVersion.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MajorVersion: %v", err)
	}
	if err = v.MinorVersion.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinorVersion: %v", err)
	}
	if err = v.Accumulator.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Accumulator: %v", err)
	}
	return nil
}
func (v PriceFeedUpdateData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Magic.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Magic: %v", err)
	}
	if err = v.MajorVersion.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MajorVersion: %v", err)
	}
	if err = v.MinorVersion.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinorVersion: %v", err)
	}
	if err = v.Accumulator.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Accumulator: %v", err)
	}
	return nil
}
func (v PriceFeedUpdateData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ParsePriceFeedUpdatesMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixParsePriceFeedUpdatesMessage); err != nil {
		return err
	}
	if err = v.UpdateData.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UpdateData: %v", err)
	}
	if err = v.PriceIds.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PriceIds: %v", err)
	}
	if err = v.MinPublishTime.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinPublishTime: %v", err)
	}
	if err = v.MaxPublishTime.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MaxPublishTime: %v", err)
	}
	if err = v.TargetAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TargetAddress: %v", err)
	}
	if v.CustomPayload, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .CustomPayload: %v", err)
	}
	return nil
}
func (v ParsePriceFeedUpdatesMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixParsePriceFeedUpdatesMessage, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.UpdateData.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UpdateData: %v", err)
	}
	if err = v.PriceIds.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PriceIds: %v", err)
	}
	if err = v.MinPublishTime.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinPublishTime: %v", err)
	}
	if err = v.MaxPublishTime.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MaxPublishTime: %v", err)
	}
	if err = v.TargetAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TargetAddress: %v", err)
	}
	if err = c.AddRef(&v.CustomPayload); err != nil {
		return fmt.Errorf("failed to .CustomPayload: %v", err)
	}
	return nil
}
func (v ParsePriceFeedUpdatesMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ParseUniquePriceFeedUpdatesMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixParseUniquePriceFeedUpdatesMessage); err != nil {
		return err
	}
	if err = v.UpdateData.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UpdateData: %v", err)
	}
	if err = v.PriceIds.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PriceIds: %v", err)
	}
	if err = v.PublishTime.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PublishTime: %v", err)
	}
	if err = v.MaxStaleness.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MaxStaleness: %v", err)
	}
	if err = v.TargetAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TargetAddress: %v", err)
	}
	if v.CustomPayload, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .CustomPayload: %v", err)
	}
	return nil
}
func (v ParseUniquePriceFeedUpdatesMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixParseUniquePriceFeedUpdatesMessage, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.UpdateData.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UpdateData: %v", err)
	}
	if err = v.PriceIds.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PriceIds: %v", err)
	}
	if err = v.PublishTime.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PublishTime: %v", err)
	}
	if err = v.MaxStaleness.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MaxStaleness: %v", err)
	}
	if err = v.TargetAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TargetAddress: %v", err)
	}
	if err = c.AddRef(&v.CustomPayload); err != nil {
		return fmt.Errorf("failed to .CustomPayload: %v", err)
	}
	return nil
}
func (v ParseUniquePriceFeedUpdatesMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *PricePoint) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Price.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Price: %v", err)
	}
	if err = v.Conf.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Conf: %v", err)
	}
	if err = v.Expo.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Expo: %v", err)
	}
	if err = v.PublishTime.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PublishTime: %v", err)
	}
	return nil
}
func (v PricePoint) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Price.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Price: %v", err)
	}
	if err = v.Conf.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Conf: %v", err)
	}
	if err = v.Expo.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Expo: %v", err)
	}
	if err = v.PublishTime.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PublishTime: %v", err)
	}
	return nil
}
func (v PricePoint) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *PricePoint) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.PublishTime.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .PublishTime: %v", err)
	}
	if err = v.Expo.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Expo: %v", err)
	}
	if err = v.Conf.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Conf: %v", err)
	}
	if err = v.Price.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Price: %v", err)
	}
	return nil
}
func (v *StoredPriceFeed) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Price.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Price: %v", err)
	}
	if err = v.EmaPrice.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .EmaPrice: %v", err)
	}
	return nil
}
func (v StoredPriceFeed) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Price.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Price: %v", err)
	}
	if err = v.EmaPrice.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .EmaPrice: %v", err)
	}
	return nil
}
func (v StoredPriceFeed) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *PriceFeedsSection) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.LatestPriceFeeds.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .LatestPriceFeeds: %v", err)
	}
	if err = v.SingleUpdateFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SingleUpdateFee: %v", err)
	}
	return nil
}
func (v PriceFeedsSection) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.LatestPriceFeeds.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .LatestPriceFeeds: %v", err)
	}
	if err = v.SingleUpdateFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SingleUpdateFee: %v", err)
	}
	return nil
}
func (v PriceFeedsSection) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *DataSourcesSection) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.IsValidDataSource.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .IsValidDataSource: %v", err)
	}
	return nil
}
func (v DataSourcesSection) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.IsValidDataSource.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .IsValidDataSource: %v", err)
	}
	return nil
}
func (v DataSourcesSection) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *GuardianSetRecord) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.ExpirationTime.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ExpirationTime: %v", err)
	}
	if err = v.GuardianKeys.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .GuardianKeys: %v", err)
	}
	return nil
}
func (v GuardianSetRecord) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.ExpirationTime.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ExpirationTime: %v", err)
	}
	if err = v.GuardianKeys.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .GuardianKeys: %v", err)
	}
	return nil
}
func (v GuardianSetRecord) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *GuardianSetsSection) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.CurrentGuardianSetIndex.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .CurrentGuardianSetIndex: %v", err)
	}
	if err = v.GuardianSets.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .GuardianSets: %v", err)
	}
	return nil
}
func (v GuardianSetsSection) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.CurrentGuardianSetIndex.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .CurrentGuardianSetIndex: %v", err)
	}
	if err = v.GuardianSets.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .GuardianSets: %v", err)
	}
	return nil
}
func (v GuardianSetsSection) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *DataSource) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.EmitterChainId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .EmitterChainId: %v", err)
	}
	if err = v.EmitterAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .EmitterAddress: %v", err)
	}
	return nil
}
func (v DataSource) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.EmitterChainId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .EmitterChainId: %v", err)
	}
	if err = v.EmitterAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .EmitterAddress: %v", err)
	}
	return nil
}
func (v DataSource) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *GovernanceSection) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.ChainId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ChainId: %v", err)
	}
	if err = v.GovernanceChainId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .GovernanceChainId: %v", err)
	}
	if err = v.GovernanceContract.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .GovernanceContract: %v", err)
	}
	if err = v.ConsumedGovernanceActions.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ConsumedGovernanceActions: %v", err)
	}
	if err = v.GovernanceDataSource.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .GovernanceDataSource: %v", err)
	}
	if err = v.LastExecutedGovernanceSequence.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .LastExecutedGovernanceSequence: %v", err)
	}
	if err = v.GovernanceDataSourceIndex.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .GovernanceDataSourceIndex: %v", err)
	}
	if err = v.UpgradeCodeHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UpgradeCodeHash: %v", err)
	}
	return nil
}
func (v GovernanceSection) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.ChainId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ChainId: %v", err)
	}
	if err = v.GovernanceChainId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .GovernanceChainId: %v", err)
	}
	if err = v.GovernanceContract.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .GovernanceContract: %v", err)
	}
	if err = v.ConsumedGovernanceActions.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ConsumedGovernanceActions: %v", err)
	}
	if err = v.GovernanceDataSource.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .GovernanceDataSource: %v", err)
	}
	if err = v.LastExecutedGovernanceSequence.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .LastExecutedGovernanceSequence: %v", err)
	}
	if err = v.GovernanceDataSourceIndex.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .GovernanceDataSourceIndex: %v", err)
	}
	if err = v.UpgradeCodeHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UpgradeCodeHash: %v", err)
	}
	return nil
}
func (v GovernanceSection) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *MainStorage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.PriceFeeds.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PriceFeeds: %v", err)
	}
	if err = v.DataSources.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .DataSources: %v", err)
	}
	if err = v.GuardianSets.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .GuardianSets: %v", err)
	}
	if err = v.Governance.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Governance: %v", err)
	}
	return nil
}
func (v MainStorage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.PriceFeeds.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PriceFeeds: %v", err)
	}
	if err = v.DataSources.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .DataSources: %v", err)
	}
	if err = v.GuardianSets.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .GuardianSets: %v", err)
	}
	if err = v.Governance.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Governance: %v", err)
	}
	return nil
}
func (v MainStorage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *PriceFeedResponseEntry) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.PriceId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PriceId: %v", err)
	}
	if err = v.PriceFeed.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PriceFeed: %v", err)
	}
	if err = v.Next.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Next: %v", err)
	}
	return nil
}
func (v PriceFeedResponseEntry) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.PriceId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PriceId: %v", err)
	}
	if err = v.PriceFeed.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PriceFeed: %v", err)
	}
	if err = v.Next.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Next: %v", err)
	}
	return nil
}
func (v PriceFeedResponseEntry) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *PriceFeedUpdateResponse) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Op.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Op: %v", err)
	}
	if err = v.PriceFeedCount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PriceFeedCount: %v", err)
	}
	if err = v.PriceFeeds.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PriceFeeds: %v", err)
	}
	if err = v.OriginalSender.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .OriginalSender: %v", err)
	}
	if v.CustomPayload, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .CustomPayload: %v", err)
	}
	return nil
}
func (v PriceFeedUpdateResponse) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Op.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Op: %v", err)
	}
	if err = v.PriceFeedCount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PriceFeedCount: %v", err)
	}
	if err = v.PriceFeeds.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PriceFeeds: %v", err)
	}
	if err = v.OriginalSender.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .OriginalSender: %v", err)
	}
	if err = c.AddRef(&v.CustomPayload); err != nil {
		return fmt.Errorf("failed to .CustomPayload: %v", err)
	}
	return nil
}
func (v PriceFeedUpdateResponse) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ErrorResponse) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(20, PrefixErrorResponse); err != nil {
		return err
	}
	if err = v.ErrorCode.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ErrorCode: %v", err)
	}
	if err = v.Operation.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Operation: %v", err)
	}
	if v.CustomPayload, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .CustomPayload: %v", err)
	}
	return nil
}
func (v ErrorResponse) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixErrorResponse, 20); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.ErrorCode.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ErrorCode: %v", err)
	}
	if err = v.Operation.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Operation: %v", err)
	}
	if err = c.AddRef(&v.CustomPayload); err != nil {
		return fmt.Errorf("failed to .CustomPayload: %v", err)
	}
	return nil
}
func (v ErrorResponse) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *SuccessResponse) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(20, PrefixSuccessResponse); err != nil {
		return err
	}
	if v.Result, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .Result: %v", err)
	}
	if v.CustomPayload, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .CustomPayload: %v", err)
	}
	return nil
}
func (v SuccessResponse) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixSuccessResponse, 20); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = c.AddRef(&v.Result); err != nil {
		return fmt.Errorf("failed to .Result: %v", err)
	}
	if err = c.AddRef(&v.CustomPayload); err != nil {
		return fmt.Errorf("failed to .CustomPayload: %v", err)
	}
	return nil
}
func (v SuccessResponse) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *PriceData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Price.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Price: %v", err)
	}
	if err = v.Conf.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Conf: %v", err)
	}
	if err = v.Expo.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Expo: %v", err)
	}
	if err = v.Timestamp.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Timestamp: %v", err)
	}
	return nil
}
func (v PriceData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Price.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Price: %v", err)
	}
	if err = v.Conf.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Conf: %v", err)
	}
	if err = v.Expo.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Expo: %v", err)
	}
	if err = v.Timestamp.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Timestamp: %v", err)
	}
	return nil
}
func (v PriceData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *PriceFeesCell) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.AssetID.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AssetID: %v", err)
	}
	if err = v.PriceData.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PriceData: %v", err)
	}
	return nil
}
func (v PriceFeesCell) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.AssetID.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AssetID: %v", err)
	}
	if err = v.PriceData.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PriceData: %v", err)
	}
	return nil
}
func (v PriceFeesCell) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *OracleResponseSuccess) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixOracleResponseSuccess); err != nil {
		return err
	}
	if err = v.SomeNum.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SomeNum: %v", err)
	}
	if err = v.PriceFeedsCell.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PriceFeedsCell: %v", err)
	}
	if err = v.InitialSender.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .InitialSender: %v", err)
	}
	if v.AfterOperation, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .AfterOperation: %v", err)
	}
	return nil
}
func (v OracleResponseSuccess) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixOracleResponseSuccess, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.SomeNum.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SomeNum: %v", err)
	}
	if err = v.PriceFeedsCell.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PriceFeedsCell: %v", err)
	}
	if err = v.InitialSender.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .InitialSender: %v", err)
	}
	if err = c.AddRef(&v.AfterOperation); err != nil {
		return fmt.Errorf("failed to .AfterOperation: %v", err)
	}
	return nil
}
func (v OracleResponseSuccess) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *GuardianSetInfo) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.ExpirationTime.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ExpirationTime: %v", err)
	}
	if v.KeysDict, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .KeysDict: %v", err)
	}
	if err = v.KeyCount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .KeyCount: %v", err)
	}
	return nil
}
func (v GuardianSetInfo) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.ExpirationTime.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ExpirationTime: %v", err)
	}
	if err = c.AddRef(&v.KeysDict); err != nil {
		return fmt.Errorf("failed to .KeysDict: %v", err)
	}
	if err = v.KeyCount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .KeyCount: %v", err)
	}
	return nil
}
func (v GuardianSetInfo) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *GuardianSetInfo) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.KeyCount.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .KeyCount: %v", err)
	}
	if v.KeysDict, err = stack.ReadCell(); err != nil {
		return fmt.Errorf("failed to read .KeysDict: %v", err)
	}
	if err = v.ExpirationTime.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ExpirationTime: %v", err)
	}
	return nil
}

func (msg UpdateGuardianSetMessage) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*MainStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg UpdatePriceFeedsMessage) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*MainStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg ExecuteGovernanceActionMessage) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*MainStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg UpgradeContractMessage) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*MainStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg ParsePriceFeedUpdatesMessage) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*MainStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg ParseUniquePriceFeedUpdatesMessage) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*MainStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}
