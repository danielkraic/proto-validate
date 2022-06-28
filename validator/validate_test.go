package validator_test

import (
	"testing"

	"github.com/danielkraic/proto-validate-plugin/example/person"
	"github.com/danielkraic/proto-validate-plugin/validator"
	"google.golang.org/protobuf/proto"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	p := &person.Person{
		Id:    9,
		Email: "felix@ca.ts",
		Name:  "Felix the Cat",
		Home: &person.Person_Location{
			Lat: 1,
			Lng: 1,
		},
	}
	pBytes, err := proto.Marshal(p)
	require.Nil(t, err)

	var (
		errUnsupportedType *validator.UnsupportedMessageType
		errInvalidMessage  *validator.InvalidMessage
		typeFullName       = p.ProtoReflect().Descriptor().FullName()
	)

	require.ErrorAs(t, validator.Validate("", nil), &errUnsupportedType)
	require.ErrorAs(t, validator.Validate("msg1", pBytes), &errUnsupportedType)
	require.ErrorAs(t, validator.Validate("msg1", []byte(`ola`)), &errUnsupportedType)
	require.ErrorAs(t, validator.Validate(string(typeFullName), []byte(`ola`)), &errInvalidMessage)

	require.Nil(t, validator.Validate(string(typeFullName), pBytes))
}

func TestInvalidMessage(t *testing.T) {
	p := &person.Person{}
	pBytes, err := proto.Marshal(p)
	require.Nil(t, err)

	var (
		errInvalidMessage *validator.InvalidMessage
		typeFullName      = p.ProtoReflect().Descriptor().FullName()
	)

	require.ErrorAs(t, validator.Validate(string(typeFullName), pBytes), &errInvalidMessage)
}
