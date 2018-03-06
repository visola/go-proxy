import { Line, LineChart, XAxis, YAxis } from 'recharts';
import { inject, observer } from 'mobx-react';
import PropTypes from 'prop-types';
import React from 'react';

function toMillis(value) {
  return Math.round( value / 1000 );
}

@inject('configurations', 'proxiedRequests')
@observer
export default class Application extends React.Component {
  static propTypes = {
    configurations: PropTypes.object.isRequired,
    proxiedRequests: PropTypes.object.isRequired,
  }

  render() {
    return <div>
      {this.renderStatistics()}
      <ul>
        {this.renderConfigurations()}
      </ul>
      {this.renderRequests()}
    </div>;
  }

  renderConfigurations() {
    let index = 0;
    const result = [];
    for (let origin in this.props.configurations.mappings) {
      const mappings = this.props.configurations.mappings[origin];
      index += 1;
      result.push(<li key={index}>
        <span className="field">
          <label>Origin:</label><span>{origin}</span>
        </span>
        <ul>
          {mappings.map((m) => this.renderMapping(m))}
        </ul>
      </li>);
    }

    return result;
  }

  renderMapping(mapping) {
    return <span className="field" key={mapping.from}>
      <label>{mapping.proxy ? 'proxy' : 'static'}</label> {mapping.from} => {mapping.to}
    </span>;
  }

  renderRequest(request, index) {
    return <tr key={index}>
      <td>{request.responseCode}</td>
      <td>{request.endTime - request.startTime}ms</td>
      <td>{request.requestedPath}</td>
      <td>{request.proxiedTo}</td>
    </tr>;
  }

  renderRequests() {
    let total = this.props.proxiedRequests.requests.length;
    let requests = this.props.proxiedRequests.requests.slice(Math.max(total - 50, 1));
    requests.reverse();
    return <div>
      <h3>Last 50 requests:</h3>
      <table className="table">
        <thead>
          <tr>
            <th>Status</th>
            <th>Duration</th>
            <th>Original Path</th>
            <th>Proxied to</th>
          </tr>
        </thead>
        <tbody>
          {requests.map((r, i) => this.renderRequest(r, i))}
        </tbody>
      </table>
    </div>;
  }

  renderStatistics() {
    const dataArray  = this.props.proxiedRequests.requestsPerTimeBucket.slice();
    return <LineChart width={1024} height={200} data={dataArray}>
      <XAxis dataKey="startString" />
      <YAxis />
      <Line activeDot={true} dataKey="count" fill="#8884d8" isAnimationActive={false} dot={false} />
    </LineChart>;
  }
}
