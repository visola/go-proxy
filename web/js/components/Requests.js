import { inject, observer } from 'mobx-react';
import { Line, LineChart, XAxis, YAxis } from 'recharts';
import { Modal, Table } from 'semantic-ui-react'
import moment from 'moment';
import PropTypes from 'prop-types';
import React from 'react';

import ProxiedRequestForm from './ProxiedRequestForm';
import RequestFilter from './RequestFilter';

@inject('proxiedRequests')
@observer
export default class Requests extends React.Component {
  static propTypes = {
    proxiedRequests: PropTypes.object.isRequired,
  }

  constructor(props) {
    super(props);
    this.state = {
      selectedRequest: null,
    }

    this.handleClose = this.handleClose.bind(this);
  }

  handleClose() {
    this.setState({ selectedRequest: null });
  }

  render() {
    return <div>
      {this.renderStatistics()}
      <RequestFilter />
      {this.renderRequests()}
    </div>;
  }

  renderModal() {
    const { selectedRequest } = this.state;

    if (!selectedRequest) {
      return null;
    }

    return <Modal
      closeIcon="close"
      onClose={this.handleClose}
      open={this.state.selectedRequest != null}
    >
      <Modal.Header>Selected Request</Modal.Header>
      <Modal.Content>
        <ProxiedRequestForm request={selectedRequest} />
      </Modal.Content>
    </Modal>;
  }

  renderRequest(request, index) {
    return <Table.Row key={index} onClick={() => this.showRequest(request)}>
      <Table.Cell>{moment(request.startTime).format('HH:mm:ss SSS')}</Table.Cell>
      <Table.Cell>{request.endTime - request.startTime}ms</Table.Cell>
      <Table.Cell>{request.method}</Table.Cell>
      <Table.Cell>{request.responseCode}</Table.Cell>
      <Table.Cell>{request.requestedURL}</Table.Cell>
    </Table.Row>;
  }

  renderRequests() {
    let total = this.props.proxiedRequests.filtered.length;
    let requests = this.props.proxiedRequests.filtered.slice(Math.max(total - 50, 1));
    requests.reverse();
    return <div className="request-list">
      <Table>
        <Table.Header>
          <Table.Row>
            <Table.HeaderCell>Timestamp</Table.HeaderCell>
            <Table.HeaderCell>Duration</Table.HeaderCell>
            <Table.HeaderCell>Method</Table.HeaderCell>
            <Table.HeaderCell>Status</Table.HeaderCell>
            <Table.HeaderCell>URL</Table.HeaderCell>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {requests.map((r, i) => this.renderRequest(r, i))}
        </Table.Body>
      </Table>
      {this.renderModal()}
    </div>;
  }

  renderStatistics() {
    const dataArray  = this.props.proxiedRequests.requestsPerTimeBucket.slice();
    return <LineChart width={800} height={200} data={dataArray}>
      <XAxis dataKey="startString" />
      <YAxis />
      <Line activeDot={true} dataKey="count" fill="#8884d8" isAnimationActive={false} dot={false} />
    </LineChart>;
  }

  showRequest(request) {
    this.setState({ selectedRequest: request });
  }
}
