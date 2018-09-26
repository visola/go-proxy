import { inject, observer } from 'mobx-react';
import PropTypes from 'prop-types';
import React from 'react';

@inject('variables')
@observer
export default class Variables extends React.Component {
  static propTypes = {
    variables: PropTypes.object.isRequired,
  }

  render() {
    const { data, loading } = this.props.variables;
    if (loading) {
      return <p>Loading...</p>;
    }

    return data.map(v => (
      <div>
        <label>{v}</label>
      </div>
    ));
  }
}
