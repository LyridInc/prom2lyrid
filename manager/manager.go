package manager

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/LyridInc/go-sdk"
	sdkModel "github.com/LyridInc/go-sdk/model"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	"github.com/pierrec/lz4"
	"prom2lyrid/logger"
	"prom2lyrid/model"
	"prom2lyrid/utils"
)

type NodeManager struct {
	ConfigFile string
	Node       model.Node
	Apps       []*sdkModel.App

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
	manager.Node = nodeconfig
	jsonFile.Close()
	// Init Lyrid SDK
	if len(manager.Node.Credential.Key) > 0 && len(manager.Node.Credential.Secret) > 0 {
		sdk.GetInstance().Initialize(manager.Node.Credential.Key, manager.Node.Credential.Secret)
	}
	if manager.Node.IsLocal {
		sdk.GetInstance().SimulateServerless(manager.Node.ServerlessUrl)
	} else {
		sdk.GetInstance().DisableSimulate()
	}
	//sdk.GetInstance().SetLyridURL("http://localhost:8080") // for testing
	manager.Apps = sdk.GetInstance().GetApps()
	name, err := os.Hostname()
	if err != nil {
		level.Error(logger.GetInstance().Logger).Log("Error", err)
		panic(err)
	}
	manager.Node.HostName = name
	manager.ResultCache = make(map[string]interface{})
	manager.WriteConfig()
	manager.AddGateway(&manager.Node)

	for _, value := range manager.Node.Endpoints {
		value.Gateway = manager.Node.ID
		value.SetUpdate(false)
		manager.UpdateExporter(value)
		go value.Run(context.Background())
	}
}

func (manager *NodeManager) ExecuteFunction(body string) {
	for _, app := range manager.Apps {
		if strings.Contains(strings.ToLower(app.Name), strings.ToLower(os.Getenv("NOC_APP_NAME"))) {
			level.Debug(logger.GetInstance().Logger).Log("App name", app.Name)
			response, _ := sdk.GetInstance().ExecuteFunctionByName(app.Name, os.Getenv("NOC_MODULE_NAME"), os.Getenv("NOC_TAG"), os.Getenv("NOC_FUNCTION_NAME"), body)
			level.Debug(logger.GetInstance().Logger).Log("Response", response)
		}
	}
}

func (manager *NodeManager) ExecuteFunctionWithURIAndMethod(method string, uri string, body string) {
	for _, app := range manager.Apps {
		if strings.Contains(strings.ToLower(app.Name), strings.ToLower(os.Getenv("NOC_APP_NAME"))) {
			level.Debug(logger.GetInstance().Logger).Log("App name", app.Name)
			response, _ := sdk.GetInstance().ExecuteApp(app.Name, os.Getenv("NOC_MODULE_NAME"), os.Getenv("NOC_TAG"), os.Getenv("NOC_FUNCTION_NAME"), uri, method, body)
			level.Debug(logger.GetInstance().Logger).Log("Response", response)
		}
	}
}

func (manager *NodeManager) AddGateway(node *model.Node) {
	addGatewayBody := utils.JsonEncode(node)
	manager.ExecuteFunctionWithURIAndMethod("POST", "/api/gateways", addGatewayBody)
}

func (manager *NodeManager) UpdateExporter(exporter *model.ExporterEndpoint) {
	addExporterBody := utils.JsonEncode(exporter)
	manager.ExecuteFunctionWithURIAndMethod("POST", "/api/exporters", addExporterBody)
}

func (manager *NodeManager) UpdateScrapeResult(scrape *model.ScrapesEndpointResult) {
	updateScrapeBody := utils.JsonEncode(scrape)
	manager.ExecuteFunctionWithURIAndMethod("POST", "/api/scrapes", updateScrapeBody)
}

func (manager *NodeManager) DeleteExporter(exporter *model.ExporterEndpoint) {
	manager.ExecuteFunctionWithURIAndMethod("DELETE", "/api/exporters/"+exporter.ID, "")
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
			level.Info(logger.GetInstance().Logger).Log("Message", "UpdateScrapeResult for endpoint", "Endpoint", endpoint.URL)
			result, err := json.Marshal(endpoint.Result)
			scrapeResult := string(result)
			if endpoint.IsCompress {
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
				ExporterID:   endpoint.ID,
				ScrapeResult: scrapeResult,
				ScrapeTime:   endpoint.LastUpdateTime.UTC(),
				IsCompress:   endpoint.IsCompress,
			}

			manager.UpdateScrapeResult(&scrapeEndpointResult)
		}
	}
}
