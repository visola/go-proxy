import { Menu } from 'semantic-ui-react';
import { BrowserRouter as Router, Link, Route } from 'react-router-dom';
import React from 'react';

import Mappings from './Mappings';
import Requests from './Requests';

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
    </React.Fragment>;
  }

  renderMenu() {
    return <Menu>
      <Menu.Item><Link to="/requests">Requests</Link></Menu.Item>
      <Menu.Item><Link to="/mappings">Mappings</Link></Menu.Item>
    </Menu>;
  }

}
