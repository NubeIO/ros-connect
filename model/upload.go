package model

type UploadResponse struct {
	Name         string `json:"name,omitempty"`
	Version      string `json:"version,omitempty"`
	UploadedOk   bool   `json:"uploaded_ok,omitempty"`
	TmpFile      string `json:"tmpFile,omitempty"`
	UploadedFile string `json:"uploadedFile,omitempty"`
}
