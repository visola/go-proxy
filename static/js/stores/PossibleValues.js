import axios from 'axios';
import { action, observable } from 'mobx';

export default class PossibleValues {
  @observable data = {};
  @observable loading = false;

  @action
  fetch() {
    this.loading = true;
    return axios.get('/api/possible-values')
      .then(({data}) => this.setData(data));
  }

  @action
  setData(data) {
    this.data = data;
    this.loading = false;
  }
}
