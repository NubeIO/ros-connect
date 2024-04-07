package rubixbioscli

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/dto"
	"github.com/NubeIO/ros-connect/backend/nresty"
)

func (inst *Client) SystemCtlAction(serviceName string, action dto.Action) (*dto.Message, error) {
	url := fmt.Sprintf("/api/systemctl/%s?unit=%s", action, serviceName)
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&dto.Message{}).
		Post(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*dto.Message), nil
}
