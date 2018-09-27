import axios from 'axios';
import { action, observable } from 'mobx';

export default class SelectedValues {
  @observable data = {};
  @observable loading = false;

  @action
  fetch() {
    this.loading = true;
    return axios.get('/api/values')
      .then(({data}) => this.setData(data));
  }

  @action
  setData(data) {
    this.data = data;
    this.loading = false;
  }

  @action
  setSelected(variable, newValue) {
    this.loading = true;
    return axios.put(`/api/values/${variable}`, JSON.stringify(newValue))
      .then(({data}) => this.setData(data));
  }
}
