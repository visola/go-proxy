import { Menu } from 'semantic-ui-react';
import React from 'react';

import Mappings from './Mappings';
import Requests from './Requests';

const tabs = [
  { id: 'requests', label: 'Requests', component: Requests },
  { id: 'mappings', label: 'Mappings', component: Mappings },
]

export default class Application extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      selectedTab: tabs[0],
    }
  }

  handleTabSelected(tabId) {
    const tab = tabs.find((t) => t.id == tabId);
    this.setState({ selectedTab: tab });
  }

  render() {
    return <div>
      <Menu>
        {tabs.map((t) => this.renderTab(t))}
      </Menu>
      {this.renderSelectedTab()}
    </div>;
  }

  renderSelectedTab() {
    const { selectedTab } = this.state;
    const Comp = selectedTab.component;
    return <Comp />
  }

  renderTab(tab) {
    return <Menu.Item
      active={tab.id == this.state.selectedTab.id}
      key={tab.id}
      onClick={() => this.handleTabSelected(tab.id)}
    >
      {tab.label}
    </Menu.Item>;
  }
}
