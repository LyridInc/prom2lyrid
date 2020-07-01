import React from 'react'

const EndpointTable = (props) => (
  <table>
    <thead>
      <tr>
        <th>Endpoint ID</th>
        <th>Endpoint URL</th>
        <th>Status</th>
        <th>Scrape Interval</th>
        <th>Scrape Timeout</th>
        <th>Last Scrape</th>
        <th>Updated On</th>
        <th>Actions</th>
      </tr>
    </thead>
    <tbody>
      {props.endpoints.length > 0 ? (
        props.endpoints.map((endpoint) => (
          <tr key={endpoint.id}>
            <td>{endpoint.id}</td>
            <td><a target="_blank" href={endpoint.url}>{endpoint.url}</a></td>
            <td>{endpoint.status}</td>
            <td>{endpoint.config.scrape_interval}</td>
            <td>{endpoint.config.scrape_timeout}</td>
            <td>{endpoint.last_scrape}</td>
            <td>{endpoint.LastUpdateTime}</td>
            <td>
              <button className="button muted-button">Edit</button>
              <button
                onClick={() => props.deleteEndpoint(endpoint.id)}
                className="button muted-button"
              >
                Delete
              </button>
            </td>
          </tr>
        ))
      ) : (
        <tr>
          <td colSpan={8}>No endpoints</td>
        </tr>
      )}
    </tbody>
  </table>
)

export default EndpointTable
