import axios from 'axios';
import { action, observable } from 'mobx';

export default class Variables {
  @observable data = [];
  @observable loading = false;
  
  @action
  fetch() {
    this.loading = true;
    return axios.get('/api/variables')
      .then(({data}) => this.setData(data));
  }

  @action
  setData(data) {
    this.data = data.sort();
    this.loading = false;
  }
}
