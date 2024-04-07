package rubixbioscli

import (
	"fmt"
	"github.com/NubeIO/ros-connect/backend/installer"
	"sync"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

const RubixOsName = "rubix-os"

var (
	mutex   = &sync.RWMutex{}
	clients = map[string]*Client{}
)

type Client struct {
	Rest          *resty.Client
	Installer     *installer.Installer
	Ip            string `json:"ip"`
	Port          int    `json:"port"`
	HTTPS         bool   `json:"https"`
	ExternalToken string `json:"external_token"`
}

type ResponseBody struct {
	Response ResponseCommon `json:"response"`
	Status   string         `json:"status"`
	Count    string         `json:"count"`
}

type ResponseCommon struct {
	UUID string `json:"uuid"`
}

func buildUrl(overrideUrl ...string) string {
	if len(overrideUrl) > 0 {
		if overrideUrl[0] != "" {
			return overrideUrl[0]
		}
	}
	return ""
}

func New(cli *Client, installer *installer.Installer) *Client {
	mutex.Lock()
	defer mutex.Unlock()
	if cli == nil {
		log.Fatal("bios client cli can not be empty")
		return nil
	}
	baseURL := getBaseUrl(cli)
	if client, found := clients[baseURL]; found {
		return client
	}
	cli.Rest = resty.New()
	cli.Installer = installer
	cli.Rest.SetBaseURL(baseURL)
	cli.SetTokenHeader(cli.ExternalToken)
	clients[baseURL] = cli
	return cli
}

func ForceNew(cli *Client, installer *installer.Installer) *Client {
	mutex.Lock()
	defer mutex.Unlock()
	if cli == nil {
		log.Fatal("bios client cli can not be empty")
		return nil
	}
	cli.Rest = resty.New()
	cli.Installer = installer
	baseURL := getBaseUrl(cli)
	cli.Rest.SetBaseURL(baseURL)
	cli.SetTokenHeader(cli.ExternalToken)
	clients[baseURL] = cli
	return cli
}

func getBaseUrl(cli *Client) string {
	if cli.Ip == "" {
		cli.Ip = "0.0.0.0"
	}
	if cli.Port == 0 {
		cli.Port = 1659
	}
	var baseURL string
	if cli.HTTPS {
		baseURL = fmt.Sprintf("https://%s:%d", cli.Ip, cli.Port)
	} else {
		baseURL = fmt.Sprintf("http://%s:%d", cli.Ip, cli.Port)
	}
	return baseURL
}

func (inst *Client) SetTokenHeader(token string) *Client {
	inst.Rest.Header.Set("Authorization", composeToken(token))
	return inst
}

func composeToken(token string) string {
	return fmt.Sprintf("External %s", token)
}

type Path struct {
	Path string
}

type Response struct {
	StatusCode int         `json:"code"`
	Message    interface{} `json:"message"`
	resty      *resty.Response
}

func (response Response) buildResponse(resp *resty.Response, err error) *Response {
	response.StatusCode = resp.StatusCode()
	response.resty = resp
	if resp.IsError() {
		response.Message = resp.Error()
	}
	if resp.StatusCode() == 0 {
		response.Message = "server is unreachable"
		response.StatusCode = 503
	}
	return &response
}
