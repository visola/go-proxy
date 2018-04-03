import { inject, observer } from 'mobx-react';
import PropTypes from 'prop-types';
import React from 'react';

import Field from './Field';

@inject('mappings')
@observer
export default class Mappings extends React.Component {
  static propTypes = {
    mappings: PropTypes.object.isRequired,
  }

  handleStatusChange(mapping, e) {
    this.props.mappings.updateMapping(mapping.mappingID, e.target.checked);
  }

  render() {
    return <ul>
      {this.renderMappings()}
    </ul>;
  }

  renderMappings() {
    let index = 0;
    const result = [];
    for (let origin in this.props.mappings.mappings) {
      const mappings = this.props.mappings.mappings[origin];
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
