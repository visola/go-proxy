import axios from 'axios';
import { action, observable } from 'mobx';

export default class Mappings {
  @observable mappings = {};

  @action
  fetch() {
    return axios.get('/mappings')
      .then(({data}) => this.setMappingsFromData(data));
  }

  @action
  updateMapping(mapping, status) {
    if (mapping.active === status) {
      return;
    }

    return axios.put(`/mappings/${mapping.mappingID}?active=${status}`)
      .then(({data}) => this.setMappingsFromData(data));
  }

  @action
  setMappingsFromData(data) {
    this.mappings = data;
  }
}
