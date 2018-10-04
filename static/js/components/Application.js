import { Icon, Menu } from 'semantic-ui-react';
import { BrowserRouter as Router, Link, Route } from 'react-router-dom';
import React from 'react';

import Configurations from './Configurations';
import Mappings from './Mappings';
import Requests from './Requests';
import Variables from './Variables';

export default class Application extends React.Component {
  render() {
    return <Router>
      <React.Fragment>
        {this.renderMenu()}
        {this.renderContent()}
      </React.Fragment>
    </Router>;
  }

  renderContent() {
    return <React.Fragment>
      <Route exact path="/" component={Requests}/>
      <Route exact path="/mappings" component={Mappings}/>
      <Route exact path="/requests" component={Requests}/>
      <Route exact path="/variables" component={Variables}/>
      <Route exact path="/configurations" component={Configurations}/>
    </React.Fragment>;
  }

  renderMenu() {
    return <Menu>
      <Menu.Item>
        <a href="https://localhost:33443" target="_blank">
          <Icon name="external" />
          Go To Server
        </a>
      </Menu.Item>
      <Menu.Item><Link to="/requests">Requests</Link></Menu.Item>
      <Menu.Item><Link to="/mappings">Mappings</Link></Menu.Item>
      <Menu.Item><Link to="/variables">Variables</Link></Menu.Item>
      <Menu.Item><Link to="/configurations">Configurations</Link></Menu.Item>
    </Menu>;
  }

}
