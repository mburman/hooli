package proposerrpc

type Message struct {
	latitude  float64
	longitude float64
	message   string
	userid    string
	timestamp string
}

type GetMessagesArgs struct {
	latitude  float64
	longitude float64
	radius    float64
}

type GetMessagesReply struct {
	messages []Message
}

type PostMessageArgs struct {
	message Message
}

type PostMessageReply struct {
}
