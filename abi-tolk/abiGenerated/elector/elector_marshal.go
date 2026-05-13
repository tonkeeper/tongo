// Code generated - DO NOT EDIT.

package abiElector

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *NewStake) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixNewStake); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.ValidatorPubkey.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ValidatorPubkey: %v", err)
	}
	if err = v.StakeAt.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .StakeAt: %v", err)
	}
	if err = v.MaxFactor.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MaxFactor: %v", err)
	}
	if err = v.AdnlAddr.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AdnlAddr: %v", err)
	}
	if err = v.Signature.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Signature: %v", err)
	}
	return nil
}
func (v NewStake) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixNewStake, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.ValidatorPubkey.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ValidatorPubkey: %v", err)
	}
	if err = v.StakeAt.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .StakeAt: %v", err)
	}
	if err = v.MaxFactor.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MaxFactor: %v", err)
	}
	if err = v.AdnlAddr.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AdnlAddr: %v", err)
	}
	if err = v.Signature.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Signature: %v", err)
	}
	return nil
}
func (v NewStake) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *NewStakeConfirmation) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixNewStakeConfirmation); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Comment.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Comment: %v", err)
	}
	return nil
}
func (v NewStakeConfirmation) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixNewStakeConfirmation, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Comment.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Comment: %v", err)
	}
	return nil
}
func (v NewStakeConfirmation) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *RecoverStakeRequest) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixRecoverStakeRequest); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}
func (v RecoverStakeRequest) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixRecoverStakeRequest, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}
func (v RecoverStakeRequest) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *RecoverStakeResponse) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixRecoverStakeResponse); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}
func (v RecoverStakeResponse) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixRecoverStakeResponse, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}
func (v RecoverStakeResponse) ToCell() (*boc.Cell, error) {
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
	if v.Code, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .Code: %v", err)
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
	if err = c.AddRef(&v.Code); err != nil {
		return fmt.Errorf("failed to .Code: %v", err)
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
func (v *UpgradeCodeResponse) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixUpgradeCodeResponse); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Op.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Op: %v", err)
	}
	return nil
}
func (v UpgradeCodeResponse) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixUpgradeCodeResponse, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Op.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Op: %v", err)
	}
	return nil
}
func (v UpgradeCodeResponse) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ConfigAccepted) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixConfigAccepted); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}
func (v ConfigAccepted) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixConfigAccepted, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}
func (v ConfigAccepted) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ConfigRejected) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixConfigRejected); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}
func (v ConfigRejected) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixConfigRejected, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}
func (v ConfigRejected) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *RegisterComplaint) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixRegisterComplaint); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.ElectionId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ElectionId: %v", err)
	}
	if err = v.Complaint.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Complaint: %v", err)
	}
	return nil
}
func (v RegisterComplaint) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixRegisterComplaint, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.ElectionId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ElectionId: %v", err)
	}
	if err = v.Complaint.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Complaint: %v", err)
	}
	return nil
}
func (v RegisterComplaint) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *VoteComplaint) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixVoteComplaint); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Signature.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Signature: %v", err)
	}
	if err = v.SignTag.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SignTag: %v", err)
	}
	if err = v.ValidatorIdx.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ValidatorIdx: %v", err)
	}
	if err = v.ElectionId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ElectionId: %v", err)
	}
	if err = v.ComplaintHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ComplaintHash: %v", err)
	}
	return nil
}
func (v VoteComplaint) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixVoteComplaint, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Signature.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Signature: %v", err)
	}
	if err = v.SignTag.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SignTag: %v", err)
	}
	if err = v.ValidatorIdx.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ValidatorIdx: %v", err)
	}
	if err = v.ElectionId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ElectionId: %v", err)
	}
	if err = v.ComplaintHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ComplaintHash: %v", err)
	}
	return nil
}
func (v VoteComplaint) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ReturnStake) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixReturnStake); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Reason.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Reason: %v", err)
	}
	return nil
}
func (v ReturnStake) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixReturnStake, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Reason.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Reason: %v", err)
	}
	return nil
}
func (v ReturnStake) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ErrorResponse) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixErrorResponse); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Op.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Op: %v", err)
	}
	return nil
}
func (v ErrorResponse) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixErrorResponse, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Op.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Op: %v", err)
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
func (v *ComplaintResponse) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixComplaintResponse); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Op.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Op: %v", err)
	}
	return nil
}
func (v ComplaintResponse) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixComplaintResponse, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Op.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Op: %v", err)
	}
	return nil
}
func (v ComplaintResponse) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *VoteComplaintResponse) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixVoteComplaintResponse); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Op.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Op: %v", err)
	}
	return nil
}
func (v VoteComplaintResponse) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixVoteComplaintResponse, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Op.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Op: %v", err)
	}
	return nil
}
func (v VoteComplaintResponse) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ValidatorComplaint) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Tag.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Tag: %v", err)
	}
	if err = v.ValidatorPubkey.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ValidatorPubkey: %v", err)
	}
	if v.Description, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .Description: %v", err)
	}
	if err = v.CreatedAt.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .CreatedAt: %v", err)
	}
	if err = v.Severity.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Severity: %v", err)
	}
	if err = v.RewardAddr.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RewardAddr: %v", err)
	}
	if err = v.Paid.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Paid: %v", err)
	}
	if err = v.SuggestedFine.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SuggestedFine: %v", err)
	}
	if err = v.SuggestedFinePart.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SuggestedFinePart: %v", err)
	}
	return nil
}
func (v ValidatorComplaint) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Tag.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Tag: %v", err)
	}
	if err = v.ValidatorPubkey.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ValidatorPubkey: %v", err)
	}
	if err = c.AddRef(&v.Description); err != nil {
		return fmt.Errorf("failed to .Description: %v", err)
	}
	if err = v.CreatedAt.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .CreatedAt: %v", err)
	}
	if err = v.Severity.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Severity: %v", err)
	}
	if err = v.RewardAddr.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RewardAddr: %v", err)
	}
	if err = v.Paid.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Paid: %v", err)
	}
	if err = v.SuggestedFine.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SuggestedFine: %v", err)
	}
	if err = v.SuggestedFinePart.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SuggestedFinePart: %v", err)
	}
	return nil
}
func (v ValidatorComplaint) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ElectorMember) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Stake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Stake: %v", err)
	}
	if err = v.Time.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Time: %v", err)
	}
	if err = v.MaxFactor.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MaxFactor: %v", err)
	}
	if err = v.SrcAddr.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SrcAddr: %v", err)
	}
	if err = v.AdnlAddr.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AdnlAddr: %v", err)
	}
	return nil
}
func (v ElectorMember) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Stake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Stake: %v", err)
	}
	if err = v.Time.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Time: %v", err)
	}
	if err = v.MaxFactor.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MaxFactor: %v", err)
	}
	if err = v.SrcAddr.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SrcAddr: %v", err)
	}
	if err = v.AdnlAddr.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AdnlAddr: %v", err)
	}
	return nil
}
func (v ElectorMember) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ParticipantListValidatorData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Stake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Stake: %v", err)
	}
	if err = v.MaxFactor.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MaxFactor: %v", err)
	}
	if err = v.Address.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Address: %v", err)
	}
	if err = v.AdnlAddr.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AdnlAddr: %v", err)
	}
	return nil
}
func (v ParticipantListValidatorData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Stake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Stake: %v", err)
	}
	if err = v.MaxFactor.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MaxFactor: %v", err)
	}
	if err = v.Address.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Address: %v", err)
	}
	if err = v.AdnlAddr.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AdnlAddr: %v", err)
	}
	return nil
}
func (v ParticipantListValidatorData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ParticipantListValidatorData) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.AdnlAddr.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .AdnlAddr: %v", err)
	}
	if err = v.Address.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Address: %v", err)
	}
	if err = v.MaxFactor.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .MaxFactor: %v", err)
	}
	if err = v.Stake.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Stake: %v", err)
	}
	return nil
}
func (v *ParticipantListExtended) ReadFromStack(stack *tlb.VmStack) (err error) {
	if v.Finished, err = stack.ReadBool(); err != nil {
		return fmt.Errorf("failed to read .Finished: %v", err)
	}
	if v.Failed, err = stack.ReadBool(); err != nil {
		return fmt.Errorf("failed to read .Failed: %v", err)
	}
	if v.Validators, err = tlb.ReadArrayFromStack[tlb.ShapedTuple2[tlb.Uint256, ParticipantListValidatorData]](stack, func(stack *tlb.VmStack) (value tlb.ShapedTuple2[tlb.Uint256, ParticipantListValidatorData], err error) {
		return (func() (result tlb.ShapedTuple2[tlb.Uint256, ParticipantListValidatorData], err error) {
			var tuple tlb.VmStkTuple
			if tuple, err = stack.ReadTuple(); err != nil {
				err = fmt.Errorf("read tensor: %w", err)
				return
			}
			var stack *tlb.VmStack
			stack, err = tuple.AsStack()
			if err != nil {
				err = fmt.Errorf("read tensor items: %w", err)
				return
			}

			if result.V1, err = tlb.ReadTupleFromStack(stack, func(stack *tlb.VmStack) (result ParticipantListValidatorData, err error) {
				err = result.ReadFromStack(stack)
				return
			}); err != nil {
				return
			}
			if err = result.V0.ReadFromStack(stack); err != nil {
				return
			}
			return
		})()
	}); err != nil {
		return fmt.Errorf("failed to read .Validators: %v", err)
	}
	if err = v.TotalStake.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .TotalStake: %v", err)
	}
	if err = v.MinStake.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .MinStake: %v", err)
	}
	if err = v.ElectClose.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ElectClose: %v", err)
	}
	if err = v.ElectAt.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ElectAt: %v", err)
	}
	return nil
}
func (v *Elect) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.ElectAt.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ElectAt: %v", err)
	}
	if err = v.ElectClose.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ElectClose: %v", err)
	}
	if err = v.MinStake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinStake: %v", err)
	}
	if err = v.TotalStake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TotalStake: %v", err)
	}
	if err = v.Members.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Members: %v", err)
	}
	if v.Failed, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .Failed: %v", err)
	}
	if v.Finished, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .Finished: %v", err)
	}
	return nil
}
func (v Elect) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.ElectAt.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ElectAt: %v", err)
	}
	if err = v.ElectClose.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ElectClose: %v", err)
	}
	if err = v.MinStake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinStake: %v", err)
	}
	if err = v.TotalStake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TotalStake: %v", err)
	}
	if err = v.Members.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Members: %v", err)
	}
	if err = c.WriteBit(v.Failed); err != nil {
		return fmt.Errorf("failed to .Failed: %v", err)
	}
	if err = c.WriteBit(v.Finished); err != nil {
		return fmt.Errorf("failed to .Finished: %v", err)
	}
	return nil
}
func (v Elect) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ElectorStorage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Elect.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Elect: %v", err)
	}
	if err = v.Credits.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Credits: %v", err)
	}
	if err = v.PastElections.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PastElections: %v", err)
	}
	if err = v.Grams.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Grams: %v", err)
	}
	if err = v.ActiveId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ActiveId: %v", err)
	}
	if err = v.ActiveHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ActiveHash: %v", err)
	}
	return nil
}
func (v ElectorStorage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Elect.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Elect: %v", err)
	}
	if err = v.Credits.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Credits: %v", err)
	}
	if err = v.PastElections.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PastElections: %v", err)
	}
	if err = v.Grams.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Grams: %v", err)
	}
	if err = v.ActiveId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ActiveId: %v", err)
	}
	if err = v.ActiveHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ActiveHash: %v", err)
	}
	return nil
}
func (v ElectorStorage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (msg NewStake) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ElectorStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg RecoverStakeRequest) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ElectorStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg UpgradeCode) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ElectorStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg ConfigAccepted) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ElectorStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg ConfigRejected) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ElectorStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg RegisterComplaint) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ElectorStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg VoteComplaint) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ElectorStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}
