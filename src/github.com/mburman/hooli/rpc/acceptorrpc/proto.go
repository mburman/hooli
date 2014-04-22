// Contains accept RPCs

package acceptorrpc

import (
	"github.com/mburman/hooli/rpc/proposerrpc"
)

type Status int

const (
	OK            Status = iota + 1 // The RPC was a success.
	PREV_ACCEPTED                   // Something else was previously accepted.
	CANCEL
)

type Proposal struct {
	Number int
	ID     int // unique server id to break clashing numbers
}

type PrepareArgs struct {
	Proposal Proposal
}

type PrepareReply struct {
	AcceptedProposalNumber int
	AcceptedMessage        proposerrpc.Message
	Status                 Status
}

type AcceptArgs struct {
	Proposal        Proposal
	ProposalMessage proposerrpc.Message
}

type AcceptReply struct {
	MinProposalNumber int
	Status            Status
}

type CommitArgs struct {
	Message proposerrpc.Message
}

type CommitReply struct {
}
