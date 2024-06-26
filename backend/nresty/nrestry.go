package nresty

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/dto"
	"github.com/go-resty/resty/v2"
)

// An Error maps to Form3 API error responses
type Error struct {
	Code    int    `json:"error_code,omitempty"`
	Message string `json:"error_message,omitempty"`
}

func FormatRestyResponse(resp *resty.Response, err error) (*resty.Response, error) {
	// it catches errors:
	// => when we don't have host server (i/o timeout)
	//    -> e.g: `Post \"http://10.8.1.9:1616/api/users/login\": dial tcp 10.8.1.9:1616: i/o timeout`
	// => when we don't have app running (connection refused) etc...
	//    -> e.g: `Post \"http://10.8.1.9:1616/api/users/login\": dial tcp 10.8.1.9:1616: connect: connection refused`
	if err != nil {
		return resp, err
	}
	if resp.IsError() {
		return resp, composeErrorMsg(resp)
	}
	return resp, nil
}

// composeErrorMsg it helps to create a clean output error message; we used to have JSON message with nested key
func composeErrorMsg(resp *resty.Response) error {
	message := dto.Message{}
	rawMessage := resp.String()
	_ = json.Unmarshal([]byte(rawMessage), &message)

	if message.Message == "" {
		// if we do not have => `{"message": <message>}`
		message.Message = fmt.Sprintf("%s %s [%d]: %s",
			resp.Request.Method,
			resp.Request.URL,
			resp.StatusCode(),
			rawMessage)
	} else if message.Message == "not found" {
		// TODO: may be we don't need this
		// this is when rubix-service returns value as status_code 404; because of FF is stopped
		message.Message = fmt.Sprintf("%s %s [%d]: %s",
			resp.Request.Method,
			resp.Request.URL,
			resp.StatusCode(),
			message.Message)
	}
	e := fmt.Errorf(message.Message)
	return e
}

func FormatRestyV2Response(resp *resty.Response, err error) (res *resty.Response, connectionError error, requestError error) {
	// it catches errors:
	// => when we don't have host server (i/o timeout)
	//    -> e.g: `Post \"http://10.8.1.9:1616/api/users/login\": dial tcp 10.8.1.9:1616: i/o timeout`
	// => when we don't have app running (connection refused) etc...
	//    -> e.g: `Post \"http://10.8.1.9:1616/api/users/login\": dial tcp 10.8.1.9:1616: connect: connection refused`
	if err != nil {
		return resp, err, nil
	}
	if resp.IsError() && resp.StatusCode() != 404 && resp.StatusCode() != 400 {
		return nil, composeErrorMsg(resp), nil
	}
	if resp.IsError() {
		return nil, nil, composeErrorMsg(resp)
	}
	return resp, nil, nil
}
