package acceptor

import (
	"fmt"
	"github.com/mburman/hooli/rpc/acceptorrpc"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

var LOGE = log.New(os.Stderr, "ERROR ", log.Lmicroseconds|log.Lshortfile)
var LOGV = log.New(ioutil.Discard, "VERBOSE ", log.Lmicroseconds|log.Lshortfile)

type acceptorObj struct {
}

// port: port to start the acceptorObj on.
func NewAcceptor(port int) *acceptorObj {
	var a acceptorObj
	setupRPC(&a, port)
	return &a
}

func (a *acceptorObj) Prepare(args *acceptorrpc.PrepareArgs, reply *acceptorrpc.PrepareReply) error {
	return nil
}

func (a *acceptorObj) Accept(args *acceptorrpc.AcceptArgs, reply *acceptorrpc.AcceptReply) error {
	return nil
}

func (a *acceptorObj) Commit(args *acceptorrpc.CommitArgs, reply *acceptorrpc.CommitReply) error {
	return nil
}

func setupRPC(a *acceptorObj, port int) {
	rpc.RegisterName("AcceptorObj", a)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if e != nil {
		LOGE.Println("listen error:", e)
	}
	go http.Serve(l, nil)
}
