package errors

type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}


func (e *Error) Error() string {
    return e.Message
}

var (
	// Client-related errors
	ClientIDRequired   = &Error{Code: "CLIENT_ID_REQUIRED", Message: "Client ID is required."}
	UserNotFound       = &Error{Code: "USER_NOT_FOUND", Message: "User not found."}
	ClientAlreadyExist = &Error{Code: "CLIENT_ALREADY_EXIST", Message: "Client Already exist in the room."}

	// Message-related errors
	MessageRequired = &Error{Code: "MESSAGE_REQUIRED", Message: "Message content is required."}
	NoMessagesFound = &Error{Code: "NO_MESSAGES_FOUND", Message: "No new messages available."}

	// General errors
	RequestTimeout      = &Error{Code: "REQUEST_TIMEOUT", Message: "Request timed out."}
	InternalServerError = &Error{Code: "INTERNAL_SERVER_ERROR", Message: "An unexpected error occurred."}
)
