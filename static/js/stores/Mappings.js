import axios from 'axios';
import { action, observable } from 'mobx';

export default class Mappings {
  @observable countsPerOrigin = {};
  @observable hasCustomSorting = false;
  @observable mappings = {};

  @action
  fetch() {
    return axios.get('/mappings')
      .then(({data}) => this.setMappingsFromData(data));
  }

  @action
  updateMapping(mapping) {
    return axios.put(`/mappings/${mapping.mappingID}`, mapping)
      .then(({data}) => this.setMappingsFromData(data));
  }

  @action
  updateMappings(mappings) {
    return axios.put('/mappings', mappings)
      .then(({data}) => this.setMappingsFromData(data));
  }

  @action
  setMappingsFromData(data) {
    this.countsPerOrigin = {};
    this.hasCustomSorting = false;
    this.mappings = data;
    data.forEach((m) => {
      const counts = this.countsPerOrigin[m.origin] || {active:0,total:0};
      counts.total++;
      if (m.active) {
        counts.active++;
      }
      if (this.hasCustomSorting === false && m.before !== '') {
        this.hasCustomSorting = true;
      }
      this.countsPerOrigin[m.origin] = counts;
    });
  }
}
