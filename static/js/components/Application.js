import { inject, observer } from 'mobx-react';
import PropTypes from 'prop-types';
import React from 'react';

@inject('configurations')
@observer
export default class Application extends React.Component {
  static propTypes = {
    configurations: PropTypes.object.isRequired,
  }

  render() {
    return <ul>
      {this.renderConfigurations()}
    </ul>;
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
}
