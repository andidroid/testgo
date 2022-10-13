package routing

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/andidroid/testgo/pkg/redis"
	"github.com/andidroid/testgo/pkg/util"
)

type MongoEventType int

const (
	MONGO_EVENT_CREATED = iota + 1
	MONGO_EVENT_READ
	MONGO_EVENT_UPDATED
	MONGO_EVENT_DELETED
)

type MongoEvent struct {
	ID         string         `json:"id" bson:"_id"`
	Collection string         `json:"collection" bson:"collection"`
	Type       MongoEventType `json:"type" bson:"type"`
	Time       time.Time      `json:"time" bson:"time"`
	//Object {} `json:"object" bson:"object"`
}

func init() {
	go func() {
		ctx := context.Background()
		sub := redis.GetClient().Subscribe(ctx, "mongoevent")
		pe := MongoEvent{}
		for {
			msg, err := sub.ReceiveMessage(ctx)
			util.CheckErr(err)
			fmt.Println(msg)

			if err := json.Unmarshal([]byte(msg.Payload), &pe); err != nil {
				util.CheckErr(err)
			} else {
				fmt.Println("received mongo event: ", pe)
			}

		}
	}()
}

func SendMongoEvent(id string, collection string, mongoEventType MongoEventType) {
	pe := MongoEvent{
		ID:         id,
		Collection: collection,
		Type:       mongoEventType,
		Time:       time.Now(),
	}
	payload, err := json.Marshal(pe)
	util.CheckErr(err)
	ret := redis.GetClient().Publish(context.Background(), "mongoevent", payload)
	fmt.Println("publish mongo event:", ret)
}
