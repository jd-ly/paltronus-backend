package controllers

import (
	"github.com/revel/revel"
	"paltronus-backend/app/models"
	"paltronus-backend/app/services"
	"time"
)

type Websocket struct {
	*revel.Controller
}

func (c *Websocket) Init(id int, user string, ws revel.ServerWebSocket) revel.Result {
	// Make sure the websocket is valid.
	if ws == nil {
		return nil
	}

	// Join the room.
	subscription := services.Subscribe(id, user)
	defer subscription.Cancel()

	services.Join(user, id)
	defer services.Leave(user, id)

	newMessages := make(chan string)
	go func() {
		var msg string
		for {
			err := ws.MessageReceive(&msg)
			if err == nil {
				version := models.Version{
					RawData:      msg,
					CreatedBy:    user,
					CreationDate: time.Now().Format("2006-01-02 15:04:05"),
					File:         id,
				}
				_, e := services.InsertVersion(version)
				if e == nil {
					newMessages <- msg
					return
				}
			}
			close(newMessages)
			return
		}
	}()

	// Listen for new events from either the websocket or the services.
	for {
		select {
		case event := <-subscription.New:
			if ws.MessageSendJSON(&event) != nil {
				// They disconnected.
				return nil
			}
			if event.Type == "join" {

			}
		case msg, ok := <-newMessages:
			// If the channel is closed, they disconnected.
			if !ok {
				return nil
			}

			// Otherwise, say something.
			services.Register(user, msg, id)
		}
	}
	return nil
}