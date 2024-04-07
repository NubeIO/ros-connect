package rubixbioscli

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/dto"
	"github.com/NubeIO/ros-connect/backend/nresty"
)

func (inst *Client) CreateDir(path string) (*dto.Message, error) {
	url := fmt.Sprintf("/api/dirs/create?path=%s", path)
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&dto.Message{}).
		Post(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*dto.Message), nil
}
