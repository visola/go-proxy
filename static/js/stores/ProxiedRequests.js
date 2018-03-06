import { action, observable } from 'mobx';
import axios from 'axios';
import moment from 'moment';

const interval = 1000;
const bucketCount = 600;
const maxRequestsToKeep = 2000;

export default class ProxiedRequests {
  @observable requests = [];
  @observable requestsPerTimeBucket = [];
  @observable statusSeen = new Set();

  constructor() {
    const socket = new WebSocket("ws://localhost:1234/requests");
    socket.onmessage = (message) => {
      this.addRequest(JSON.parse(message.data));
    };

    this.fetchRequests();
    setInterval(() => { this.calculateRequestsPerTimeBucket()}, interval);
  }

  @action
  addRequest(request) {
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

    if (this.requests.length > maxRequestsToKeep) {
      this.requests = this.requests.slice(this.requests.length - maxRequestsToKeep);
    }

    this.requests.forEach((r) => {
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
      .then(({data}) => {
        this.requests = data;
      });
  }
}
