import React, { useState, useEffect } from 'react'
import TagsInput from 'react-tagsinput'
import 'react-tagsinput/react-tagsinput.css'

const EditEndpointForm = (props) => {

  const [endpoint, setEndpoint] = useState(props.currentEndpoint)

  const handleInputChange = (event) => {
    const { name, value } = event.target

    setEndpoint({ ...endpoint, [name]: value })
  }
  
  const handleLabelChange = (labels) => {
    setEndpoint({ ...endpoint, "additional_labels": labels })
  }
  
  useEffect(() => {
    setEndpoint(props.currentEndpoint)
  }, [props])
  
  return (
    <form
      onSubmit={(event) => {
        event.preventDefault()

        props.updateEndpoint(endpoint.id, endpoint)
      }}
    >
      <label>Endpoint URL</label>
      <input
        type="text"
        name="url"
        value={endpoint.url}
        onChange={handleInputChange}
      />
      <label>Endpoint Label</label>
      <TagsInput name="additional_labels" value={endpoint.additional_labels} onChange={handleLabelChange} />
      <button>Update endpoint</button>
      <button
        onClick={() => props.setEditing(false)}
        className="button muted-button"
      >
        Cancel
      </button>
    </form>
  )
}

export default EditEndpointForm
