package domain

import "time"

type EntityType string

type Message struct {
	Name      string
	Type      EntityType
	Data      interface{}
	CreatedAt time.Time
}

type Repositories struct {
	Entity EntityRepository
}

type KeeperClient interface {
	Stop()
	Send(mess *Message)
}
