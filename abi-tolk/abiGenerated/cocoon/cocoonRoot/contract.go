// Code generated - DO NOT EDIT.

package abiGenerated

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

const PrefixAddWorkerType uint64 = 0xe34b1c60

type AddWorkerType struct {
	QueryId    tlb.Uint64
	WorkerHash tlb.Uint256
}

const PrefixDelWorkerType uint64 = 0x8d94a79a

type DelWorkerType struct {
	QueryId    tlb.Uint64
	WorkerHash tlb.Uint256
}

const PrefixAddModelType uint64 = 0xc146134d

type AddModelType struct {
	QueryId   tlb.Uint64
	ModelHash tlb.Uint256
}

const PrefixDelModelType uint64 = 0x92b11c18

type DelModelType struct {
	QueryId   tlb.Uint64
	ModelHash tlb.Uint256
}

const PrefixAddProxyType uint64 = 0x71860e80

type AddProxyType struct {
	QueryId   tlb.Uint64
	ProxyHash tlb.Uint256
}

const PrefixDelProxyType uint64 = 0x3c41d0b2

type DelProxyType struct {
	QueryId   tlb.Uint64
	ProxyHash tlb.Uint256
}

const PrefixRegisterProxy uint64 = 0x927c7cb5

type RegisterProxy struct {
	QueryId   tlb.Uint64
	ProxyInfo tlb.Any
}

const PrefixUnregisterProxy uint64 = 0x6d49eaf2

type UnregisterProxy struct {
	QueryId tlb.Uint64
	Seqno   tlb.Uint32
}

const PrefixUpdateProxy uint64 = 0x9c7924ba

type UpdateProxy struct {
	QueryId   tlb.Uint64
	Seqno     tlb.Uint32
	ProxyAddr tlb.Any
}

const PrefixChangeFees uint64 = 0xc52ed8d4

type ChangeFees struct {
	QueryId           tlb.Uint64
	PricePerToken     tlb.Coins
	WorkerFeePerToken tlb.Coins
}

const PrefixChangeParams uint64 = 0x022fa189

type ChangeParams struct {
	QueryId                tlb.Uint64
	PricePerToken          tlb.Coins
	WorkerFeePerToken      tlb.Coins
	ProxyDelayBeforeClose  tlb.Uint32
	ClientDelayBeforeClose tlb.Uint32
	MinProxyStake          tlb.Coins
	MinClientStake         tlb.Coins
}

const PrefixUpgradeContracts uint64 = 0xa2370f61

type UpgradeContracts struct {
	QueryId    tlb.Uint64
	ProxyCode  boc.Cell
	WorkerCode boc.Cell
	ClientCode boc.Cell
}

const PrefixUpgradeCode uint64 = 0x11aefd51

type UpgradeCode struct {
	QueryId tlb.Uint64
	NewCode boc.Cell
}

const PrefixResetRoot uint64 = 0x563c1d96

type ResetRoot struct {
	QueryId tlb.Uint64
}

const PrefixUpgradeFull uint64 = 0x4f7c5789

type UpgradeFull struct {
	QueryId tlb.Uint64
	NewData boc.Cell
	NewCode boc.Cell
}

const PrefixChangeOwner uint64 = 0xc4a1ae54

type ChangeOwner struct {
	QueryId         tlb.Uint64
	NewOwnerAddress tlb.InternalAddress
}
type RootMessageKind uint

const (
	RootMessageKind_AddWorkerType    RootMessageKind = 0xe34b1c60
	RootMessageKind_DelWorkerType    RootMessageKind = 0x8d94a79a
	RootMessageKind_AddModelType     RootMessageKind = 0xc146134d
	RootMessageKind_DelModelType     RootMessageKind = 0x92b11c18
	RootMessageKind_AddProxyType     RootMessageKind = 0x71860e80
	RootMessageKind_DelProxyType     RootMessageKind = 0x3c41d0b2
	RootMessageKind_RegisterProxy    RootMessageKind = 0x927c7cb5
	RootMessageKind_UnregisterProxy  RootMessageKind = 0x6d49eaf2
	RootMessageKind_UpdateProxy      RootMessageKind = 0x9c7924ba
	RootMessageKind_ChangeFees       RootMessageKind = 0xc52ed8d4
	RootMessageKind_ChangeParams     RootMessageKind = 0x022fa189
	RootMessageKind_UpgradeContracts RootMessageKind = 0xa2370f61
	RootMessageKind_UpgradeCode      RootMessageKind = 0x11aefd51
	RootMessageKind_ResetRoot        RootMessageKind = 0x563c1d96
	RootMessageKind_UpgradeFull      RootMessageKind = 0x4f7c5789
	RootMessageKind_ChangeOwner      RootMessageKind = 0xc4a1ae54
)

type RootMessage struct { // tagged union
	SumType          RootMessageKind
	AddWorkerType    *AddWorkerType
	DelWorkerType    *DelWorkerType
	AddModelType     *AddModelType
	DelModelType     *DelModelType
	AddProxyType     *AddProxyType
	DelProxyType     *DelProxyType
	RegisterProxy    *RegisterProxy
	UnregisterProxy  *UnregisterProxy
	UpdateProxy      *UpdateProxy
	ChangeFees       *ChangeFees
	ChangeParams     *ChangeParams
	UpgradeContracts *UpgradeContracts
	UpgradeCode      *UpgradeCode
	ResetRoot        *ResetRoot
	UpgradeFull      *UpgradeFull
	ChangeOwner      *ChangeOwner
}

