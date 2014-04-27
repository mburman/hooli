package proposer

import (
	"fmt"
	//"github.com/mburman/hooli/acceptor"
	"github.com/mburman/hooli/rpc/acceptorrpc"
	"github.com/mburman/hooli/rpc/proposerrpc"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"
	"code.google.com/p/gorest"
)

var LOGE = log.New(os.Stderr, "ERROR ", log.Lmicroseconds|log.Lshortfile)
var LOGV = log.New(ioutil.Discard, "VERBOSE ", log.Lmicroseconds|log.Lshortfile)

type proposerObj struct {
	port              int
	acceptorPorts     []string
	messageQueue      chan *proposerrpc.Message // messages to be handled
	messages          []proposerrpc.Message     // list of messages
	acceptorList      []*rpc.Client
	maxProposalNumber int // max proposal number server has seen
	id                int // proposer id

	//REST things
	gorest.RestService `root:"/proposer/" consumes:"application/json" produces:"application/json"`
	listMessages gorest.EndPoint `method:"GET" path:"/messages/" output:"[]Message"`
	postMessage gorest.EndPoint `method:"POST" path:"/messages/" postdata:"Message"`
}

// port: port for proposer to listen to client requests on.
// acceptorPorts: ports contact acceptors on.
var p proposerObj
func NewProposer(port int, acceptorPorts []string) *proposerObj {
//	var p proposerObj
	p.port = port
	p.acceptorPorts = acceptorPorts
	p.messageQueue = make(chan *proposerrpc.Message, 100)
	p.acceptorList = make([]*rpc.Client, 0)

	rand.Seed(time.Now().Unix())
	p.id = rand.Intn(100)               // Random server id.
	p.maxProposalNumber = rand.Intn(10) // Randomize initial round number.

	setupREST(&p, port)
//	setupRPC(&p,port)
	connectToAcceptors(&p)
	fmt.Println("done connecting to acceptors")
	go processMessages(&p) // start processing incoming messa
	return &p
}

type Message struct {
	Latitude  float64
	Longitude float64
	MessageText   string
	Author    string
}

//REST method handler for posting message
func(serv proposerObj) PostMessage(PostData Message){
	// Promise to handle. Don't block.
	fmt.Printf("Received Message to post:  %+v\n", PostData)
	rpcMess := new(proposerrpc.Message)
	rpcMess.MessageText = PostData.MessageText
	rpcMess.Author = PostData.Author
	rpcMess.Latitude = PostData.Latitude
	rpcMess.Longitude = PostData.Longitude
//	serv.ResponseBuilder().Created("http://localhost:9009/proposer/messages/"+string(m.author)) //Created, http 201
	serv.ResponseBuilder().Created("http://localhost:9009/proposer/messages/") //Created, http 201
	go handleMessage(&p, rpcMess)
}

func(serv proposerObj) ListMessages() []Message {
	fmt.Printf("Received request for messages\n");
	//get messages from acceptor
	acceptorArgs := &proposerrpc.GetMessagesArgs{
		Latitude: 0,
		Longitude: 0,
		Radius: 0,
	}
	acceptorReply := new(proposerrpc.GetMessagesReply)
	error := p.GetMessages(acceptorArgs,acceptorReply)
	if error != nil {
		LOGE.Println("error REST GETting messages")
		return nil
	}
	messArray := make([]Message,0)
	for _,v := range acceptorReply.Messages {
		mess := new(Message)
		mess.MessageText = v.MessageText
		mess.Latitude = v.Latitude
		mess.Longitude = v.Longitude
		mess.Author = v.Author
		fmt.Println("appending to messArray: ",mess)
		messArray = append(messArray, *mess)
	}
//	messArray = append(messArray,Message{Latitude:37.33233141,Longitude:-122.03121860,MessageText:"test1",Author:"Dylan Koenig"})
//	messArray = append(messArray,Message{Latitude:37.33233141,Longitude:-122.03121860,MessageText:"test2",Author:"Dylan Koenig"})
//	fmt.Println("messages:",messArray)
	return messArray
}

// Client calls this to post a message.

func (p *proposerObj) PostMessageRPC(args *proposerrpc.PostMessageArgs, reply *proposerrpc.PostMessageReply) error {
	// Promise to handle. Don't block.
	fmt.Printf("Received Message to post:  %+v\n", args.Message)
	go handleMessage(p, &args.Message)
	return nil
}


func (p *proposerObj) GetMessages(args *proposerrpc.GetMessagesArgs, reply *proposerrpc.GetMessagesReply) error {
	LOGV.Println("Getting messages")
	// Get messages from a storage server (acceptor)
	if len(p.acceptorList) < 1 {
		// panic return error
	}

	client := p.acceptorList[0]
	request := &acceptorrpc.GetMessagesArgs{
		Latitude:  args.Latitude,
		Longitude: args.Longitude,
		Radius:    args.Radius,
	}
	var messagesReply acceptorrpc.GetMessagesReply
	err := client.Call("AcceptorObj.GetMessages", request, &messagesReply)
	if err != nil {
		LOGE.Println("rpc error:", err)
	}
	reply.Messages = messagesReply.Messages
	return nil
}

func handleMessage(p *proposerObj, message *proposerrpc.Message) {
	fmt.Println("handling message:",*message)
	fmt.Printf("pobj: %+v\n", p)
	p.messageQueue <- message
	fmt.Println("handled message:",*message)
}

