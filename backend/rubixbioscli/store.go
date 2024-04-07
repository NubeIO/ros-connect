package rubixbioscli

//
//func (inst *Client) GetPluginsStorePlugins() ([]installer.BuildDetails, error, error) {
//	pluginStore := inst.Installer.GetPluginsStorePath()
//	url := fmt.Sprintf("/api/files/list?path=%s", pluginStore)
//	resp, connectionErr, requestErr := nresty.FormatRestyV2Response(inst.Rest.R().
//		SetResult(&[]fileutils.FileDetails{}).
//		Get(url))
//	if connectionErr != nil || requestErr != nil {
//		return nil, connectionErr, requestErr
//	}
//
//	files := *resp.Result().(*[]fileutils.FileDetails)
//	plugins := make([]installer.BuildDetails, 0)
//	for _, file := range files {
//		plugins = append(plugins, *inst.Installer.GetZipBuildDetails(file.Name))
//	}
//	return plugins, nil, nil
//}
//
//func (inst *Client) UploadPluginStorePlugin(fileName string, reader io.Reader) error {
//	tmpDir := inst.Installer.GetEmptyNewTmpFolder()
//	_, err := inst.CreateDir(tmpDir)
//	if err != nil {
//		return errors.New(fmt.Sprintf("directory creation failed: %s", err.Error()))
//	}
//
//	_, err = inst.UploadFile(tmpDir, fileName, reader)
//	if err != nil {
//		return errors.New(fmt.Sprintf("upload plugin: %s", err.Error()))
//	}
//
//	from := path.Join(tmpDir, fileName)
//	to := inst.Installer.GetPluginsStoreWithFile(fileName)
//	_, err = inst.MoveFile(from, to)
//	if err != nil {
//		return err
//	}
//	_, _, _ = inst.DeleteFiles(tmpDir)
//	return nil
//}
