import axios from 'axios';
import { Provider } from 'mobx-react';
import React from 'react';
import { render } from 'react-dom';

import Application from './components/Application';
import stores from './stores';

axios.defaults.headers.post['Content-Type'] = 'application/json';
axios.defaults.headers.put['Content-Type'] = 'application/json';

const ApplicationWithState = () => (
  <Provider {...stores}>
    <Application />
  </Provider>
);

render(<ApplicationWithState />, document.getElementById('container'));
