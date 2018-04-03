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
  updateMapping(mappingID, status) {
    for (let origin in this.mappings) {
      this.mappings[origin].forEach((m) => {
        if (m.mappingID === mappingID && m.active !== status) {
          return axios.put(`/mappings/${mappingID}?active=${status}`)
            .then(({data}) => this.setMappingsFromData(data));
        }
      })
    }
  }

  @action
  setMappingsFromData(data) {
    const result = {};
    data.forEach((mapping) => {
      const mappings = result[mapping.origin] || [];
      mappings.push(mapping);
      result[mapping.origin] = mappings;
    });
    this.mappings = result;
  }
}
