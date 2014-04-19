// Contains accept RPCs

package acceptrpc

import (
	"github.com/mburman/hooli/rpc/proposerrpc"
)

type PrepareArgs struct {
	ProposalNumber int
}

type PrepareReply struct {
	AcceptedProposalNumber int
	AcceptedMessage        proposerrpc.Message
}

type AcceptArgs struct {
	ProposalNumber  int
	ProposalMessage proposerrpc.Message
}

type AcceptReply struct {
	MinProposal int
}
