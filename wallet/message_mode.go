package wallet

type MessageMode int

// For detailed information about message modes take a look at https://docs.ton.org/develop/smart-contracts/messages.
const (
	// AttachAllRemainingBalance means that a wallet will transfer all the remaining balance to the destination
	// instead of the value originally indicated in the message.
	AttachAllRemainingBalance MessageMode = 128
	// AttachAllRemainingBalanceOfInboundMessage means that
	// a wallet will transfer all the remaining value of the inbound message in addition to the value initially indicated
	// in the new message
	AttachAllRemainingBalanceOfInboundMessage MessageMode = 64
	// DestroyAccount means that current account must be destroyed if its resulting balance is zero (often used with Mode 128).
	DestroyAccount MessageMode = 32
)

func IsMessageModeSet(modeValue int, mode MessageMode) bool {
	if modeValue&int(mode) == int(mode) {
		return true
	}
	return false
}
