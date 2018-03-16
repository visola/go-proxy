import { inject, observer } from 'mobx-react';
import PropTypes from 'prop-types';
import React from 'react';

import Field from './Field';

@inject('configurations')
@observer
export default class Configurations extends React.Component {
  static propTypes = {
    configurations: PropTypes.object.isRequired,
  }

  handleStatusChange(mapping, e) {
    this.props.configurations.updateMapping(mapping.mappingID, e.target.checked);
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
        <Field label="Origin:" value={origin} />
        <ul>
          {mappings.map((m) => this.renderMapping(m))}
        </ul>
      </li>);
    }

    return result;
  }

  renderMapping(mapping) {
    return <div className="field" key={mapping.from}>
      <input
        checked={mapping.active} 
        onChange={this.handleStatusChange.bind(this, mapping)}
        type="checkbox"
      />&nbsp;
      <label>{mapping.proxy ? 'proxy' : 'static'}</label>
      <span>{`${mapping.from || mapping.regexp} => ${mapping.to}`}</span>
    </div>;
  }
}
