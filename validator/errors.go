package validator

import "fmt"

type UnsupportedMessageType struct {
	MessageType string
	Internal    error
}

func (e *UnsupportedMessageType) Error() string {
	return fmt.Sprintf("unsupported message type %q: %s", e.MessageType, e.Internal)
}

type InvalidMessage struct {
	MessageType string
	Internal    error
}

func (e *InvalidMessage) Error() string {
	return fmt.Sprintf("invalid message %q: %s", e.MessageType, e.Internal)
}
