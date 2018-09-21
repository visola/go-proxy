import { inject, observer } from 'mobx-react';
import PropTypes from 'prop-types';
import React from 'react';
import { SortableContainer, SortableElement, SortableHandle } from 'react-sortable-hoc';
import { Button, Table } from 'semantic-ui-react';

const DragHandle = SortableHandle(() => <span className="handle">::</span>);

@observer
class Row extends React.Component {
  render() {
    const {counts, handleOriginStatusChange, handleStatusChange, mapping} = this.props;
    return <Table.Row>
      <Table.Cell>
        <DragHandle />
        <input
          checked={mapping.active}
          onChange={(e) => handleStatusChange(mapping, e)}
          type="checkbox"
        />
      </Table.Cell>
      <Table.Cell>{mapping.proxy ? 'proxy' : 'static'}</Table.Cell>
      <Table.Cell>{mapping.from !== '' ? mapping.from : mapping.regexp}</Table.Cell>
      <Table.Cell>{mapping.to}</Table.Cell>
      <Table.Cell>
        <input
          checked={counts.active === counts.total}
          onChange={(e) => handleOriginStatusChange(mapping.origin, e)}
          type="checkbox"
        />&nbsp;
        {mapping.origin}
      </Table.Cell>
    </Table.Row>;
  }
}
const SortableRow = SortableElement(Row);

@observer
class TableBody extends React.Component {
  render() {
    const { handleOriginStatusChange, handleStatusChange, mappingsStore } = this.props;
    const { countsPerOrigin, mappings } = mappingsStore
    return <Table.Body>
      {mappings.map((mapping, index) => {
        return <SortableRow
          counts={countsPerOrigin[mapping.origin]}
          index={index}
          handleOriginStatusChange={handleOriginStatusChange}
          handleStatusChange={handleStatusChange}
          key={mapping.mappingID}
          mapping={mapping}
        />;
      })}
    </Table.Body>;
  }
}
const SortableTableBody = SortableContainer(TableBody);

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

  handleResetCustomSorting() {
    const { mappings } = this.props.mappings;
    mappings.filter((m) => m.before = null);
    this.props.mappings.updateMappings(mappings);
  }

  handleSortEnd({oldIndex, newIndex}) {
    const mappingsStore = this.props.mappings;

    const currentMapping = mappingsStore.mappings[oldIndex];
    const otherMapping = mappingsStore.mappings[newIndex];

    currentMapping.before = otherMapping.mappingID
    mappingsStore.updateMapping(currentMapping);
  }

  handleStatusChange(mapping, e) {
    mapping.active = e.target.checked;
    this.props.mappings.updateMapping(mapping);
  }

  render() {
    const { hasCustomSorting, loading } = this.props.mappings;
    if (loading) {
      return <p>Loading...</p>;
    }

    return <React.Fragment>
      <Button
        disabled={!hasCustomSorting}
        onClick={this.handleResetCustomSorting.bind(this)}
        primary
      >
        Reset Custom Sorting
      </Button>

      <Table>
        <Table.Header>
          <Table.Row>
            <Table.HeaderCell></Table.HeaderCell>
            <Table.HeaderCell>Type</Table.HeaderCell>
            <Table.HeaderCell>From</Table.HeaderCell>
            <Table.HeaderCell>To</Table.HeaderCell>
            <Table.HeaderCell>Origin</Table.HeaderCell>
          </Table.Row>
        </Table.Header>
        <SortableTableBody
          handleOriginStatusChange={this.handleOriginStatusChange.bind(this)}
          handleSortEnd={this.handleSortEnd.bind(this)}
          handleStatusChange={this.handleStatusChange.bind(this)}
          mappingsStore={this.props.mappings}
          onSortEnd={this.handleSortEnd.bind(this)}
          useDragHandle={true}
        />
      </Table>
    </React.Fragment>;
  }
}
