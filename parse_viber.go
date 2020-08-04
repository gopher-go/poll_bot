package poll_bot

import (
	"encoding/json"
	"fmt"
)

// Message - Viber Message
type Message struct {
	Text string `json:"text,omitempty"`
	Type string `json:"type,omitempty"`
}

// User - Viber User
type User struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Country  string `json:"country,omitempty"`
	Language string `json:"language"`
	MNC      int    `json:"mnc"`
	MCC      int    `json:"mcc"`
}

// ViberCallback - Viber Callback
type ViberCallback struct {
	Event string `json:"event,omitempty"`
	User  User   `json:"user,omitempty"`

	Message      Message `json:"message,omitempty"`
	Context      string  `json:"context"`
	MessageToken int     `json:"message_token,omitempty"`
}

// ViberCallbackMessage - Viber Callback Message
type ViberCallbackMessage struct {
	User User `json:"sender,omitempty"`
}

// ViberSeenMessage - Viber Seen Message
type ViberSeenMessage struct {
	UserID string `json:"user_id,omitempty"`
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
		ret.User.ID = m.UserID
		return ret, err
	}

	return ret, err
}
