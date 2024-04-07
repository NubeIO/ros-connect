package rubixbioscli

import (
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/dto"
	"github.com/NubeIO/ros-connect/backend/nresty"

	"io"
)

func (inst *Client) ListFilesV2(path string) ([]fileutils.FileDetails, error, error) {
	url := fmt.Sprintf("/api/files/list?path=%s", path)
	resp, connectionErr, requestErr := nresty.FormatRestyV2Response(inst.Rest.R().
		SetResult(&[]fileutils.FileDetails{}).
		Get(url))
	if connectionErr != nil || requestErr != nil {
		return nil, connectionErr, requestErr
	}
	return *resp.Result().(*[]fileutils.FileDetails), nil, nil
}

func (inst *Client) UploadFile(dir, file string, reader io.Reader) (*dto.UploadResponse, error) {
	url := fmt.Sprintf("/api/files/upload?destination=%s", dir)
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&dto.UploadResponse{}).
		SetFileReader("file", file, reader).
		Post(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*dto.UploadResponse), nil
}

func (inst *Client) CopyFile(from, to string) (*dto.Message, error) {
	url := fmt.Sprintf("/api/files/copy?from=%s&to=%s", from, to)
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&dto.Message{}).
		Post(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*dto.Message), nil
}

func (inst *Client) MoveFile(from, to string) (*dto.Message, error) {
	url := fmt.Sprintf("/api/files/move?from=%s&to=%s", from, to)
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&dto.Message{}).
		Post(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*dto.Message), nil
}

func (inst *Client) DeleteFiles(path string) (*dto.Message, error, error) {
	url := fmt.Sprintf("/api/files/delete-all?path=%s", path)
	resp, connectionErr, requestErr := nresty.FormatRestyV2Response(inst.Rest.R().
		SetResult(&dto.Message{}).
		Delete(url))
	if connectionErr != nil || requestErr != nil {
		return nil, connectionErr, requestErr
	}
	return resp.Result().(*dto.Message), nil, nil
}
