import { BehaviorSubject } from 'rxjs';

class UpstreamsService extends BehaviorSubject {
  constructor() {
    super();
    this.next({loading: false, data: null});
  }

  fetch() {
    this.next({loading: true, data: null});
    return fetch('/api/upstreams')
      .then((response) => response.json())
      .then((data) => {
        this.next({loading: false, data});
      });
  }
}

const upstreamsService = new UpstreamsService();
upstreamsService.fetch();
export default upstreamsService;
