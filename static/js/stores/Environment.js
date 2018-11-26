import axios from 'axios';
import { action, observable } from 'mobx';

export default class Environment {
  @observable data = {};
  @observable loading = false;

  @action
  fetch() {
    this.loading = true;
    return axios.get('/api/environment')
      .then(({data}) => this.setData(data));
  }

  @action
  setData(data) {
    this.data = data;
    this.loading = false;
  }
}
