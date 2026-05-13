tolk 1.4

// Tolk ABI for the TON elector contract (FunC source: governance-contract/elector-code.fc).
// This is a stub

contract elector {
    incomingMessages: ElectorMessage
    storage: ElectorStorage
    outgoingMessages:
        VoteComplaintResponse |
        ComplaintResponse |
        UpgradeCodeResponse |
        ErrorResponse |
        ReturnStake |
        RecoverStakeResponse |
        NewStakeConfirmation
}

type ElectorMessage =
    NewStake |
    RecoverStakeRequest |
    UpgradeCode |
    ConfigAccepted |
    ConfigRejected |
    RegisterComplaint |
    VoteComplaint;

struct (0x4e73744b) NewStake {
    queryId: uint64
    validatorPubkey: bits256
    stakeAt: uint32
    maxFactor: uint32
    adnlAddr: bits256
    signature: Cell<bits512>
}

struct (0xf374484c) NewStakeConfirmation {
    queryId: uint64
    comment: uint32
}

struct (0x47657424) RecoverStakeRequest {
    queryId: uint64
}

struct (0xf96f7324) RecoverStakeResponse {
    queryId: uint64
}

struct (0x4e436f64) UpgradeCode {
    queryId: uint64
    code: cell
}

struct (0xce436f64) UpgradeCodeResponse {
    queryId: uint64
    op: uint32
}

struct (0xee764f4b) ConfigAccepted {
    queryId: uint64
}

struct (0xee764f6f) ConfigRejected {
    queryId: uint64
}

struct (0x52674370) RegisterComplaint {
    queryId: uint64
    electionId: uint32
    complaint: ValidatorComplaint
}

struct (0x56744370) VoteComplaint {
    queryId: uint64
    signature: bits512
    signTag: uint32
    validatorIdx: uint16
    electionId: uint32
    complaintHash: uint256
}

struct (0xee6f454c) ReturnStake {
    queryId: uint64
    reason: uint32
}

struct (0xfffffffe) ErrorResponse {
    queryId: uint64
    op: uint32
}

struct (0xf2676350) ComplaintResponse {
    queryId: uint64
    op: uint32
}

struct (0xd6745240) VoteComplaintResponse {
    queryId: uint64
    op: uint32
}

struct ValidatorComplaint {
    tag: uint8
    validatorPubkey: uint256
    description: cell
    createdAt: uint32
    severity: uint8
    rewardAddr: uint256
    paid: coins
    suggestedFine: coins
    suggestedFinePart: uint32
}

struct ElectorMember {
    stake: coins
    time: uint32
    maxFactor: uint32
    srcAddr: uint256
    adnlAddr: uint256
}

struct ParticipantListValidatorData {
    stake: coins
    maxFactor: uint32
    address: uint256
    adnlAddr: uint256
}

struct ParticipantListExtended {
    electAt: uint32
    electClose: uint32
    minStake: coins
    totalStake: coins
    validators: array<[uint256, ParticipantListValidatorData]>
    failed: bool
    finished: bool
}

fun ParticipantListExtended.unpackFromSlice(mutate s: slice): ParticipantListExtended {
    return ParticipantListExtended {
        electAt: 0,
        electClose: 0,
        minStake: 0,
        totalStake: 0,
        validators: [],
        failed: false,
        finished: false,
    };
}

fun ParticipantListExtended.packToBuilder(self, mutate b: builder) {
}

struct Elect {
    electAt: uint32
    electClose: uint32
    minStake: coins
    totalStake: coins
    members: map<uint256, ElectorMember>
    failed: bool
    finished: bool
}

struct ElectorStorage {
    elect: Cell<Elect>?
    credits: map<uint256, coins>
    pastElections: map<uint32, slice>
    grams: coins
    activeId: uint32
    activeHash: uint256
}

fun onInternalMessage(in: InMessage) {
}

get fun active_election_id(): int {
    return 0;
}

get fun participates_in(validator_key: uint256): int {
    return 0;
}

get fun compute_returned_stake(wallet_addr: uint256): int {
    return 0;
}

get fun participant_list_extended(): ParticipantListExtended {
    throw 0;
}
