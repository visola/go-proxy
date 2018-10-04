import { Icon, Label } from 'semantic-ui-react';
import { inject, observer } from 'mobx-react';
import PropTypes from 'prop-types';
import React from 'react';

@inject('configurations')
@observer
export default class Variables extends React.Component {
  static propTypes = {
    configurations: PropTypes.object.isRequired,
  }

  render() {
    const { data, loading } = this.props.configurations;
    if (loading) {
      return <p>Loading...</p>;
    }

    return <div>
      <h3>Base Directories</h3>
      <p><Icon name="add circle" size="large" /></p>
      {data.BaseDirectories.map(this.renderBaseDirectory)}
    </div>;
  }

  renderBaseDirectory(baseDir) {
    return <div key={baseDir}>
      <Label>
        {baseDir}
        <Icon name="delete" />
      </Label>
    </div>;
  }
}
