package main

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Text string `json:"text,omitempty"`
	Type string `json:"type,omitempty"`
}

type User struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Country string `json:"country,omitempty"`
}

type ViberCallback struct {
	Event string `json:"event,omitempty"`
	User  User   `json:"user,omitempty"`

	Message      Message `json:"message,omitempty"`
	Context      string  `json:"context"`
	MessageToken int     `json:"message_token,omitempty"`
}

type ViberCallbackMessage struct {
	User User `json:"sender,omitempty"`
}

type ViberSeenMessage struct {
	UserId string `json:"user_id,omitempty"`
}

func parseCallback(b []byte) (*ViberCallback, error) {
	ret := &ViberCallback{}
	err := json.Unmarshal(b, ret)
	if err != nil {
		return nil, fmt.Errorf("Invalid json: %v", err)
	}
	if ret.Event == "subscribed" || ret.Event == "conversation_started" {
		return ret, nil
	}
	if ret.Event == "message" {
		m := &ViberCallbackMessage{}
		err = json.Unmarshal(b, m)
		if err != nil {
			return nil, err
		}
		ret.User = m.User
		return ret, err
	}
	if ret.Event == "delivered" || ret.Event == "seen" || ret.Event == "unsubscribed" {
		m := &ViberSeenMessage{}
		err = json.Unmarshal(b, m)
		if err != nil {
			return nil, err
		}
		ret.User.Id = m.UserId
		return ret, err
	}

	return ret, err
}
