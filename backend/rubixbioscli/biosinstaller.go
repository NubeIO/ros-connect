package rubixbioscli

import (
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/dto"
	"github.com/NubeIO/ros-connect/backend/bmodel"
	"github.com/NubeIO/ros-connect/backend/namings"
	"github.com/NubeIO/ros-connect/backend/nresty"
	"github.com/NubeIO/ros-connect/constants"
	"github.com/hashicorp/go-version"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var backupExcludeDataFiles = map[string]bool{
	"plugins":   true,
	"modules":   true,
	"images":    true,
	"snapshots": true,
}

func (inst *Client) RubixOsUpload(body *dto.FileUpload) (*dto.Message, error) {
	uploadLocation := inst.Installer.GetAppDownloadPathWithVersion(constants.RubixOs, body.Version)
	url := fmt.Sprintf("/api/dirs/create?path=%s", uploadLocation)
	_, _ = nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&dto.Message{}).
		Post(url))

	url = fmt.Sprintf("/api/files/upload?destination=%s", uploadLocation)
	reader, err := os.Open(body.File)
	if err != nil {
		return nil, err
	}
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&bmodel.UploadResponse{}).
		SetFileReader("file", filepath.Base(body.File), reader).
		Post(url))
	if err != nil {
		return nil, err
	}
	upload := resp.Result().(*bmodel.UploadResponse)

	url = fmt.Sprintf("/api/zip/unzip?source=%s&destination=%s", upload.Destination, uploadLocation)
	resp, err = nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&[]fileutils.FileDetails{}).
		Post(url))
	if err != nil {
		return nil, err
	}
	unzippedFiles := resp.Result().(*[]fileutils.FileDetails)

	url = fmt.Sprintf("/api/files/delete?file=%s", upload.Destination)
	resp, err = nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&dto.Message{}).
		Delete(url))
	if err != nil {
		return nil, err
	}

	for _, f := range *unzippedFiles {
		from := path.Join(uploadLocation, f.Name)
		to := path.Join(uploadLocation, "app")
		url = fmt.Sprintf("/api/files/move?from=%s&to=%s", from, to)
		resp, err = nresty.FormatRestyResponse(inst.Rest.R().
			SetResult(&dto.Message{}).
			Post(url))
		if err != nil {
			return nil, err
		}
	}
	return &dto.Message{Message: "successfully uploaded the rubix-os in bios device"}, nil
}

func (inst *Client) RubixOsInstall(rubixOsVersion, rubixOsSystemdContent string) (*dto.Message, error) {
	// delete installed files
	installationDirectory := inst.Installer.GetAppInstallPath(constants.RubixOs)
	url := fmt.Sprintf("/api/files/delete-all?path=%s", installationDirectory)
	_, _ = nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&dto.Message{}).
		Delete(url))
	log.Println("deleted installed files, if any")

	downloadedFile := path.Join(inst.Installer.GetAppDownloadPathWithVersion(constants.RubixOs, rubixOsVersion), "app")
	installationFile := path.Join(inst.Installer.GetAppInstallPathWithVersion(constants.RubixOs, rubixOsVersion), "app")

	// create installation directory
	installationDirectoryWithVersion := inst.Installer.GetAppInstallPathWithVersion(constants.RubixOs, rubixOsVersion)
	url = fmt.Sprintf("/api/dirs/create?path=%s", installationDirectoryWithVersion)
	_, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&dto.Message{}).
		Post(url))
	if err != nil {
		return nil, err
	}
	log.Info("created installation directory: ", installationDirectoryWithVersion)

	// move downloaded file to installation directory
	url = fmt.Sprintf("/api/files/move?from=%s&to=%s", downloadedFile, installationFile)
	_, err = nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&dto.Message{}).
		Post(url))
	if err != nil {
		return nil, err
	}
	log.Info("moved downloaded file to installation directory")

	pluginDownloadedPath := inst.Installer.GetROSPluginDownloadPath()
	pluginPath := inst.Installer.GetROSPluginInstallPath()
	_, _, _ = inst.DeleteFiles(pluginPath)
	_, _ = inst.MoveFile(pluginDownloadedPath, pluginPath) // ignore error: sometimes from folder will be empty
	log.Info("moved plugins in data directory")

	wd := inst.Installer.GetAppInstallPathWithVersion(constants.RubixOs, rubixOsVersion)
	log.Infof("systemd file with working directory: %s", wd)
	rubixOsSystemdContent = strings.Replace(rubixOsSystemdContent, "<working_dir>", wd, -1)

	serverFlag, err := inst.getServerFlag(rubixOsVersion)
	if err != nil {
		return nil, err
	}
	rubixOsSystemdContent = strings.Replace(rubixOsSystemdContent, "<server_flag>", serverFlag, -1)

	message, err := inst.installServiceFile(constants.RubixOs, rubixOsSystemdContent)
	if err != nil {
		return message, err
	}
	return &dto.Message{Message: "successfully installed the rubix-os in bios device"}, nil
}

