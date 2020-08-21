package manager

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"github.com/LyridInc/go-sdk"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	"github.com/pierrec/lz4"
	"io"
	"io/ioutil"
	"log"
	"os"
	"prom2lyrid/logger"
	"prom2lyrid/model"
	"prom2lyrid/utils"
	"strings"
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
	manager.ConfigFile = os.Getenv("CONFIG_DIR") + "/config.json"
	manager.isUploading = false
	var nodeconfig model.Node

	jsonFile, err := os.Open(manager.ConfigFile)
	// if we os.Open returns an error then handle it
	if err != nil {
		// file does not exist
		level.Info(logger.GetInstance().Logger).Log("Message", "Config file not found, generating a new one")
		nodeconfig.ID = uuid.New().String()
		nodeconfig.IsLocal = true
		nodeconfig.ServerlessUrl = "http://localhost:8080"
	} else {
		level.Info(logger.GetInstance().Logger).Log("Message", "Config file found, loading")
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal([]byte(byteValue), &nodeconfig)
	}

	if nodeconfig.Endpoints == nil {
		nodeconfig.Endpoints = make(map[string]*model.ExporterEndpoint)
	}
	if manager.Node.IsLocal {
		sdk.GetInstance().SimulateServerless(manager.Node.ServerlessUrl)
	} else {
		sdk.GetInstance().DisableSimulate()
	}
	jsonFile.Close()

	manager.Node = nodeconfig

	name, err := os.Hostname()
	if err != nil {
		level.Error(logger.GetInstance().Logger).Log("Error", err)
		panic(err)
	}
	manager.Node.HostName = name
	manager.ResultCache = make(map[string]interface{})
	manager.WriteConfig()
	sdk.GetInstance().Initialize(manager.Node.Credential.Key, manager.Node.Credential.Secret)
	if manager.Node.IsLocal {
		sdk.GetInstance().SimulateServerless(manager.Node.ServerlessUrl)
	}
	sdk.GetInstance().ExecuteFunction(os.Getenv("FUNCTION_ID"), "LYR", utils.JsonEncode(model.LyFnInputParams{Command: "AddGateway", Gateway: manager.Node}))
	for _, value := range manager.Node.Endpoints {
		value.Gateway = manager.Node.ID
		value.SetUpdate(false)
		sdk.GetInstance().ExecuteFunction(os.Getenv("FUNCTION_ID"), "LYR", utils.JsonEncode(model.LyFnInputParams{Command: "AddExporter", Exporter: *value}))
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

	duration, _ := time.ParseDuration(os.Getenv("UPLOAD_INTERVAL"))
	for c := time.Tick(duration); ; {
		if !manager.isUploading {
			manager.isUploading = true

			manager.mux.Lock()
			// dump the cache temporarily

			manager.mux.Unlock()
			// check every n-seconds for all the metrics that is collected and updated, dump it together to lyrid serverless
			//
			//response, _ := sdk.GetInstance().ExecuteFunction("2054f61c-2d57-489f-a172-79fc15c6c20c", "LYR", utils.JsonEncode(model.LyFnInputParams{Command: "UpdateScrapeResult"}))
			//log.Println(string(response))

			level.Info(logger.GetInstance().Logger).Log("Message", "Uploading scrapes to gateway")
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
	_ = os.Mkdir(os.Getenv("CONFIG_DIR"), 0755)
	backupFile := os.Getenv("CONFIG_DIR") + "/config.json.bak." + time.Now().UTC().String()
	os.Rename(manager.ConfigFile, backupFile)
	f, _ := json.MarshalIndent(manager.Node, "", " ")
	_ = ioutil.WriteFile(manager.ConfigFile, f, 0644)
	manager.mux.Unlock()
}

func (manager *NodeManager) Upload() {
	for _, endpoint := range manager.Node.Endpoints {
		if endpoint.IsUpdated {
			level.Info(logger.GetInstance().Logger).Log("Message", "UpdateScrapeResult for endpoint", "Endpoint",  endpoint.URL)
			result, err := json.Marshal(endpoint.Result)
			scrapeResult := string(result)
			if endpoint.IsCompress{
				var writebuffer bytes.Buffer
				w := lz4.NewWriter(&writebuffer)
				io.Copy(w, strings.NewReader(string(result)))
				w.Close()
				scrapeResult = b64.StdEncoding.EncodeToString(writebuffer.Bytes())
			}
			if err != nil {
				log.Println(err)
				return
			}

			scrapeEndpointResult := model.ScrapesEndpointResult{
				ExporterID: endpoint.ID,
				ScrapeResult: scrapeResult,
				ScrapeTime: endpoint.LastUpdateTime.UTC(),
				IsCompress: endpoint.IsCompress,
			}
			response, _ := sdk.GetInstance().ExecuteFunction(os.Getenv("FUNCTION_ID"), "LYR", utils.JsonEncode(model.LyFnInputParams{Command: "UpdateScrapeResult", Exporter: *endpoint, ScapeResult: scrapeEndpointResult}))
			level.Debug(logger.GetInstance().Logger).Log("Response", response)
		}
	}
}
