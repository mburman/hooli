// Contains accept RPCs

package acceptrpc

import (
	"github.com/mburman/hooli/rpc/proposerrpc"
)

type PrepareArgs struct {
	proposalNumber int
}

type PrepareReply struct {
	acceptedProposalNumber int
	acceptedMessage        proposerrpc.Message
}

type AcceptArgs struct {
	proposalNumber  int
	proposalMessage proposerrpc.Message
}

type AcceptReply struct {
	minProposal int
}
