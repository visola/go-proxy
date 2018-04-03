import { action, observable } from 'mobx';
import axios from 'axios';
import moment from 'moment';

const interval = 1000;
const bucketCount = 300;
const maxRequestsToKeep = 2000;

function isSelected (value = '', filterString = '') {
  return value.indexOf(filterString) >= 0;
}

export default class ProxiedRequests {
  @observable filter = '';
  @observable filtered = [];
  @observable requests = [];
  @observable requestsPerTimeBucket = [];
  @observable statusSeen = new Set();

  constructor() {
    const socket = new WebSocket(`ws://${window.location.host}/requests`);
    socket.onmessage = (message) => {
      this.addRequest(JSON.parse(message.data));
    };

    this.fetchRequests();
    setInterval(() => { this.calculateRequestsPerTimeBucket()}, interval);
  }

  @action
  addRequest(request) {
    if (isSelected(request.requestedURL, this.filter)) {
      this.filtered.push(request);
    }
    this.requests.push(request);
  }

  @action
  calculateRequestsPerTimeBucket() {
    const buckets = [];
    const now = Date.now();
    for (let time = now; time > now - (bucketCount * interval); time -= interval) {
      buckets.push({ start: time, startString: moment(time).format("HH:mm:ss.SSS"), end: time + interval, count: 0 });
    }
    buckets.reverse();

    let requests = this.filtered;
    if (requests.length > maxRequestsToKeep) {
      this.filtered = requests.slice(requests.length - maxRequestsToKeep);
      requests = this.filtered;
    }

    requests.forEach((r) => {
      for(let i = 0; i < buckets.length; i += 1) {
        const { start, end } = buckets[i];
        if (r.startTime >= start && r.startTime < end) {
          buckets[i].count += 1;
        }
      }
    });

    this.requestsPerTimeBucket = buckets;
  }

  @action
  fetchRequests() {
    return axios.get("/requests")
      .then(({data}) => this.setCollection(data));
  }

  @action
  setCollection(data) {
    this.requests = data;
    this.filtered = data.filter((v) => isSelected(v.requestedURL, this.filter));
  }

  @action
  setFilter(newValue) {
    this.filter = newValue;
    this.filtered = this.requests.filter((v) => isSelected(v.requestedURL, this.filter));
    this.calculateRequestsPerTimeBucket();
  }
}
