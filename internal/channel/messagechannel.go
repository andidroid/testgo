package channel

type MessageChannelEntity struct {
	Id      int
	Message string
}

func NewMessageChannelEntity(someParameter string) *MessageChannelEntity {
    p := new(MessageChannelEntity)
    p.Message = someParameter
    p.Id = 1 // <- a very sensible default value
    return p
}
