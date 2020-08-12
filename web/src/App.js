import React, { useState, useEffect } from 'react'
import EndpointTable from './tables/EndpointTable'
import AddEndpointForm from './forms/AddEndpointForm'
import EditEndpointForm from './forms/EditEndpointForm'
import CredentialForm from './forms/CredentialForm'

const App = () => {
  const endpointsData = []
  //const ROOT_URL = 'http://localhost:8081'
  const ROOT_URL = '';
  const [time, setTime] = useState(Date.now())
  const [endpoints, setEndpoints] = useState(endpointsData)
  const [credential, setCredential] = useState({key: '', secret: ''})
  const [editing, setEditing] = useState(false)
  const [isLocal, setLocal] = useState(true)
  const [serverlessUrl, setServerlessUrl] = useState("http://localhost:8080")
  const initialFormState = { id: null, url: '', additional_labels: [], config: {scrape_interval: "1m", scrape_timeout: "10m"}, is_compress:false}
  const [currentEndpoint, setCurrentEndpoint] = useState(initialFormState)
  
  const [lyridConnection, setLyridConnection] = useState({"status":"Checking Lyrid account ..."})
  
  const editRow = (endpoint) => {
    setEditing(true)
    let tags = []
    if (endpoint.additional_labels) {
      tags = endpoint.additional_labels
    }
    setCurrentEndpoint({ id: endpoint.id, url: endpoint.url, additional_labels: tags, config: endpoint.config, is_compress: endpoint.is_compress})
  }
  
  const updateEndpoint = (id, updatedEndpoint) => {
    const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updatedEndpoint)
    };
    fetch(ROOT_URL+'/endpoints/update/'+id+'/labels', requestOptions)
    .then((res) => {
        if(!res.ok) { throw new Error(res.status) } else {return res.json()}
    })
    .then(
      (result) => {
        //console.log(result)
        setEditing(false)
        setEndpoints(endpoints.map((endpoint) => (endpoint.id === id ? result : endpoint)))
      },
      (error) => {
        console.log(error)
      }
    )
  }
  
  const addEndpoint = (endpoint) => {
    const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(endpoint)
    };
    fetch(ROOT_URL+'/endpoints/add', requestOptions)
    .then((res) => {
        if(!res.ok) { throw new Error(res.status) } else {return res.json()}
    })
    .then(
      (result) => {
        console.log(result)
        setEndpoints([...endpoints, result])
      },
      (error) => {
        console.log(error)
      }
    )
  }
  
  const updateCredential = (credential) => {
    const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(credential)
    };
    fetch(ROOT_URL+'/config/credential', requestOptions)
    .then((res) => {
        if(!res.ok) { throw new Error(res.status) } else {return res.json()}
    })
    .then(
      (result) => {
        console.log(result)
        setCredential(result)
        setLyridConnection({"status":"Checking Lyrid account ..."})
        checkLyridConnection()
      },
      (error) => {
        console.log(error)
      }
    )
  }
  
  const deleteEndpoint = (id) => {
    const requestOptions = {
      method: 'DELETE'
    };
    fetch(ROOT_URL+'/endpoints/delete/'+id, requestOptions)
    setEndpoints(endpoints.filter((endpoint) => endpoint.id !== id))
  }
  
  const restartEndpoint = (id) => {
    fetch(ROOT_URL+"/endpoints/restart/"+id)
    .then((res) => {
        if(!res.ok) { throw new Error(res.status) } else {return res.json()}
    })
    .then((result) => {
        setEndpoints(endpoints.map((endpoint) => (endpoint.id === id ? result : endpoint)))
    })
  }
  
  const stopEndpoint = (id) => {
    fetch(ROOT_URL+"/endpoints/stop/"+id)
    .then((res) => {
        if(!res.ok) { throw new Error(res.status) } else {return res.json()}
    })
    .then((result) => {
        setEndpoints(endpoints.map((endpoint) => (endpoint.id === id ? result : endpoint)))
    })
  }
  
  const toggleLocal = () => {
    setLocal(!isLocal)
    const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({is_local: !isLocal})
    };
    fetch(ROOT_URL+"/config/local", requestOptions)
    .then((res) => {
        if(!res.ok) { throw new Error(res.status) } else {return res.json()}
    })
    .then((result) => {
        //setServerless(result)
    }) 
  }
  
  const postServerlessUrl = (target) => {
    setServerlessUrl(target.value)
    const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({url: target.value})
    };
    fetch(ROOT_URL+"/config/serverless", requestOptions)
    .then((res) => {
        if(!res.ok) { throw new Error(res.status) } else {return res.json()}
    })
    .then((result) => {
        //setServerless(result)
    })
  }
  
  const checkLyridConnection = () => {
    fetch(ROOT_URL+"/config/credential/status")
    .then((res) => {
        if(!res.ok) { throw new Error(res.status) } else {return res.json()}
    })
    .then((result) => {
        setLyridConnection({status: "Connected to Lyrid under accout name " + result[0].name + "."})
    })
  }
  
  useEffect(() => {
    const interval = setInterval(() => setTime(Date.now()), 60000)
    
    fetch(ROOT_URL+"/endpoints/list")
    .then((res) => {
        if(!res.ok) { throw new Error(res.status) } else {return res.json()}
    })
    .then(
      (result) => {
        //console.log(result)
        const keys = Object.keys(result)
        let eps = [];
        for (const key of keys) {
          eps.push(result[key])
        }
        setEndpoints(eps)
      },
      (error) => {
        console.log(error)
      }
    )
    
    checkLyridConnection()
    
    return () => {
        clearInterval(interval);
      }
  }, [time])
  
  useEffect(() => {
    fetch(ROOT_URL+"/config/credential")
    .then((res) => {
        if(!res.ok) { throw new Error(res.status) } else {return res.json()}
    })
    .then((result) => {
        setCredential(result)
    })
    
    fetch(ROOT_URL+"/config/local")
    .then((res) => {
        if(!res.ok) { throw new Error(res.status) } else {return res.json()}
    })
    .then((result) => {
        setLocal(result)
    })
    
    fetch(ROOT_URL+"/config/serverless")
    .then((res) => {
        if(!res.ok) { throw new Error(res.status) } else {return res.json()}
    })
    .then((result) => {
        setServerlessUrl(result)
    })
  }, []);
  
  return (
    <div className="container">
      <h1>prom2lyrid configuration page</h1>
      <div className="flex-row">
        <div className="flex-large">
          <label className="switch">
            <input type="checkbox" checked={isLocal} onChange={toggleLocal} />
            <div className="slider"></div>
          </label>
        </div>
      </div>
      <div className="flex-row">
        <div className="flex-large">
            { !isLocal ? (
            <div>
            <h2>Lyrid account</h2>
            <small>{lyridConnection.status}</small>
            <CredentialForm updateCredential={updateCredential} credential={credential} />
            </div>
            ) : (
            <div>
            <h2>Local URL</h2>
            <input type="text" name="url" value={serverlessUrl} onChange={postServerlessUrl}/>
            </div>
            )}
        </div>
      </div>
      <div className="flex-row">
        <div className="flex-large">
          {editing ? (
            <div>
              <h2>Edit endpoint</h2>
              <EditEndpointForm
                setEditing={setEditing}
                currentEndpoint={currentEndpoint}
                updateEndpoint={updateEndpoint}
              />
            </div>
          ) : (
            <div>
              <h2>Add endpoint</h2>
              <AddEndpointForm addEndpoint={addEndpoint} />
            </div>
          )}
        </div>
        
        <div className="flex-large">
          <h2>List endpoints 
            <button
              onClick={() => setTime(Date.now())}
              className="button muted-button"
            >
              Refresh
            </button>
          </h2>
          <EndpointTable endpoints={endpoints} editRow={editRow} deleteEndpoint={deleteEndpoint} restartEndpoint={restartEndpoint} stopEndpoint={stopEndpoint}/>
        </div>
      </div>
    </div>
  )
}

export default App
