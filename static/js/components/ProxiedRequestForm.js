import PropTypes from 'prop-types';
import React from 'react';

import Field from './Field';

export default class ProxiedRequestForm extends React.Component {
  static propTypes = {
    request: PropTypes.object.isRequired,
  }

  constructor(props) {
    super(props);
    this.state = {
      showRequest: false,
    }

    this.toggleRequestData = this.toggleRequestData.bind(this);
  }

  render() {
    const { request } = this.props;
    return <div className="request">
      <Field label="Method:" value={request.method} />
      <Field label="Requested path:" value={request.requestedURL} />
      <Field label="Executed path:" value={request.executedURL} />
      <Field label="Status:" value={request.responseCode} />
      {this.renderRequestData(request)}
    </div>;
  }

  renderBody(body) {
    if (!body) {
      return null;
    }

    return <div>
      <h4>Body</h4>
      <pre className="request-body">{body}</pre>
    </div>;
  }

  renderHeaders(headers) {
    const renderedHeaders = [];

    for (const headerName in headers) {
      renderedHeaders.push(<li key={headerName}>
        <Field label={`${headerName}:`} value={headers[headerName].join(", ")} />
      </li>);
    }

    return <ul>{renderedHeaders}</ul>;
  }

  renderRequestData(request) {
    if (!this.state.showRequest) {
      return <h3 onClick={this.toggleRequestData}>
        <span className="glyphicon glyphicon-chevron-down" />
        Request
      </h3>;
    }

    return <div>
      <h3 onClick={this.toggleRequestData}>
        <span className="glyphicon glyphicon-chevron-up" />
        Request
      </h3>
      {this.renderHeaders(request.requestData.headers)}
      {this.renderBody(request.requestData.body)}
    </div>;
  }

  toggleRequestData() {
    this.setState({ showRequest: !this.state.showRequest });
  }
}
