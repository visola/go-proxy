import axios from 'axios';
import { action, observable } from 'mobx';

export default class Configurations {
  @observable mappings = [];

  @action
  fetch() {
    return axios.get('/configurations')
      .then(({data}) => {
        const result = {};
        data.forEach((mapping) => {
          const mappings = result[mapping.origin] || [];
          mappings.push(mapping);
          result[mapping.origin] = mappings;
        });
        this.mappings = result;
      });
  }
}
