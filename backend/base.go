package backend

import (
	"errors"
	"fmt"
	"github.com/NubeIO/ros-connect/backend/config"
	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var db = "./db/db.yaml"

type Instance struct {
	Name          string `yaml:"name" json:"name"`
	UUID          string `yaml:"uuid" json:"uuid"`
	IP            string `yaml:"ip" json:"ip"`
	ExternalToken string `yaml:"external_token" json:"externalToken"`
	GRPCPort      int    `yaml:"grpcPort" json:"grpcPort"`
	Port          int    `yaml:"port" json:"port"`
	HTTPS         bool   `yaml:"https" json:"https"`
}

type App struct {
	Instances  map[string]*Instance
	mutex      sync.Mutex
	upgrader   websocket.Upgrader
	clients    map[*websocket.Conn]bool
	polling    map[string]*PollingInstance
	config     *config.Config
	configPath string
}

type PollingInstance struct {
	UUID         string          `json:"uuid"`
	InstanceUUID string          `json:"instanceUUID"`
	Timeout      int             `json:"timeout"`
	RefreshRate  time.Duration   `json:"refreshRate"`
	Options      *PollingOptions `json:"options"`
	ticker       *time.Ticker
	stopChan     chan struct{}
	counter      int
}

type PollingOptions struct {
}

type Options struct {
	Config *config.Config
}

func New(dbPath, configPath string, opts *Options) (*App, error) {
	db = dbPath
	file, err := os.Open(db)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	instances := make(map[string]*Instance)
	err = yaml.Unmarshal(data, &instances)
	if err != nil {
		return nil, err
	}
	return &App{
		Instances:  instances,
		mutex:      sync.Mutex{},
		clients:    make(map[*websocket.Conn]bool),
		polling:    make(map[string]*PollingInstance),
		config:     opts.Config,
		configPath: configPath,
	}, nil
}

func (inst *App) CreateInstance(instance *Instance) (*Instance, error) {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	id := newUUID()
	instance.UUID = id
	inst.Instances[id] = instance
	inst.saveToYAML()
	return instance, nil
}

func (inst *App) UpdateInstance(instance *Instance) error {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	if _, exists := inst.Instances[instance.UUID]; !exists {
		return errors.New("instance not found")
	}
	inst.Instances[instance.UUID] = instance
	inst.saveToYAML()
	return nil
}

func (inst *App) GetInstance(id string) (*Instance, error) {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	if instance, exists := inst.Instances[id]; !exists {
		return nil, errors.New("instance not found")
	} else {
		return instance, nil
	}
}

func (inst *App) DeleteInstance(hostUUID string) error {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()

	if _, exists := inst.Instances[hostUUID]; !exists {
		return errors.New("instance not found")
	}

	delete(inst.Instances, hostUUID)
	inst.saveToYAML()

	return nil
}

func (inst *App) GetAllInstances() map[string]*Instance {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	return inst.Instances
}

func (inst *App) saveToYAML() {
	data, err := yaml.Marshal(inst.Instances)
	if err != nil {
		fmt.Println("Error marshaling YAML:", err)
		return
	}

	err = ioutil.WriteFile(db, data, 0644)
	if err != nil {
		fmt.Println("Error writing YAML file:", err)
	}
}
