import axios from 'axios';
import { action, observable } from 'mobx';

export default class Configurations {
  @observable data = {};
  @observable loading = false;

  @action
  fetch() {
    this.loading = true;
    return axios.get('/api/configurations')
      .then(({data}) => this.setData(data));
  }

  @action
  save(configurations) {
    this.loading = true;
    return axios.put('/api/configurations', configurations)
      .then(({data}) => this.setData(data));
  }

  @action
  setData(data) {
    this.data = data;
    this.loading = false;
  }
}
