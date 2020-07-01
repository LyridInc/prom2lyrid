import React, { useState } from 'react'

const AddEndpointForm = (props) => {
  const initialFormState = {id: null, url: '', config:{scrape_interval:"1m", scrape_timeout: "10m"}, last_scrape:'', LastUpdateTime: ''}
  const [endpoint, setEndpoint] = useState(initialFormState)

  const handleInputChange = (event) => {
    const { name, value } = event.target

    setEndpoint({ ...endpoint, [name]: value })
  }
  return (
    <form
      onSubmit={(event) => {
        event.preventDefault()
        if (!endpoint.url) return

        props.addEndpoint(endpoint)
        setEndpoint(initialFormState)
      }}
    >
      <label>Endpoint URL</label>
      <input type="text" name="url" value={endpoint.url} onChange={handleInputChange}/>
      <button>Add new endpoint</button>
    </form>
  )
}

export default AddEndpointForm
