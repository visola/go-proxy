import PropTypes from 'prop-types';
import React from 'react';

export default class ProxiedRequestForm extends React.Component {
  static propTypes = {
    request: PropTypes.object.isRequired,
  }

  render() {
    const { request } = this.props;
    return <div>
      <div><label>Method:</label> {request.method}</div>
      <div><label>Requested path:</label> {request.requestedURL}</div>
      <div><label>Executed path:</label> {request.executedURL}</div>
      <div><label>Status:</label> {request.responseCode}</div>
      <h3>Request</h3>
      <div>
        <h4>Body</h4>
        <pre className="request-body">
          {request.requestData.body}
        </pre>
      </div>
    </div>;
  }
}
