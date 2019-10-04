import axios from 'axios';
import { action, observable } from 'mobx';

export default class Mappings {
  @observable countsPerOrigin = {};
  @observable countsPerTag = {};
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
    this.countsPerTag = {};
    this.hasCustomSorting = false;
    this.mappings = data;
    data.forEach((m) => {
      m.tags.forEach(tag => {
        const countsForTag = this.countsPerTag[tag] || {active:0,total:0};
        countsForTag.total++;
        this.countsPerTag[tag] = countsForTag;
      });

      const countsForOrigin = this.countsPerOrigin[m.origin] || {active:0,total:0};
      countsForOrigin.total++;

      if (m.active) {
        countsForOrigin.active++;

        m.tags.forEach(tag => {
          const countsForTag = this.countsPerTag[tag];
          countsForTag.active++;
          this.countsPerTag[tag] = countsForTag;
        });
      }

      if (this.hasCustomSorting === false && m.before !== '') {
        this.hasCustomSorting = true;
      }
      this.countsPerOrigin[m.origin] = countsForOrigin;
    });
    this.loading = false;
  }
}
