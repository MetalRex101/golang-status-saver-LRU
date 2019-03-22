package keeper

import "gitlab.com/artilligence/http-db-saver/domain"

type client struct {
	send chan<- *domain.Message
}

func (c *client) Stop() {
	close(c.send)
}

func (c *client) Send(mess *domain.Message) {
	c.send <- mess
}


