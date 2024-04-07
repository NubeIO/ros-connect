package rubixbioscli

import (
	"fmt"
	"github.com/NubeIO/ros-connect/backend/nresty"
	"github.com/NubeIO/ros-connect/model"
)

func (inst *Client) BiosArch() (*model.Arch, error) {
	url := fmt.Sprintf("api/system/arch")
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&model.Arch{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Arch), nil
}
