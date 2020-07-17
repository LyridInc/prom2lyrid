import React, { useState, useEffect } from 'react'
import TagsInput from 'react-tagsinput'
import 'react-tagsinput/react-tagsinput.css'
import {ObjectInput} from 'react-object-input'

const EditEndpointForm = (props) => {

  const [endpoint, setEndpoint] = useState(props.currentEndpoint)

  const handleInputChange = (event) => {
    const { name, value } = event.target

    setEndpoint({ ...endpoint, [name]: value })
  }
  
  const setLabels = (obj) => {
    setEndpoint({ ...endpoint, "additional_labels": obj })
  }
  
  const handleLabelChange = (labels) => {
    setEndpoint({ ...endpoint, "additional_labels": labels })
  }
  
  useEffect(() => {
    setEndpoint(props.currentEndpoint)
  }, [props])
  
  const [value, setValue] = useState(endpoint.additional_labels)
  return (
    <div>
      <ObjectInput
      obj={value}
      onChange={setValue}
      renderItem={(key, value, updateKey, updateValue, deleteProperty) => (
        // render an editor row for an item, using the provided callbacks
        // to propagate changes
        <div className="additional-label">
          <input 
            className="label-key"
            type="text"
            value={key}
            placeholder="key"
            onChange={e => updateKey(e.target.value)}
          />
          :
          <input 
            className="label-value"
            placeholder="value"
            type="text"
            value={value || ''} // value will be undefined for new rows
            onChange={e => updateValue(e.target.value)}
          />
          <button onClick={deleteProperty}>x</button>
        </div>
      )}
    />
    
    <form
      onSubmit={(event) => {
        event.preventDefault()
        handleLabelChange(value)
        endpoint.additional_labels = value
        console.log(value)
        console.log(endpoint)
        props.updateEndpoint(endpoint.id, endpoint)
      }}
    >
    <div>
      <button>Update endpoint</button>
      <button
        onClick={() => props.setEditing(false)}
        className="button muted-button"
      >
        Cancel
      </button>
    </div>
    </form>
    </div>
    
  )
}

export default EditEndpointForm
