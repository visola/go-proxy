import { Icon, Input, Label } from 'semantic-ui-react';
import { inject, observer } from 'mobx-react';
import PropTypes from 'prop-types';
import React from 'react';

@inject('configurations')
@observer
export default class Variables extends React.Component {
  static propTypes = {
    configurations: PropTypes.object.isRequired,
  }

  constructor(props) {
    super(props);
    this.state = {
      adding: false,
      addedText: "",
    };
  }

  handleAddedTextChange(e) {
    const newValue = e.target.value;
    this.setState({ addedText: newValue });
  }

  handleClickAddBaseDirectory() {
    this.setState({ adding: true });
  }

  handleAddBaseDirectoryKeyPress(e) {
    if (e.key === 'Enter') {
      this.handleSaveBaseDirectory();
    }
  }

  handleRemoveBaseDirectory(toRemove) {
    const { configurations } = this.props;
    const { data } = configurations;
    const indexOf = data.BaseDirectories.indexOf(toRemove);
    data.BaseDirectories.splice(indexOf, 1);
    configurations.save(data);
  }

  handleSaveBaseDirectory() {
    const { configurations } = this.props;
    const { data } = configurations;
    data.BaseDirectories.push(this.state.addedText);
    configurations.save(data);
    this.setState({ adding: false, addedText: "" });
  }

  render() {
    const { data, loading } = this.props.configurations;
    if (loading) {
      return <p>Loading...</p>;
    }

    return <div>
      <h3>Base Directories</h3>
      {this.renderAddBaseDirectory()}
      {data.BaseDirectories.map((baseDir) => this.renderBaseDirectory(baseDir))}
    </div>;
  }

  renderAddBaseDirectory() {
    const { adding, addedText } = this.state;
    if (adding) {
      return <div className="addBaseDirectory">
        <Input
          onChange={this.handleAddedTextChange.bind(this)}
          onKeyPress={this.handleAddBaseDirectoryKeyPress.bind(this)}
          ref={(input) => { input ? input.focus() : null }}
          value={addedText}
        />
        <Icon
          color="teal"
          name="save"
          onClick={this.handleSaveBaseDirectory.bind(this)}
          size="large"
        />
        <Icon
          color="red"
          name="cancel"
          onClick={() => this.setState({ adding: false })}
          size="large"
        />
      </div>;
    }

    return <div className="addBaseDirectory">
      <Icon
        className="clickable"
        name="add circle"
        onClick={this.handleClickAddBaseDirectory.bind(this)}
        size="large"
      />
    </div>;
  }

  renderBaseDirectory(baseDir) {
    return <div className="baseDirectory" key={baseDir}>
      <Label size="large">
        {baseDir}
        <Icon
          color="red"
          name="delete"
          onClick={this.handleRemoveBaseDirectory.bind(this, baseDir)}
        />
      </Label>
    </div>;
  }
}
