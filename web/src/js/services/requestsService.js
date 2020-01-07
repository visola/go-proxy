import { BehaviorSubject } from 'rxjs';

import webSocketService from './webSocketService';

class RequestsService extends BehaviorSubject {
  constructor() {
    super();
    this.requests = {};
    this.publish();

    webSocketService.subscribe((message) => {
      if (message == null || message.type != "request") {
        return;
      }

      const newRequest = message.data;
      this.requests[newRequest.id] = newRequest;
      this.publish();
    });
  }

  fetch() {
    this.next({loading: true, data: null});
    return fetch('/api/requests')
      .then((response) => response.json())
      .then((data) => {
        this.requests = data;
        this.publish();
      });
  }

  publish() {
    const dataArray = Object.values(this.requests);
    dataArray.sort((r1, r2) => {
      return r2.timings.started - r1.timings.started;
    });
    this.next({loading: false, data: dataArray});
  }
}

const requestsService = new RequestsService();
export default requestsService;
