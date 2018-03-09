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
      showResponse: false,
    }

    this.toggleRequestData = this.toggleRequestData.bind(this);
    this.toggleResponseData = this.toggleResponseData.bind(this);
  }

  render() {
    const { request } = this.props;
    return <div className="request">
      <Field label="Method:" value={request.method} />
      <Field label="Requested path:" value={request.requestedURL} />
      <Field label="Executed path:" value={request.executedURL} />
      <Field label="Status:" value={request.responseCode} />
      {this.renderHTTPData("Request", this.state.showRequest, this.toggleRequestData, request.requestData)}
      {this.renderHTTPData("Response", this.state.showResponse, this.toggleResponseData, request.responseData)}
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

  renderHTTPData(title, expanded, toggleExpand, requestData) {
    if (!expanded) {
      return <h3 onClick={toggleExpand}>
        <span className="glyphicon glyphicon-chevron-down" />
        {title}
      </h3>;
    }

    return <div>
      <h3 onClick={toggleExpand}>
        <span className="glyphicon glyphicon-chevron-up" />
        {title}
      </h3>
      {this.renderHeaders(requestData.headers)}
      {this.renderBody(requestData.body)}
    </div>;
  }

  toggleRequestData() {
    this.setState({ showRequest: !this.state.showRequest });
  }

  toggleResponseData() {
    this.setState({ showResponse: !this.state.showResponse });
  }
}
