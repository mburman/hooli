package acceptor

import (
	"fmt"
	"github.com/mburman/hooli/rpc/acceptrpc"
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

func (a *acceptorObj) Prepare(args *acceptrpc.PrepareArgs, reply *acceptrpc.PrepareReply) {

}

func (a *acceptorObj) Accept(args *acceptrpc.AcceptArgs, reply *acceptrpc.AcceptReply) {

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