func (v *RootMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = RootMessageKind(prefix)
	switch v.SumType {
	case RootMessageKind_AddWorkerType:
		v.AddWorkerType = new(AddWorkerType)
		return v.AddWorkerType.UnmarshalTLB(c, decoder)
	case RootMessageKind_DelWorkerType:
		v.DelWorkerType = new(DelWorkerType)
		return v.DelWorkerType.UnmarshalTLB(c, decoder)
	case RootMessageKind_AddModelType:
		v.AddModelType = new(AddModelType)
		return v.AddModelType.UnmarshalTLB(c, decoder)
	case RootMessageKind_DelModelType:
		v.DelModelType = new(DelModelType)
		return v.DelModelType.UnmarshalTLB(c, decoder)
	case RootMessageKind_AddProxyType:
		v.AddProxyType = new(AddProxyType)
		return v.AddProxyType.UnmarshalTLB(c, decoder)
	case RootMessageKind_DelProxyType:
		v.DelProxyType = new(DelProxyType)
		return v.DelProxyType.UnmarshalTLB(c, decoder)
	case RootMessageKind_RegisterProxy:
		v.RegisterProxy = new(RegisterProxy)
		return v.RegisterProxy.UnmarshalTLB(c, decoder)
	case RootMessageKind_UnregisterProxy:
		v.UnregisterProxy = new(UnregisterProxy)
		return v.UnregisterProxy.UnmarshalTLB(c, decoder)
	case RootMessageKind_UpdateProxy:
		v.UpdateProxy = new(UpdateProxy)
		return v.UpdateProxy.UnmarshalTLB(c, decoder)
	case RootMessageKind_ChangeFees:
		v.ChangeFees = new(ChangeFees)
		return v.ChangeFees.UnmarshalTLB(c, decoder)
	case RootMessageKind_ChangeParams:
		v.ChangeParams = new(ChangeParams)
		return v.ChangeParams.UnmarshalTLB(c, decoder)
	case RootMessageKind_UpgradeContracts:
		v.UpgradeContracts = new(UpgradeContracts)
		return v.UpgradeContracts.UnmarshalTLB(c, decoder)
	case RootMessageKind_UpgradeCode:
		v.UpgradeCode = new(UpgradeCode)
		return v.UpgradeCode.UnmarshalTLB(c, decoder)
	case RootMessageKind_ResetRoot:
		v.ResetRoot = new(ResetRoot)
		return v.ResetRoot.UnmarshalTLB(c, decoder)
	case RootMessageKind_UpgradeFull:
		v.UpgradeFull = new(UpgradeFull)
		return v.UpgradeFull.UnmarshalTLB(c, decoder)
	case RootMessageKind_ChangeOwner:
		v.ChangeOwner = new(ChangeOwner)
		return v.ChangeOwner.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}

const PrefixReturnExcessesBack uint64 = 0x2565934c

type ReturnExcessesBack struct {
	QueryId tlb.Uint64
}

const PrefixPayout uint64 = 0xc59a7cd3

type Payout struct {
	QueryId tlb.Uint64
}

func (v *AddWorkerType) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixAddWorkerType {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.WorkerHash.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *DelWorkerType) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixDelWorkerType {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.WorkerHash.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *AddModelType) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixAddModelType {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.ModelHash.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *DelModelType) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixDelModelType {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.ModelHash.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *AddProxyType) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixAddProxyType {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.ProxyHash.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *DelProxyType) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixDelProxyType {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.ProxyHash.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *RegisterProxy) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixRegisterProxy {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if v.ProxyInfo, err = (func() (tlb.Any, error) {
		cc := c.CopyRemaining()
		return tlb.Any(*cc), nil
	})(); err != nil {
		return err
	}
	return nil
}
func (v *UnregisterProxy) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixUnregisterProxy {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.Seqno.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *UpdateProxy) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixUpdateProxy {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.Seqno.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if v.ProxyAddr, err = (func() (tlb.Any, error) {
		cc := c.CopyRemaining()
		return tlb.Any(*cc), nil
	})(); err != nil {
		return err
	}
	return nil
}
func (v *ChangeFees) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixChangeFees {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.PricePerToken.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.WorkerFeePerToken.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *ChangeParams) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixChangeParams {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.PricePerToken.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.WorkerFeePerToken.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.ProxyDelayBeforeClose.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.ClientDelayBeforeClose.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.MinProxyStake.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.MinClientStake.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *UpgradeContracts) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixUpgradeContracts {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if v.ProxyCode, err = c.NextRefV(); err != nil {
		return err
	}
	if v.WorkerCode, err = c.NextRefV(); err != nil {
		return err
	}
	if v.ClientCode, err = c.NextRefV(); err != nil {
		return err
	}
	return nil
}
func (v *UpgradeCode) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixUpgradeCode {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if v.NewCode, err = c.NextRefV(); err != nil {
		return err
	}
	return nil
}
func (v *ResetRoot) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixResetRoot {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *UpgradeFull) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixUpgradeFull {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if v.NewData, err = c.NextRefV(); err != nil {
		return err
	}
	if v.NewCode, err = c.NextRefV(); err != nil {
		return err
	}
	return nil
}
func (v *ChangeOwner) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixChangeOwner {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.NewOwnerAddress.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *ReturnExcessesBack) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixReturnExcessesBack {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *Payout) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixPayout {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
