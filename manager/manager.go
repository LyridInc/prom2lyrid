package manager

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/LyridInc/go-sdk"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
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

	isUploading bool
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
	manager.isUploading = false
	var nodeconfig model.Node

	jsonFile, err := os.Open(manager.ConfigFile)
	// if we os.Open returns an error then handle it
	if err != nil {
		// file does not exist
		log.Println("Config file not found, generating a new one")
		nodeconfig.ID = uuid.New().String()
		nodeconfig.IsLocal = true
		nodeconfig.ServerlessUrl = "http://localhost:8080"
	} else {
		log.Println("Config file found, loading")
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal([]byte(byteValue), &nodeconfig)
	}

	if nodeconfig.Endpoints == nil {
		nodeconfig.Endpoints = make(map[string]*model.ExporterEndpoint)
	}

	jsonFile.Close()

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

	sdk.GetInstance().Initialize(manager.Node.Credential.Key, manager.Node.Credential.Secret)
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

	duration, _ := time.ParseDuration(os.Getenv("UPLOAD_INTERVAL"))
	for c := time.Tick(duration); ; {
		if !manager.isUploading {
			manager.isUploading = true

			manager.mux.Lock()
			// dump the cache temporarily

			manager.mux.Unlock()
			// check every n-seconds for all the metrics that is collected and updated, dump it together to lyrid serverless
			//

			log.Println("Uploading scrapes to gateway: ")
			manager.Upload()
			//result := manager.dumpresult()

			//
			//fmt.Println(result)
			//manager.ResultCache = make(map[string]interface{})
			manager.isUploading = false
		}

		select {
		case <-c:
			continue
		case <-ctx.Done():
			return

		}
	}
}

func (manager *NodeManager) WriteConfig() {
	manager.mux.Lock()
	//file, _ := os.OpenFile(manager.ConfigFile, os.O_CREATE, os.ModePerm)
	//defer file.Close()

	//encoder := json.NewEncoder(file)
	//encoder.Encode(manager.Node)

	f, _ := json.MarshalIndent(manager.Node, "", " ")
	_ = ioutil.WriteFile(manager.ConfigFile, f, 0644)
	manager.mux.Unlock()
}

func (manager *NodeManager) Upload() {
	for _, endpoint := range manager.Node.Endpoints {
		if endpoint.IsUpdated {

			// todo: Change to lyrid-sdk later
			//if (manager.Node.IsLocal) { }

			url := manager.Node.ServerlessUrl

			request := make(map[string]interface{})
			request["Command"] = "UpdateScrapeResult"
			scrapeResult := make(map[string]interface{})
			scrapeResult["ExporterID"] = endpoint.ID
			result, _ := json.Marshal(endpoint.Result)
			scrapeResult["ScrapeResult"] = string(result)
			scrapeResult["ScrapeTime"] = endpoint.LastUpdateTime.UTC()
			request["ScapeResult"] = scrapeResult

			jsonreq, _ := json.Marshal(request)
			fmt.Println()
			req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonreq))
			req.Header.Add("content-type", "application/json")
			response, err := http.DefaultClient.Do(req)
			if err != nil {
				return
			}

			body, _ := ioutil.ReadAll(response.Body)
			defer response.Body.Close()

			fmt.Println(string(body))
			/*
				var grant_json = "{\"grant_type\":\"client_credentials\"," +
					"\"client_id\": \"" + os.Getenv("AUTH0_CLIENTID") + "\"," +
					"\"client_secret\": \"" + os.Getenv("AUTH0_CLIENTSECRET") + "\"," +
					"\"audience\": \"https://" + os.Getenv("AUTH0_DOMAIN") + "/api/v2/\"}"

				req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(grant_json)))
				req.Header.Add("content-type", "application/json")
				response, err := http.DefaultClient.Do(req)
				if err != nil {
					sentry.CaptureException(err)
					return "", err
				}

				body, _ := ioutil.ReadAll(response.Body)
				defer response.Body.Close()

				var tokenjson map[string]interface{}
				json.Unmarshal(body, &tokenjson)
			*/
		}
	}
}
