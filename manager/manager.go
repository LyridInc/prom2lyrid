package manager

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"prom2lyrid/model"
	"sync"
)

type NodeManager struct {
	ConfigFile string
	Node       model.Node

	ResultCache map[string]interface{}
}

var instance *NodeManager
var once sync.Once

func GetInstance() *NodeManager {
	once.Do(func() {
		instance = &NodeManager{}
	})
	return instance
}

func (manager *NodeManager) Init() {
	manager.ConfigFile = os.Getenv("CONFIG_FILE")

	jsonFile, err := os.Open(manager.ConfigFile)
	// if we os.Open returns an error then handle it
	if err != nil {

		// file does not exist
	}
	defer jsonFile.Close()

	var nodeconfig model.Node

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &nodeconfig)

	manager.Node = nodeconfig

	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	manager.Node.HostName = name
	manager.ResultCache = make(map[string]interface{})
}

func (manager *NodeManager) Run() {

}
