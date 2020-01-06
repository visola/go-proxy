import { BehaviorSubject } from 'rxjs';

class RequestsService extends BehaviorSubject {
  constructor() {
    super();
    this.next({loading: false, data: null});
  }

  fetch() {
    this.next({loading: true, data: null});
    return fetch('/api/requests')
      .then((response) => response.json())
      .then((data) => {
        this.next({loading: false, data});
      });
  }
}

const requestsService = new RequestsService();
export default requestsService;
