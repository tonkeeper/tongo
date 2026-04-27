// Code generated - DO NOT EDIT.

package abiCocoon

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *AddWorkerType) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixAddWorkerType); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.WorkerHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .WorkerHash: %v", err)
	}
	return nil
}
func (v AddWorkerType) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixAddWorkerType, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.WorkerHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .WorkerHash: %v", err)
	}
	return nil
}
func (v AddWorkerType) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *DelWorkerType) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixDelWorkerType); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.WorkerHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .WorkerHash: %v", err)
	}
	return nil
}
func (v DelWorkerType) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixDelWorkerType, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.WorkerHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .WorkerHash: %v", err)
	}
	return nil
}
func (v DelWorkerType) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *AddModelType) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixAddModelType); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.ModelHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ModelHash: %v", err)
	}
	return nil
}
func (v AddModelType) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixAddModelType, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.ModelHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ModelHash: %v", err)
	}
	return nil
}
func (v AddModelType) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *DelModelType) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixDelModelType); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.ModelHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ModelHash: %v", err)
	}
	return nil
}
func (v DelModelType) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixDelModelType, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.ModelHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ModelHash: %v", err)
	}
	return nil
}
func (v DelModelType) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *AddProxyType) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixAddProxyType); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.ProxyHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProxyHash: %v", err)
	}
	return nil
}
func (v AddProxyType) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixAddProxyType, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.ProxyHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProxyHash: %v", err)
	}
	return nil
}
func (v AddProxyType) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *DelProxyType) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixDelProxyType); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.ProxyHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProxyHash: %v", err)
	}
	return nil
}
func (v DelProxyType) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixDelProxyType, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.ProxyHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProxyHash: %v", err)
	}
	return nil
}
func (v DelProxyType) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *RegisterProxy) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixRegisterProxy); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if v.ProxyInfo, err = (func() (boc.Cell, error) {
		remain := c.CopyRemaining()
		if remain == nil {
			return boc.Cell{}, nil
		}
		return *remain, nil
	})(); err != nil {
		return fmt.Errorf("failed to read .ProxyInfo: %v", err)
	}
	return nil
}
func (v RegisterProxy) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixRegisterProxy, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = c.AddRef(&v.ProxyInfo); err != nil {
		return fmt.Errorf("failed to .ProxyInfo: %v", err)
	}
	return nil
}
func (v RegisterProxy) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *UnregisterProxy) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixUnregisterProxy); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Seqno.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Seqno: %v", err)
	}
	return nil
}
func (v UnregisterProxy) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixUnregisterProxy, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Seqno.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Seqno: %v", err)
	}
	return nil
}
func (v UnregisterProxy) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *UpdateProxy) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixUpdateProxy); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Seqno.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Seqno: %v", err)
	}
	if v.ProxyAddr, err = (func() (boc.Cell, error) {
		remain := c.CopyRemaining()
		if remain == nil {
			return boc.Cell{}, nil
		}
		return *remain, nil
	})(); err != nil {
		return fmt.Errorf("failed to read .ProxyAddr: %v", err)
	}
	return nil
}
func (v UpdateProxy) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixUpdateProxy, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Seqno.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Seqno: %v", err)
	}
	if err = c.AddRef(&v.ProxyAddr); err != nil {
		return fmt.Errorf("failed to .ProxyAddr: %v", err)
	}
	return nil
}
func (v UpdateProxy) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ChangeFees) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixChangeFees); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.PricePerToken.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PricePerToken: %v", err)
	}
	if err = v.WorkerFeePerToken.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .WorkerFeePerToken: %v", err)
	}
	return nil
}
func (v ChangeFees) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixChangeFees, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.PricePerToken.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PricePerToken: %v", err)
	}
	if err = v.WorkerFeePerToken.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .WorkerFeePerToken: %v", err)
	}
	return nil
}
func (v ChangeFees) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ChangeParams) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixChangeParams); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.PricePerToken.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PricePerToken: %v", err)
	}
	if err = v.WorkerFeePerToken.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .WorkerFeePerToken: %v", err)
	}
	if err = v.ProxyDelayBeforeClose.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProxyDelayBeforeClose: %v", err)
	}
	if err = v.ClientDelayBeforeClose.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ClientDelayBeforeClose: %v", err)
	}
	if err = v.MinProxyStake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinProxyStake: %v", err)
	}
	if err = v.MinClientStake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinClientStake: %v", err)
	}
	return nil
}
func (v ChangeParams) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixChangeParams, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.PricePerToken.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PricePerToken: %v", err)
	}
	if err = v.WorkerFeePerToken.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .WorkerFeePerToken: %v", err)
	}
	if err = v.ProxyDelayBeforeClose.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProxyDelayBeforeClose: %v", err)
	}
	if err = v.ClientDelayBeforeClose.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ClientDelayBeforeClose: %v", err)
	}
	if err = v.MinProxyStake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinProxyStake: %v", err)
	}
	if err = v.MinClientStake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinClientStake: %v", err)
	}
	return nil
}
func (v ChangeParams) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *UpgradeContracts) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixUpgradeContracts); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if v.ProxyCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .ProxyCode: %v", err)
	}
	if v.WorkerCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .WorkerCode: %v", err)
	}
	if v.ClientCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .ClientCode: %v", err)
	}
	return nil
}
func (v UpgradeContracts) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixUpgradeContracts, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = c.AddRef(&v.ProxyCode); err != nil {
		return fmt.Errorf("failed to .ProxyCode: %v", err)
	}
	if err = c.AddRef(&v.WorkerCode); err != nil {
		return fmt.Errorf("failed to .WorkerCode: %v", err)
	}
	if err = c.AddRef(&v.ClientCode); err != nil {
		return fmt.Errorf("failed to .ClientCode: %v", err)
	}
	return nil
}
func (v UpgradeContracts) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *UpgradeCode) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixUpgradeCode); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if v.NewCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .NewCode: %v", err)
	}
	return nil
}
func (v UpgradeCode) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixUpgradeCode, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = c.AddRef(&v.NewCode); err != nil {
		return fmt.Errorf("failed to .NewCode: %v", err)
	}
	return nil
}
func (v UpgradeCode) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ResetRoot) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixResetRoot); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}
func (v ResetRoot) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixResetRoot, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}
func (v ResetRoot) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *UpgradeFull) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixUpgradeFull); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if v.NewData, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .NewData: %v", err)
	}
	if v.NewCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .NewCode: %v", err)
	}
	return nil
}
func (v UpgradeFull) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixUpgradeFull, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = c.AddRef(&v.NewData); err != nil {
		return fmt.Errorf("failed to .NewData: %v", err)
	}
	if err = c.AddRef(&v.NewCode); err != nil {
		return fmt.Errorf("failed to .NewCode: %v", err)
	}
	return nil
}
func (v UpgradeFull) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ChangeOwner) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixChangeOwner); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.NewOwnerAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .NewOwnerAddress: %v", err)
	}
	return nil
}
func (v ChangeOwner) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixChangeOwner, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.NewOwnerAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewOwnerAddress: %v", err)
	}
	return nil
}
func (v ChangeOwner) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *RootData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.ProxyHashes.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProxyHashes: %v", err)
	}
	if err = v.RegisteredProxies.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RegisteredProxies: %v", err)
	}
	if err = v.LastProxySeqno.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .LastProxySeqno: %v", err)
	}
	if err = v.WorkerHashes.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .WorkerHashes: %v", err)
	}
	if err = v.ModelHashes.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ModelHashes: %v", err)
	}
	return nil
}
func (v RootData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.ProxyHashes.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProxyHashes: %v", err)
	}
	if err = v.RegisteredProxies.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RegisteredProxies: %v", err)
	}
	if err = v.LastProxySeqno.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .LastProxySeqno: %v", err)
	}
	if err = v.WorkerHashes.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .WorkerHashes: %v", err)
	}
	if err = v.ModelHashes.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ModelHashes: %v", err)
	}
	return nil
}
func (v RootData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *RootStorage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.OwnerAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .OwnerAddress: %v", err)
	}
	if err = v.Data.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Data: %v", err)
	}
	if err = v.Params.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Params: %v", err)
	}
	if err = v.Version.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Version: %v", err)
	}
	return nil
}
func (v RootStorage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.OwnerAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .OwnerAddress: %v", err)
	}
	if err = v.Data.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Data: %v", err)
	}
	if err = v.Params.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Params: %v", err)
	}
	if err = v.Version.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Version: %v", err)
	}
	return nil
}
func (v RootStorage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CocoonData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Version.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Version: %v", err)
	}
	if err = v.LastProxySeqno.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .LastProxySeqno: %v", err)
	}
	if err = v.ParamsVersion.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ParamsVersion: %v", err)
	}
	if err = v.UniqueId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UniqueId: %v", err)
	}
	if v.IsTest, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .IsTest: %v", err)
	}
	if err = v.PricePerToken.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PricePerToken: %v", err)
	}
	if err = v.WorkerFeePerToken.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .WorkerFeePerToken: %v", err)
	}
	if err = v.MinProxyStake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinProxyStake: %v", err)
	}
	if err = v.MinClientStake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinClientStake: %v", err)
	}
	if err = v.OwnerAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .OwnerAddress: %v", err)
	}
	return nil
}
func (v CocoonData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Version.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Version: %v", err)
	}
	if err = v.LastProxySeqno.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .LastProxySeqno: %v", err)
	}
	if err = v.ParamsVersion.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ParamsVersion: %v", err)
	}
	if err = v.UniqueId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UniqueId: %v", err)
	}
	if err = c.WriteBit(v.IsTest); err != nil {
		return fmt.Errorf("failed to .IsTest: %v", err)
	}
	if err = v.PricePerToken.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PricePerToken: %v", err)
	}
	if err = v.WorkerFeePerToken.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .WorkerFeePerToken: %v", err)
	}
	if err = v.MinProxyStake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinProxyStake: %v", err)
	}
	if err = v.MinClientStake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinClientStake: %v", err)
	}
	if err = v.OwnerAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .OwnerAddress: %v", err)
	}
	return nil
}
func (v CocoonData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CocoonData) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.OwnerAddress.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .OwnerAddress: %v", err)
	}
	if err = v.MinClientStake.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .MinClientStake: %v", err)
	}
	if err = v.MinProxyStake.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .MinProxyStake: %v", err)
	}
	if err = v.WorkerFeePerToken.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .WorkerFeePerToken: %v", err)
	}
	if err = v.PricePerToken.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .PricePerToken: %v", err)
	}
	if v.IsTest, err = stack.ReadBool(); err != nil {
		return fmt.Errorf("failed to read .IsTest: %v", err)
	}
	if err = v.UniqueId.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .UniqueId: %v", err)
	}
	if err = v.ParamsVersion.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ParamsVersion: %v", err)
	}
	if err = v.LastProxySeqno.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .LastProxySeqno: %v", err)
	}
	if err = v.Version.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Version: %v", err)
	}
	return nil
}
func (v *CurrentCocoonParams) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.ParamsVersion.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ParamsVersion: %v", err)
	}
	if err = v.UniqueId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UniqueId: %v", err)
	}
	if v.IsTest, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .IsTest: %v", err)
	}
	if err = v.PricePerToken.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PricePerToken: %v", err)
	}
	if err = v.WorkerFeePerToken.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .WorkerFeePerToken: %v", err)
	}
	if err = v.CachedTokensPriceMultiplier.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .CachedTokensPriceMultiplier: %v", err)
	}
	if err = v.ReasoningTokensPriceMultiplier.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ReasoningTokensPriceMultiplier: %v", err)
	}
	if err = v.ProxyDelayBeforeClose.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProxyDelayBeforeClose: %v", err)
	}
	if err = v.ClientDelayBeforeClose.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ClientDelayBeforeClose: %v", err)
	}
	if err = v.MinProxyStake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinProxyStake: %v", err)
	}
	if err = v.MinClientStake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinClientStake: %v", err)
	}
	if err = v.ProxyCodeHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProxyCodeHash: %v", err)
	}
	if err = v.WorkerCodeHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .WorkerCodeHash: %v", err)
	}
	if err = v.ClientCodeHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ClientCodeHash: %v", err)
	}
	return nil
}
func (v CurrentCocoonParams) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.ParamsVersion.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ParamsVersion: %v", err)
	}
	if err = v.UniqueId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UniqueId: %v", err)
	}
	if err = c.WriteBit(v.IsTest); err != nil {
		return fmt.Errorf("failed to .IsTest: %v", err)
	}
	if err = v.PricePerToken.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PricePerToken: %v", err)
	}
	if err = v.WorkerFeePerToken.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .WorkerFeePerToken: %v", err)
	}
	if err = v.CachedTokensPriceMultiplier.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .CachedTokensPriceMultiplier: %v", err)
	}
	if err = v.ReasoningTokensPriceMultiplier.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ReasoningTokensPriceMultiplier: %v", err)
	}
	if err = v.ProxyDelayBeforeClose.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProxyDelayBeforeClose: %v", err)
	}
	if err = v.ClientDelayBeforeClose.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ClientDelayBeforeClose: %v", err)
	}
	if err = v.MinProxyStake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinProxyStake: %v", err)
	}
	if err = v.MinClientStake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinClientStake: %v", err)
	}
	if err = v.ProxyCodeHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProxyCodeHash: %v", err)
	}
	if err = v.WorkerCodeHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .WorkerCodeHash: %v", err)
	}
	if err = v.ClientCodeHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ClientCodeHash: %v", err)
	}
	return nil
}
func (v CurrentCocoonParams) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CurrentCocoonParams) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.ClientCodeHash.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ClientCodeHash: %v", err)
	}
	if err = v.WorkerCodeHash.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .WorkerCodeHash: %v", err)
	}
	if err = v.ProxyCodeHash.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ProxyCodeHash: %v", err)
	}
	if err = v.MinClientStake.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .MinClientStake: %v", err)
	}
	if err = v.MinProxyStake.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .MinProxyStake: %v", err)
	}
	if err = v.ClientDelayBeforeClose.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ClientDelayBeforeClose: %v", err)
	}
	if err = v.ProxyDelayBeforeClose.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ProxyDelayBeforeClose: %v", err)
	}
	if err = v.ReasoningTokensPriceMultiplier.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ReasoningTokensPriceMultiplier: %v", err)
	}
	if err = v.CachedTokensPriceMultiplier.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .CachedTokensPriceMultiplier: %v", err)
	}
	if err = v.WorkerFeePerToken.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .WorkerFeePerToken: %v", err)
	}
	if err = v.PricePerToken.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .PricePerToken: %v", err)
	}
	if v.IsTest, err = stack.ReadBool(); err != nil {
		return fmt.Errorf("failed to read .IsTest: %v", err)
	}
	if err = v.UniqueId.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .UniqueId: %v", err)
	}
	if err = v.ParamsVersion.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ParamsVersion: %v", err)
	}
	return nil
}

func (msg AddWorkerType) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*RootStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg DelWorkerType) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*RootStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg AddModelType) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*RootStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg DelModelType) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*RootStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg AddProxyType) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*RootStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg DelProxyType) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*RootStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg RegisterProxy) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*RootStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg UnregisterProxy) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*RootStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg UpdateProxy) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*RootStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg ChangeFees) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*RootStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg ChangeParams) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*RootStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg UpgradeContracts) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*RootStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg UpgradeCode) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*RootStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg ResetRoot) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*RootStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg UpgradeFull) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*RootStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg ChangeOwner) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*RootStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}
