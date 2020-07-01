import React, { useState } from 'react'
import EndpointTable from './tables/EndpointTable'
import AddEndpointForm from './forms/AddEndpointForm'

const App = () => {

  const endpointsData = [
      {
        "id": "087730ab-bcc6-4f3a-87ba-07924b8f966d",
        "url": "http://10.1.17.49:9182/metrics",
        "config":{
          "scrape_interval": "1m",
          "scrape_timeout": "10m"
        },
        "status": "Warning",
        "last_scrape": "0001-01-01T00:00:00Z",
        "additional_labels": null,
        "message": "",
        "is_updated": false,
        "DurationSinceLastUpdate": 9223372036854775807,
        "LastUpdateTime": "0001-01-01T00:00:00Z"
      },
      {
        "id": "687cdaee-3976-49dc-8507-0dcf34bc5cdc",
        "url": "http://10.1.17.63:9216/metrics",
        "config":{
          "scrape_interval": "1m",
          "scrape_timeout": "10m"
        },
        "status": "Warning",
        "last_scrape": "2020-06-08T23:37:04.8704619Z",
        "additional_labels": null,
        "message": "",
        "is_updated": false,
        "DurationSinceLastUpdate": 1909224241062098,
        "LastUpdateTime": "2020-06-08T23:37:04.8704619Z"
      },
      {
        "id": "687cdaee-3976-49dc-8507-0dcf34bc5cde",
        "url": "http://localhost:9100/metrics",
        "config":{
          "scrape_interval": "1m",
          "scrape_timeout": "10m"
        },
        "status": "Warning",
        "last_scrape": "2020-06-08T23:37:04.8704619Z",
        "additional_labels": null,
        "message": "",
        "is_updated": false,
        "DurationSinceLastUpdate": 1909224241062296,
        "LastUpdateTime": "2020-06-08T23:37:04.8704619Z"
      }
    ]

  const [endpoints, setEndpoints] = useState(endpointsData)
  
  const addEndpoint = (endpoint) => {
    endpoint.id = endpoints.length + 1
    endpoint.config.scrape_interval = '1m'
    endpoint.config.scrape_timeout = '10m'
    endpoint.status = "OK"
    endpoint.last_scrape = "2020-06-08T23:37:04.8704619Z"
    endpoint.LastUpdateTime = "2020-06-08T23:37:04.8704619Z"
    setEndpoints([...endpoints, endpoint])
  }
  
  const deleteEndpoint = (id) => {
    setEndpoints(endpoints.filter((endpoint) => endpoint.id !== id))
  }
  
  return (
    <div className="container">
      <h1>prom2lyrid configuration page</h1>
      <div className="flex-row">
        <div className="flex-large">
          <h2>Add endpoint</h2>
          <AddEndpointForm addEndpoint={addEndpoint} />
        </div>
        <div className="flex-large">
          <h2>List endpoints</h2>
          <EndpointTable endpoints={endpoints} deleteEndpoint={deleteEndpoint}/>
        </div>
      </div>
    </div>
  )
}

export default App
