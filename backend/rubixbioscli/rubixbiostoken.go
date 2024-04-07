package rubixbioscli

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/externaltoken"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/user"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/dto"
	"github.com/NubeIO/ros-connect/backend/nresty"
)

func (inst *Client) RubixBiosLogin(body *user.User) (*dto.TokenResponse, error) {
	url := "/api/users/login"
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetBody(body).
		SetResult(&dto.TokenResponse{}).
		Post(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*dto.TokenResponse), nil
}

func (inst *Client) RubixBiosTokens(jwtToken string) (*[]externaltoken.ExternalToken, error) {
	url := "/api/tokens"
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetHeader("Authorization", jwtToken).
		SetResult(&[]externaltoken.ExternalToken{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	data := resp.Result().(*[]externaltoken.ExternalToken)
	return data, nil
}

func (inst *Client) RubixBiosToken(jwtToken string, uuid string) (*externaltoken.ExternalToken, error) {
	url := fmt.Sprintf("/api/tokens/%s", uuid)
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetHeader("Authorization", jwtToken).
		SetResult(&externaltoken.ExternalToken{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	data := resp.Result().(*externaltoken.ExternalToken)
	return data, nil
}

func (inst *Client) RubixBiosTokenGenerate(jwtToken string, name string) (*externaltoken.ExternalToken, error) {
	url := "/api/tokens/generate"
	body := externaltoken.ExternalToken{Name: name, Blocked: false}
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetHeader("Authorization", jwtToken).
		SetBody(body).
		SetResult(&externaltoken.ExternalToken{}).
		Post(url))
	if err != nil {
		return nil, err
	}
	data := resp.Result().(*externaltoken.ExternalToken)
	return data, nil
}

func (inst *Client) RubixBiosTokenBlock(jwtToken string, uuid string, block bool) (*externaltoken.ExternalToken, error) {
	url := fmt.Sprintf("/api/tokens/%s/block", uuid)
	body := map[string]bool{"blocked": block}
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetHeader("Authorization", jwtToken).
		SetBody(body).
		SetResult(&externaltoken.ExternalToken{}).
		Put(url))
	if err != nil {
		return nil, err
	}
	data := resp.Result().(*externaltoken.ExternalToken)
	return data, nil
}

func (inst *Client) RubixBiosTokenRegenerate(jwtToken string, uuid string) (*externaltoken.ExternalToken, error) {
	url := fmt.Sprintf("/api/tokens/%s/regenerate", uuid)
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetHeader("Authorization", jwtToken).
		SetResult(&externaltoken.ExternalToken{}).
		Put(url))
	if err != nil {
		return nil, err
	}
	data := resp.Result().(*externaltoken.ExternalToken)
	return data, nil
}

func (inst *Client) RubixBiosTokenDelete(jwtToken string, uuid string) (bool, error) {
	url := fmt.Sprintf("/api/tokens/%s", uuid)
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetHeader("Authorization", jwtToken).
		Delete(url))
	if err != nil {
		return false, err
	}
	return resp.String() == "true", nil
}

func (inst *Client) RubixBiosUpdateUser(jwtToken, username, password string) (bool, error) {
	url := "/api/users"
	_, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetHeader("Authorization", jwtToken).
		SetBody(map[string]string{"username": username, "password": password}).
		Put(url))
	if err != nil {
		return false, err
	}
	return true, nil
}
