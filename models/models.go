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

type MessagesResponse struct {
	Messages     []*Chat   `json:"messages,omitempty"`
	Message      string    `json:"message,omitempty"`
	ResponseTime time.Time `json:"ReponseTime,omitempty"`
	ID           string    `json:"userId,omitempty"`
}

type Logs struct {
	ClientID         string        `json:"clientID,omitempty"`
	Message          string        `json:"message,omitempty"`
	StartTime        time.Time     `json:"startTime,omitempty"`
	EndTime          time.Time     `json:"endTime,omitempty"`
	MethodName       string        `json:"methodName,omitempty"`
	RequestURL       string        `json:"requestURL,omitempty"`
	RequestBody      string        `json:"requestBody,omitempty"`
	RequestHeaders   string        `json:"requestHeaders,omitempty"`
	QueryParameters  string        `json:"queryParameters,omitempty"`
	ResponseDuration time.Duration `json:"responseDuration,omitempty"`
	ResponseBody     string        `json:"responseBody,omitempty"`
	ResponseHeaders  string        `json:"responseHeaders,omitempty"`
	StatusCode       int           `json:"statusCode,omitempty"`
}
