package errors

type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

var (
	Clinet_ID_REQUIRED = &Error{Code: "Client_ID_REQUIRED", Message: "Client ID is Required "}
)
