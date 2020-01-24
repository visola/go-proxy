import { BehaviorSubject } from 'rxjs';

class ListenersService extends BehaviorSubject {
  constructor() {
    super();
    this.data = {
      data: null,
      loading: false,
      saving: false,
    };
    this.next(this.data);
  }

  fetch() {
    this.data.loading = true;
    this.data.data = null;
    this.next(this.data);
    return fetch('/api/listeners')
      .then((response) => response.json())
      .then((data) => {
        this.data.loading = false;
        data.sort((l1, l2) => l1.port - l2.port);
        this.data.data = data;
        this.next(this.data);
      });
  }

  save(listener) {
    this.data.saving = true;
    this.next(this.data);

    const options = {
      body: JSON.stringify(listener),
      headers: { 'Content-type': 'application/json'},
      method: 'PUT',
    };
    return fetch(`/api/listeners/${listener.name}`, options)
      .then((r) => r.json())
      .then((data) => {
        const newListener = data;
        const indexOf = this.data.data.indexOf((l) => l.name == newListener.name);
        this.data.data[indexOf] = newListener;
        this.data.saving = false;
        this.next(this.data);
      });
  }
}

const listenersService = new ListenersService();
listenersService.fetch();
export default listenersService;
