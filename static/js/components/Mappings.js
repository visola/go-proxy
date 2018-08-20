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

  handleOriginStatusChange(origin, e) {
    const { mappings } = this.props.mappings;
    const newStatus = e.target.checked;
    mappings.filter((m) => m.origin === origin)
      .forEach((m) => m.active = newStatus);
    this.props.mappings.updateMappings(mappings);
  }

  handleStatusChange(mapping, e) {
    mapping.active = e.target.checked;
    this.props.mappings.updateMapping(mapping);
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
    const { countsPerOrigin, mappings } = this.props.mappings;
    return mappings.map((mapping, index) => {
      const count = countsPerOrigin[mapping.origin];
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
        <Table.Cell>
          <input
            checked={count.active === count.total}
            onChange={this.handleOriginStatusChange.bind(this, mapping.origin)}
            type="checkbox"
          />&nbsp;
          {mapping.origin}
        </Table.Cell>

      </Table.Row>;
    });
  }
}
