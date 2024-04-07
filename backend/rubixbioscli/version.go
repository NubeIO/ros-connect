package rubixbioscli

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/ros-connect/backend/bmodel"
	"github.com/NubeIO/ros-connect/backend/nresty"
	"github.com/NubeIO/ros-connect/constants"
)

func (inst *Client) GetRubixOsVersion() (*bmodel.Version, error) {
	installLocation := inst.Installer.GetAppInstallPath(constants.RubixOs)
	url := fmt.Sprintf("/api/files/list?path=%s", installLocation)
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&[]fileutils.FileDetails{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	versions := resp.Result().(*[]fileutils.FileDetails)
	if versions != nil && len(*versions) > 0 {
		return &bmodel.Version{Version: (*versions)[0].Name}, nil
	}
	return nil, errors.New("doesn't found the installation file")
}

func (inst *Client) GetRubixOsVersionV2() (*bmodel.Version, error, error) {
	installLocation := inst.Installer.GetAppInstallPath(constants.RubixOs)
	url := fmt.Sprintf("/api/files/list?path=%s", installLocation)
	resp, connectionErr, requestErr := nresty.FormatRestyV2Response(inst.Rest.R().
		SetResult(&[]fileutils.FileDetails{}).
		Get(url))
	if connectionErr != nil || requestErr != nil {
		return nil, connectionErr, requestErr
	}
	versions := resp.Result().(*[]fileutils.FileDetails)
	if versions != nil && len(*versions) > 0 {
		return &bmodel.Version{Version: (*versions)[0].Name}, nil, nil
	}
	return nil, nil, errors.New("doesn't found the installation file")
}
