import { Form, Input } from 'semantic-ui-react'
import { inject, observer } from 'mobx-react';
import React from 'react';

@inject('proxiedRequests')
@observer
export default class RequestFilter extends React.Component {
  constructor(props) {
    super(props);
    this.handleChange = this.handleChange.bind(this);
  }

  handleChange(e) {
    this.props.proxiedRequests.setFilter(e.target.value);
  }

  render() {
    const filter = this.props.proxiedRequests.filter;
    return <Form>
      <Form.Field>
        <Input
          onChange={this.handleChange}
          placeholder="Start typing to filter..."
          value={filter}
        />
      </Form.Field>
    </Form>;
  }
}
