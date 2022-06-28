package validator

import (
	"encoding/base64"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func Validate(typeName string, msgBytes []byte) error {
	msgType, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(typeName))
	if err != nil {
		return &UnsupportedMessageType{MessageType: typeName, Internal: err}
	}

	var msg = msgType.New().Interface()
	err = proto.Unmarshal(msgBytes, msg)
	if err != nil {
		return &InvalidMessage{
			MessageType: typeName,
			Internal:    fmt.Errorf("%s: %s", err, base64.StdEncoding.EncodeToString(msgBytes)),
		}
	}

	// inspired by https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/validator/validator.go
	switch v := msg.(type) {
	case validatorLegacy:
		if err := v.Validate(); err != nil {
			return &InvalidMessage{MessageType: typeName, Internal: err}
		}
	case validator:
		if err := v.Validate(false); err != nil {
			return &InvalidMessage{MessageType: typeName, Internal: err}
		}
	}

	return nil
}

// The validate interface starting with protoc-gen-validate v0.6.0.
// See https://github.com/envoyproxy/protoc-gen-validate/pull/455.
type validator interface {
	Validate(all bool) error
}

// The validate interface prior to protoc-gen-validate v0.6.0.
type validatorLegacy interface {
	Validate() error
}
