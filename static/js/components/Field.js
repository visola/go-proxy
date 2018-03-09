import PropTypes from 'prop-types';
import React from 'react';

export default class Field extends React.Component {
  static propTypes = {
    label: PropTypes.string.isRequired,
    value: PropTypes.oneOfType([
      PropTypes.string,
      PropTypes.number,
      PropTypes.bool,
    ]).isRequired
  }

  render() {
    return <div className="field">
      <label>{this.props.label}</label>
      <span>{this.props.value}</span>
    </div>;
  }
}
