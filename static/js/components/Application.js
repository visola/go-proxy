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
    return this.props.configurations.collection.map((config, index) => {
      return <li key={index}>
        <span className="field">
          <label>Origin:</label><span>{config.origin}</span>
        </span>
        <span className="field">
          <label>{config.proxy ? 'proxy' : 'static'}</label> {config.from} => {config.to
        }</span>
      </li>
    });
  }
}
