package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/danielkraic/proto-validate-plugin/validator"

	_ "github.com/danielkraic/proto-validate-plugin/example/person" // register protobuf type

	pdk "github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
)

const (
	PluginVersion  = "0.0.1"
	PluginPriority = 1
)

type Config struct {
	IgnoreUnknownTypes bool `json:"ignore_unknown_types"`
}

func main() {
	server.StartServer(New, PluginVersion, PluginPriority)
}

func New() interface{} {
	return &Config{}
}

func (cfg *Config) Access(kong *pdk.PDK) {
	logReq(kong)

	isGRPCReq, err := isGRPCRequest(kong)
	if err != nil {
		kong.Log.Err(err.Error())
		return
	}

	if !isGRPCReq {
		return
	}

	err = cfg.validateGRPCReq(kong)
	if err != nil {
		kong.Log.Err(err.Error())
		kong.Response.ExitStatus(http.StatusBadRequest)
		return
	}
}

func isGRPCRequest(kong *pdk.PDK) (bool, error) {
	method, err := kong.Request.GetMethod()
	if err != nil {
		return false, fmt.Errorf("get request method: %w", err)
	}

	if method != http.MethodPost {
		return false, nil
	}

	contentType, err := kong.Request.GetHeader("Content-Type")
	if err != nil {
		return false, fmt.Errorf("get request Content-Type: %w", err)
	}

	if !strings.HasPrefix(contentType, "application/grpc") {
		return false, nil
	}

	return true, nil
}

func (cfg *Config) validateGRPCReq(kong *pdk.PDK) error {
	urlPath, err := kong.Request.GetPath()
	if err != nil {
		return fmt.Errorf("get request url path: %w", err)
	}

	body, err := kong.Request.GetRawBody()
	if err != nil {
		return fmt.Errorf("get request body: %w", err)
	}

	kong.Log.Info("validating request for " + urlPath)

	if urlPath != "" {
		//TODO: setup getting of request type from gRPC service method
		urlPath = "github.com.danielkraic.proto_validate_plugin.example.person.Person"
	}

	err = validator.Validate(urlPath, body)
	if err != nil {
		var uknownMsgType *validator.UnsupportedMessageType
		if cfg.IgnoreUnknownTypes && errors.As(err, &uknownMsgType) {
			kong.Log.Info("validation skipped for unknown type " + urlPath)
			return nil
		}

		return fmt.Errorf("validate request: %w", err)
	}

	return nil
}

func logReq(kong *pdk.PDK) {
	h, err := kong.Request.GetHeaders(100)
	if err != nil {
		kong.Log.Err("get request headers: " + err.Error())
		return
	}

	kong.Log.Info("headers count " + fmt.Sprint(len(h)))
	for k, v := range h {
		kong.Log.Info("header " + k + ": " + strings.Join(v, "; "))
	}

	body, err := kong.Request.GetRawBody()
	if err != nil {
		kong.Log.Err("get request body: " + err.Error())
		return
	}

	kong.Log.Info("request body " + base64.StdEncoding.EncodeToString(body))
}
