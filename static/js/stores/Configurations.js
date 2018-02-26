import axios from 'axios';
import { action, observable } from 'mobx';

export default class Configurations {
  @observable collection = [];

  @action
  fetch() {
    return axios.get('/configurations')
      .then(({data}) => {
        this.collection = data;
      });
  }
}
