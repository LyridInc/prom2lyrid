import React, { useState, useEffect } from 'react'
import EndpointTable from './tables/EndpointTable'
import AddEndpointForm from './forms/AddEndpointForm'

const App = () => {

  const endpointsData = [
    

  ]

  const [endpoints, setEndpoints] = useState(endpointsData)
  
  const addEndpoint = (endpoint) => {
    const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(endpoint)
    };
    fetch('http://localhost:8081/endpoints/add', requestOptions)
    .then(res => res.json())
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
  
  const deleteEndpoint = (id) => {
    const requestOptions = {
      method: 'DELETE'
    };
    fetch('http://localhost:8081/endpoints/delete/'+id, requestOptions)
    .then(res => res.json())
    .then(
      (result) => {
        console.log(result)
      },
      (error) => {
        console.log(error)
      }
    )
    //setEndpoints(endpoints.filter((endpoint) => endpoint.id !== id))
  }
  
  useEffect(() => {
    fetch("http://localhost:8081/endpoints/list")
    .then(res => res.json())
    .then(
      (result) => {
        console.log(result)
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
  }, [])
  
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
