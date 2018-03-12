import { Line, LineChart, XAxis, YAxis } from 'recharts';
import { inject, observer } from 'mobx-react';
import PropTypes from 'prop-types';
import React from 'react';

import Configurations from './Configurations';
import Field from './Field';
import Modal from './Modal';
import ProxiedRequestForm from './ProxiedRequestForm';

function toMillis(value) {
  return Math.round( value / 1000 );
}

@inject('proxiedRequests')
@observer
export default class Application extends React.Component {
  static propTypes = {
    proxiedRequests: PropTypes.object.isRequired,
  }

  constructor(props) {
    super(props);
    this.state = {
      selectedRequest: null,
    }
  }

  render() {
    return <div>
      {this.renderStatistics()}
      <Configurations />
      {this.renderRequests()}
    </div>;
  }

  renderModal() {
    const { selectedRequest } = this.state;

    if (!selectedRequest) {
      return null;
    }

    return <Modal title="Selected Request" onClose={() => this.setState({ selectedRequest: null})}>
      <ProxiedRequestForm request={selectedRequest} />
    </Modal>;
  }

  renderRequest(request, index) {
    return <tr key={index} onClick={() => this.showRequest(request)}>
      <td>{request.method}</td>
      <td>{request.requestedURL}</td>
    </tr>;
  }

  renderRequests() {
    let total = this.props.proxiedRequests.requests.length;
    let requests = this.props.proxiedRequests.requests.slice(Math.max(total - 50, 1));
    requests.reverse();
    return <div className="request-list">
      <h3>Last 50 requests:</h3>
      <table className="table">
        <thead>
          <tr>
            <th>Method</th>
            <th>Requested URL</th>
          </tr>
        </thead>
        <tbody>
          {requests.map((r, i) => this.renderRequest(r, i))}
        </tbody>
      </table>
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
