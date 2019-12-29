import { BehaviorSubject } from 'rxjs';

class ListenersService extends BehaviorSubject {
  constructor() {
    super();
    this.next({loading: false, data: null});
  }

  fetch() {
    this.next({loading: true, data: null});
    return fetch('/listeners')
      .then((response) => response.json())
      .then((data) => {
        this.next({loading: false, data});
      });
  }
}

export default new ListenersService();
