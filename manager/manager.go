package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"os"
	"prom2lyrid/model"
	"sync"
	"time"
)

type NodeManager struct {
	ConfigFile string
	Node       model.Node

	ResultCache map[string]interface{}
	mux         sync.Mutex
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
	var nodeconfig model.Node
	jsonFile, err := os.Open(manager.ConfigFile)
	// if we os.Open returns an error then handle it
	if err != nil {
		// file does not exist
		log.Println("Config file not found, generating a new one")
		nodeconfig.ID = uuid.New().String()
	} else {
		log.Println("Config file found, loading")
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal([]byte(byteValue), &nodeconfig)
	}

	if nodeconfig.Endpoints == nil {
		nodeconfig.Endpoints = make(map[string]*model.ExporterEndpoint)
	}

	defer jsonFile.Close()

	manager.Node = nodeconfig

	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	manager.Node.HostName = name
	manager.ResultCache = make(map[string]interface{})
	manager.WriteConfig()

	for _, value := range manager.Node.Endpoints {
		value.SetUpdate(false)
		go value.Run(context.Background())
	}
}

func (manager *NodeManager) dumpresult() []interface{} {
	result := make([]interface{}, 0)

	index := 0
	for _, val := range manager.Node.Endpoints {
		if val.IsUpdated {
			result = append(result, val)
			val.SetUpdate(false)
		}
		index++
	}

	return result
}

func (manager *NodeManager) Run(ctx context.Context) {

	duration, _ := time.ParseDuration("30s")
	for c := time.Tick(duration); ; {

		manager.mux.Lock()
		// check every n-seconds for all the metrics that is collected and updated, dump it together to lyrid serverless
		//
		result := manager.dumpresult()
		fmt.Println(result)
		//manager.ResultCache = make(map[string]interface{})

		manager.mux.Unlock()

		select {
		case <-c:
			continue
		case <-ctx.Done():
			return
		}
	}
}

func (manager *NodeManager) WriteConfig() {
	file, _ := os.OpenFile(manager.ConfigFile, os.O_CREATE, os.ModePerm)
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.Encode(manager.Node)
}
