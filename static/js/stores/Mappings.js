import axios from 'axios';
import { action, observable } from 'mobx';

export default class Mappings {
  @observable countsPerOrigin = {};
  @observable hasCustomSorting = false;
  @observable loading = false;
  @observable mappings = {};

  @action
  fetch() {
    this.loading = true;
    return axios.get('/api/mappings')
      .then(({data}) => this.setMappingsFromData(data));
  }

  @action
  updateMapping(mapping) {
    return axios.put(`/api/mappings/${mapping.mappingID}`, mapping)
      .then(({data}) => this.setMappingsFromData(data));
  }

  @action
  updateMappings(mappings) {
    return axios.put('/api/mappings', mappings)
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
    this.loading = false;
  }
}