func (inst *Client) getServerFlag(rubixOsVersion string) (string, error) {
	v, err := version.NewVersion(rubixOsVersion)
	if err != nil {
		return "", err
	}
	mV, _ := version.NewVersion("v0.3.4")
	if v.GreaterThan(mV) {
		return "--server", nil
	}
	return "", nil
}

func (inst *Client) installServiceFile(appName, rubixOsSystemdContent string) (*dto.Message, error) {
	serviceFileName := namings.GetServiceNameFromAppName(appName)
	serviceFile := path.Join(constants.ServiceDir, serviceFileName)
	symlinkServiceFile := path.Join(constants.ServiceDirSoftLink, serviceFileName)
	url := fmt.Sprintf("/api/files/upload?destination=%s", constants.ServiceDir)
	reader := strings.NewReader(rubixOsSystemdContent)
	if _, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetFileReader("file", serviceFileName, reader).
		SetResult(&bmodel.UploadResponse{}).
		Post(url)); err != nil {
		return nil, err
	}
	log.Info("service file is uploaded successfully")

	url = fmt.Sprintf("/api/syscall/unlink?path=%s", symlinkServiceFile)
	if _, err := nresty.FormatRestyResponse(inst.Rest.R().Post(url)); err != nil {
		log.Error(err)
	}
	log.Infof("soft un-linked %s", symlinkServiceFile)

	url = fmt.Sprintf("/api/syscall/link?path=%s&link=%s", serviceFile, symlinkServiceFile)
	if _, err := nresty.FormatRestyResponse(inst.Rest.R().Post(url)); err != nil {
		log.Error(err)
	}
	log.Infof("soft linked %s to %s", serviceFile, symlinkServiceFile)

	url = "/api/systemctl/daemon-reload"
	if _, err := nresty.FormatRestyResponse(inst.Rest.R().Post(url)); err != nil {
		log.Error(err)
	}
	log.Infof("daemon reloaded")

	url = fmt.Sprintf("/api/systemctl/enable?unit=%s", serviceFileName)
	if _, err := nresty.FormatRestyResponse(inst.Rest.R().Post(url)); err != nil {
		log.Error(err)
	}
	log.Infof("enabled service %s", serviceFileName)

	url = fmt.Sprintf("/api/systemctl/restart?unit=%s", serviceFileName)
	if _, err := nresty.FormatRestyResponse(inst.Rest.R().Post(url)); err != nil {
		log.Error(err)
	}
	log.Infof("started service %s", serviceFileName)
	return nil, nil
}

func (inst *Client) ListPlugins() ([]dto.Plugin, error, error) {
	p := inst.Installer.GetROSPluginInstallPath()
	files, connectionErr, requestErr := inst.ListFilesV2(p)
	if connectionErr != nil || requestErr != nil {
		return nil, connectionErr, requestErr
	}
	var plugins []dto.Plugin
	for _, file := range files {
		plugins = append(plugins, *inst.Installer.GetPluginDetails(file.Name))
	}
	return plugins, nil, nil
}

func (inst *Client) BackupRubixOsDataDir() error {
	rubixOsVersion, connectionErr, requestErr := inst.GetRubixOsVersionV2()
	if requestErr != nil {
		// it takes here on the first time of ROS installation, when we don't have files
		log.Warnf(requestErr.Error())
		return nil
	} else if connectionErr != nil {
		return connectionErr
	}

	from := inst.Installer.GetAppDataDataPath(constants.RubixOs)
	to := inst.Installer.GetAppBackupPath(constants.RubixOs, rubixOsVersion.Version)

	files, connectionErr, requestErr := inst.ListFilesV2(from)
	if requestErr != nil {
		// it takes here only if we try to list files on files instead of folder --almost no chance
		log.Warnf("backup process from %s to %s got failed (%s)", from, to, requestErr)
		return nil
	} else if connectionErr != nil {
		return connectionErr
	}

	log.Infof("backing up from %s to %s", from, to)
	for _, file := range files {
		if backupExcludeDataFiles[file.Name] {
			continue
		}
		fromFile := path.Join(from, file.Name)
		toFile := path.Join(to, file.Name)
		url := fmt.Sprintf("/api/files/copy?from=%s&to=%s", fromFile, toFile)
		_, err := nresty.FormatRestyResponse(inst.Rest.R().Post(url))
		if err != nil {
			log.Warnf("backup process from %s to %s got failed (%s)", fromFile, toFile, requestErr)
			return err
		}
	}
	log.Infof("backup process has been completed from %s to %s", from, to)
	return nil
}
