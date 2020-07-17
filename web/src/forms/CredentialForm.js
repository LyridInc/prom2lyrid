import React, { useState, useEffect } from 'react'


const CredentialForm = (props) => {
  const initialFormState = {key: '', secret: ''}
  const [credential, setCredential] = useState(initialFormState)

  const handleInputChange = (event) => {
    const { name, value } = event.target
    setCredential({ ...credential, [name]: value })
  }
  
  useEffect(() => {
    setCredential(props.credential)
  }, [props])
  
  return (
    <form
      onSubmit={(event) => {
        event.preventDefault()
        props.updateCredential(credential)
      }}
    >
      <label>Key</label>
      <input type="text" name="key" value={credential.key} onChange={handleInputChange}/>
      <label>Secret</label>
      <input type="text" name="secret" value={credential.secret} onChange={handleInputChange}/>
      <button>Update</button>
    </form>
  )
}

export default CredentialForm
