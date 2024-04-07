package rubixbioscli

import (
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/ros-connect/backend/nresty"
)

func (inst *Client) Unzip(source, destination string) (*[]fileutils.FileDetails, error) {
	url := fmt.Sprintf("/api/zip/unzip?source=%s&destination=%s", source, destination)
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&[]fileutils.FileDetails{}).
		Post(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*[]fileutils.FileDetails), nil
}
