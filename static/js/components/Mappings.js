import { inject, observer } from 'mobx-react';
import PropTypes from 'prop-types';
import React from 'react';
import { Table } from 'semantic-ui-react';

@inject('mappings')
@observer
export default class Mappings extends React.Component {
  static propTypes = {
    mappings: PropTypes.object.isRequired,
  }

  handleStatusChange(mapping, e) {
    this.props.mappings.updateMapping(mapping, e.target.checked);
  }

  render() {
    return <Table>
      <Table.Header>
        <Table.Row>
          <Table.HeaderCell></Table.HeaderCell>
          <Table.HeaderCell>Type</Table.HeaderCell>
          <Table.HeaderCell>From</Table.HeaderCell>
          <Table.HeaderCell>To</Table.HeaderCell>
          <Table.HeaderCell>Origin</Table.HeaderCell>
        </Table.Row>
      </Table.Header>
      <Table.Body>
        {this.renderMappings()}
      </Table.Body>
    </Table>;
  }

  renderMappings() {
    return this.props.mappings.mappings.map((mapping, index) => {
      return <Table.Row key={mapping.mappingID}>
        <Table.Cell>
          <input
            checked={mapping.active}
            onChange={this.handleStatusChange.bind(this, mapping)}
            type="checkbox"
          />
        </Table.Cell>
        <Table.Cell>{mapping.proxy ? 'proxy' : 'static'}</Table.Cell>
        <Table.Cell>{mapping.from}</Table.Cell>
        <Table.Cell>{mapping.to}</Table.Cell>
        <Table.Cell>{mapping.origin}</Table.Cell>

      </Table.Row>;
    });
  }
}
