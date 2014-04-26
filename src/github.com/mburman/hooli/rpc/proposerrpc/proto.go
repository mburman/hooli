package proposerrpc


type Message struct {
	Latitude  float32
	Longitude float32
	MessageText   string
	Author    string
	Timestamp string
}


type GetMessagesArgs struct {
	Latitude  float64
	Longitude float64
	Radius    float64
}

type GetMessagesReply struct {
	Messages []Message
}

type PostMessageArgs struct {
	Message Message
}

type PostMessageReply struct {
}
