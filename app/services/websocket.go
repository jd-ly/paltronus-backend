package services

import (
	"container/list"
	"time"
)

type Channel struct {
	FileId 		 int
	User    	 string
	Subscription chan<- Subscription
}

type Subscriber struct {
	FileId 	int
	User    	 string
	Events 	chan Event
}

type Event struct {
	Type      string // "join", "leave", or "message"
	User      string
	Timestamp int    // Unix timestamp (secs)
	Text      string // What the user said (if Type == "message")
	FileId	  int
}

type Subscription struct {
	fileId 	int
	New     <-chan Event // New events coming in.
}

// Owner of a subscription must cancel it when they stop listening to events.
func (s Subscription) Cancel() {
	unsubscribe <- s.New // Unsubscribe the channel.
	drain(s.New)         // Drain it, just in case there was a pending publish.
}

func newEvent(typ, user, msg string, fileId int) Event {
	return Event{typ, user, int(time.Now().Unix()), msg, fileId}
}

func Subscribe(id int, user string) Subscription {
	resp := make(chan Subscription)
	subscribe <- Channel{id, user,resp}
	return <-resp
}

func Join(user string, id int) {
	publish <- newEvent("join", user, "", id)
}

func Register(user, message string, id int) {
	publish <- newEvent("message", user, message, id)
}

func Leave(user string, id int) {
	publish <- newEvent("leave", user, "", id)
}

const archiveSize = 10

var (
	// Send a channel here to get room events back.  It will send the entire
	// archive initially, and then new messages as they come in.
	subscribe = make(chan Channel)
	// Send a channel here to unsubscribe.
	unsubscribe = make(chan (<-chan Event))
	// Send events here to publish them.
	publish = make(chan Event)
)

func getUsers(subs *list.List, id int) string {
	var users = ""
	for ch := subs.Front(); ch != nil; ch = ch.Next() {
		value := ch.Value.(Subscriber)
		if value.FileId == id {
			if users == "" {
				users = value.User
			} else {
				users = users + "," + value.User
			}
		}
	}
	return users
}

// This function loops forever
func room() {
	subscribers := list.New()

	for {
		select {
		case ch := <-subscribe:
			subscriber := make(chan Event)
			subscribers.PushBack(Subscriber{ch.FileId, ch.User, subscriber})
			ch.Subscription <- Subscription{0, subscriber}

		case event := <-publish:
			for ch := subscribers.Front(); ch != nil; ch = ch.Next() {
				value := ch.Value.(Subscriber)
				if value.FileId == event.FileId {
					if event.Type == "join" {
						if value.User == event.User {
							// Send content
							version, err := QueryLastVersion(event.FileId)
							if err == nil {
								value.Events <- newEvent("message", event.User, version.RawData, event.FileId)
							}
						}
						currentUsers := getUsers(subscribers, event.FileId)
						if currentUsers != "" {
							value.Events <- newEvent("join", event.User, currentUsers, event.FileId)
						}
					} else if event.Type == "leave" {
						unsubscribe <- value.Events
						currentUsers := getUsers(subscribers, event.FileId)
						if currentUsers != "" {

						}
					} else {
						if value.User != event.User {
							// Send to everyone else
							value.Events <- event
						}
					}
				}
			}


		case unsub := <-unsubscribe:
			for ch := subscribers.Front(); ch != nil; ch = ch.Next() {
				value := ch.Value.(Subscriber)
				if value.Events == unsub {
					subscribers.Remove(ch)
					break
				}
			}
		}
	}
}

func init() {
	go room()
}

// Helpers

// Drains a given channel of any messages.
func drain(ch <-chan Event) {
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		default:
			return
		}
	}
}