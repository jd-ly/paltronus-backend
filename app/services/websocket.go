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

var (
	// Send a channel here to get room events back.  It will send the entire
	// archive initially, and then new messages as they come in.
	subscribe = make(chan Channel)
	// Send a channel here to unsubscribe.
	unsubscribe = make(chan (<-chan Event))
	// Send events here to publish them.
	publish = make(chan Event)
)

func getUsers(subs *list.List) string {
	var users = ""
	for ch := subs.Front(); ch != nil; ch = ch.Next() {
		value := ch.Value.(Subscriber).User
		if users == "" {
			users = value
		} else {
			users = users + "," + value
		}
	}
	return users
}

// This function loops forever
func room() {
	subscribers := make(map[int]*list.List)

	for {
		select {
		case ch := <-subscribe:
			subscriber := make(chan Event)
			if subscribers[ch.FileId] == nil {
				subscribers[ch.FileId] = list.New()
			}
			subscribers[ch.FileId].PushBack(Subscriber{ch.User, subscriber})
			ch.Subscription <- Subscription{0, subscriber}

		case event := <-publish:
			fileSubs := subscribers[event.FileId]
			if fileSubs != nil {
				currentUsers := getUsers(fileSubs)
				for ch := fileSubs.Front(); ch != nil; ch = ch.Next() {
					value := ch.Value.(Subscriber)
					if event.Type == "join" {
						if value.User == event.User {
							// Send content
							version, err := QueryLastVersion(event.FileId)
							if err == nil {
								value.Events <- newEvent("message", event.User, version.RawData, event.FileId)
							}
						}
						if currentUsers != "" {
							value.Events <- newEvent("join", event.User, currentUsers, event.FileId)
						}
					} else if event.Type == "leave" {
						// Remove the users
						for cha := fileSubs.Front(); cha != nil; cha = cha.Next() {
							sub := cha.Value.(Subscriber)
							if sub.User == event.User {
								fileSubs.Remove(cha)
								break
							}
						}
						currentUsers := getUsers(fileSubs)
						if currentUsers != "" {
							for cha := fileSubs.Front(); cha != nil; cha = cha.Next() {
								cha.Value.(Subscriber).Events <- newEvent("leave", event.User, currentUsers, event.FileId)
							}
						}
					} else if value.User != event.User {
							value.Events <- event
					}
				}
			}


		case unsub := <-unsubscribe:
			for _, value := range subscribers {
				fileSubs := value
				for ch := fileSubs.Front(); ch != nil; ch = ch.Next() {
					value := ch.Value.(Subscriber)
					if value.Events == unsub {
						fileSubs.Remove(ch)
						break
					}
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