func connectToAcceptors(p *proposerObj) {
	// Connect to acceptors.
	for _, acceptorPort := range p.acceptorPorts {
		// Keep redialing till we connect to the server.
		client, err := rpc.DialHTTP("tcp", ":"+acceptorPort)
		if err != nil {
			LOGE.Println("redialing:", err)
			dialTicker := time.NewTicker(time.Second)
		DialLoop:
			for {
				select {
				case <-dialTicker.C:
					client, err = rpc.DialHTTP("tcp", ":"+acceptorPort)
					if err == nil {
						break DialLoop
					}
				}
			}
		}
		p.acceptorList = append(p.acceptorList, client)
		fmt.Println("acceptor connected")
	}
}

func generateProposal(p *proposerObj) *acceptorrpc.Proposal {
	p.maxProposalNumber++
	return &acceptorrpc.Proposal{
		Number: p.maxProposalNumber,
		ID:     p.id,
	}
}

func sendPrepare(p *proposerObj, client *rpc.Client, proposal *acceptorrpc.Proposal) acceptorrpc.PrepareReply {
	request := &acceptorrpc.PrepareArgs{
		Proposal: *proposal,
	}
	var reply acceptorrpc.PrepareReply
	err := client.Call("AcceptorObj.Prepare", request, &reply)
	if err != nil {
		LOGE.Println("rpc error:", err)
	}
	return reply
}

func sendAccept(p *proposerObj, client *rpc.Client, proposal *acceptorrpc.Proposal,
	acceptedMessage *proposerrpc.Message) acceptorrpc.AcceptReply {
	request := &acceptorrpc.AcceptArgs{
		Proposal:        *proposal,
		ProposalMessage: *acceptedMessage,
	}
	var reply acceptorrpc.AcceptReply
	err := client.Call("AcceptorObj.Accept", request, &reply)
	if err != nil {
		LOGE.Println("rpc error:", err)
	}
	return reply
}

func sendCommit(p *proposerObj, client *rpc.Client, message *proposerrpc.Message, index int) acceptorrpc.CommitReply {
	request := &acceptorrpc.CommitArgs{
		Message: *message,
		Index:   index,
	}
	var reply acceptorrpc.CommitReply
	err := client.Call("AcceptorObj.Commit", request, &reply)
	if err != nil {
		LOGE.Println("rpc error:", err)
	}
	return reply
}

// Continuously reads messages from queue and Paxos' them
func processMessages(pGlobal *proposerObj) {
	for {
		fmt.Println("waiting for message to process")
		message := <-pGlobal.messageQueue
		fmt.Println("processing message: ", *message)
		for {
			delayMills := rand.Intn(75) // Random delay to avoid conflicts.
			time.Sleep(time.Millisecond * time.Duration(delayMills))
			proposal := generateProposal(pGlobal)

			// Send prepares. TODO: needs to be done async since nodes might go down.
			acceptCount := 0
			highestCancelProposalNumber := -1
			highestIndex := 0
			for _, a := range pGlobal.acceptorList {
				// Send prepare message to all the acceptors
				prepareReply := sendPrepare(pGlobal, a, proposal)
				if prepareReply.Status == acceptorrpc.OK {
					if highestIndex < prepareReply.Index {
						highestIndex = prepareReply.Index
					}
					acceptCount++
				} else if prepareReply.Status == acceptorrpc.PREV_ACCEPTED {
					// IGNORE...
				} else if prepareReply.Status == acceptorrpc.CANCEL {
					if prepareReply.AcceptedProposalNumber >= highestCancelProposalNumber {
						highestCancelProposalNumber = prepareReply.AcceptedProposalNumber
					}
				}
			}

			// If a majority have not accepted - this is not the leader
			// Try again.
			if acceptCount <= len(pGlobal.acceptorList)/2 {
				pGlobal.maxProposalNumber = highestCancelProposalNumber
				continue
			}

			fmt.Println("Sending accepts")
			// LEADER. Send accepts.
			acceptCount = 0
			highestCancelProposalNumber = -1
			for _, a := range pGlobal.acceptorList {
				acceptReply := sendAccept(pGlobal, a, proposal, message)
				if acceptReply.Status == acceptorrpc.OK {
					acceptCount++
				} else if acceptReply.Status == acceptorrpc.CANCEL {
					if acceptReply.MinProposalNumber >= highestCancelProposalNumber {
						highestCancelProposalNumber = acceptReply.MinProposalNumber
					}
				} else {
					// ASSERT FALSE should never happen
				}
			}

			if acceptCount <= len(pGlobal.acceptorList)/2 {
				pGlobal.maxProposalNumber = highestCancelProposalNumber
				continue
			}

			fmt.Println("COMMITTING")
			// Value has been chosen.
			for _, a := range pGlobal.acceptorList {
				sendCommit(pGlobal, a, message, highestIndex)
			}

			break // get a new message.
		}
	}
}

func setupRPC(a *proposerObj, port int) {
	fmt.Println("Proposer rpc:", port)
	rpc.RegisterName("ProposerObj", a)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if e != nil {
		LOGE.Println("listen error:", e)
	}
	go http.Serve(l, nil)
}

func setupREST(a *proposerObj, port int) {
	fmt.Println("Proposer RESTing:", port)
	gorest.RegisterService(a)
	http.Handle("/",gorest.Handle())
	go http.ListenAndServe(fmt.Sprintf(":%d", port),nil)
}
