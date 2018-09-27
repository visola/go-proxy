import { Button, Icon, Input, Label, Message, Modal, Table } from 'semantic-ui-react';
import { inject, observer } from 'mobx-react';
import PropTypes from 'prop-types';
import React from 'react';

@inject('possibleValues', 'selectedValues', 'variables')
@observer
export default class Variables extends React.Component {
  static propTypes = {
    possibleValues: PropTypes.object.isRequired,
    selectedValues: PropTypes.object.isRequired,
    variables: PropTypes.object.isRequired,
  }

  constructor(props) {
    super(props);
    this.state = {
      newValueVariableName: null,
      newValueVariableValue: "",
      showNewValueModal: false,
      selection: this.getSelection(),
    }
  }

  componentWillMount() {
    Promise.all([
      this.props.possibleValues.fetch(),
      this.props.selectedValues.fetch(),
    ]).then(() => this.setState({ selection: this.getSelection() }));
  }

  getSelection() {
    const selection = {};
    const selectedValues = this.props.selectedValues.data;
    for (let n in selectedValues) {
      selection[n] = selectedValues[n];
    }
    return selection;
  }

  handleAddValue(variable) {
    this.setState({
      newValueVariableName: variable,
      showNewValueModal: true,
    });
  }

  handleAddValueSave() {
    const { newValueVariableName, newValueVariableValue } = this.state;
    this.handleSelectionChange(newValueVariableName, newValueVariableValue);
    this.handleResetNewValueModal();
  }

  handleDeleteValue(e, variable, value) {
    e.preventDefault();
    e.stopPropagation();
    const values = (this.props.possibleValues.data[variable] || [])
      .filter((v) => v !== value);
    this.props.possibleValues.setPossibleValues(variable, values)
  }

  handleNewValueKeyPress(e) {
    if (e.key === 'Enter') {
      this.handleAddValueSave();
    }
  }

  handleResetNewValueModal() {
    this.setState({
      newValueVariableName: null,
      newValueVariableValue: "",
      showNewValueModal: false
    });
  }

  handleSelectionChange(variable, newValue) {
    const { selection } = this.state;
    if (newValue === selection[variable]) {
      return;
    }

    selection[variable] = newValue;
    this.setState({ selection });
    this.props.selectedValues.setSelected(variable, newValue)
      .then(() => this.props.possibleValues.fetch());
  }

  render() {
    if (this.props.possibleValues.loading || this.props.variables.loading) {
      return <p>Loading...</p>;
    }

    return <React.Fragment>
      <Message warning>
        Variables come from mappings. Add a variable on a mapping and they'll automatically show up here.
        If some mapping has a variable with no value set, it will fail with a 500.
      </Message>

      {this.renderNewValueModal()}

      <Table celled collapsing>
        <Table.Header>
          <Table.Row>
            <Table.HeaderCell>Variable</Table.HeaderCell>
            <Table.HeaderCell>Value</Table.HeaderCell>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {this.renderVariables()}
        </Table.Body>
      </Table>
    </React.Fragment>
  }

  renderNewValueModal() {
    return <Modal open={this.state.showNewValueModal} size="small">
      <Modal.Header>Type a new value for '{this.state.newValueVariableName}'</Modal.Header>
      <Modal.Content>
        <Input
          onChange={(e) => this.setState({newValueVariableValue: e.target.value})}
          onKeyPress={(e) => this.handleNewValueKeyPress(e)}
          placeholder="New Value"
          ref={(input) => {input ? input.focus() : null;} }
          type="text"
          value={this.state.newValueVariableValue}
        />
      </Modal.Content>
      <Modal.Actions>
        <Button negative onClick={() => this.handleResetNewValueModal() }>Cancel</Button>
        <Button onClick={() => this.handleAddValueSave()} positive>Save</Button>
      </Modal.Actions>
    </Modal>
  }

  renderValues(variable) {
    const values = (this.props.possibleValues.data[variable] || []).sort();
    const selectedValue = this.props.selectedValues.data[variable];
    
    const result = values.map(value => {
      if (value === selectedValue) {
        return <Label key={value} color='olive'>{value}</Label>
      } else {
        return <Label className="clickable" key={value} onClick={(e) => this.handleSelectionChange(variable, value)}>
          {value}
          <Icon color="red" name="delete" onClick={(e) => this.handleDeleteValue(e, variable, value)} />
        </Label>
      }
    });

    result.push(<Icon key="add" name="add circle" onClick={() => this.handleAddValue(variable)} />);
    return result;
  }

  renderVariables() {
    const variables = this.props.variables.data;
    return variables.map(v => (
      <Table.Row key={v}>
        <Table.Cell>{v}</Table.Cell>
        <Table.Cell className="values">{this.renderValues(v)}</Table.Cell>
      </Table.Row>
    ));
  }
}
