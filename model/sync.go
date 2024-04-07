package model

type SyncModel struct {
	UUID    string  `json:"uuid"`
	IsError bool    `json:"isError"`
	Message *string `json:"message"`
}
