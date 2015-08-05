package otr3

import "time"

// How long after sending a packet should we wait to send a heartbeat?
const heartbeatInterval = 60 * time.Second

type heartbeatContext struct {
	lastSent time.Time
}

func (c *Conversation) updateLastSent() {
	c.heartbeat.lastSent = time.Now()
}

func (c *Conversation) maybeHeartbeat(plain, toSend messageWithHeader, err error) ([]byte, messageWithHeader, messageWithHeader, error) {
	if err != nil {
		return nil, nil, nil, err
	}
	tsExtra, e := c.potentialHeartbeat(plain)
	return plain, toSend, tsExtra, e
}

func (c *Conversation) potentialHeartbeat(plain []byte) (toSend messageWithHeader, err error) {
	if plain != nil {
		now := time.Now()
		if c.heartbeat.lastSent.Before(now.Add(-heartbeatInterval)) {
			dataMsg, err := c.genDataMsgWithFlag(nil, messageFlagIgnoreUnreadable)
			if err != nil {
				return nil, err
			}
			toSend, err = c.wrapMessageHeader(msgTypeData, dataMsg.serialize())
			if err != nil {
				return nil, err
			}
			c.updateLastSent()
			messageEventHeartbeatSent(c)
		}
	}
	return
}
