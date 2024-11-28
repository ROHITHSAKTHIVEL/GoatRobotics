package models

import "time"

type Chat struct {
	ClientID string    `json:"clientID,omitempty"`
	Message  string    `json:"message,omitempty"`
	UserName string    `json:"userName,omitempty"`
	SentTime time.Time `json:"sentTime,omitempty"`
}

type Request struct {
	ClientID string `json:"clientID,omitempty"`
}
