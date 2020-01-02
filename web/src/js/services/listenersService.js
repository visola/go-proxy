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
    return fetch('/listeners')
      .then((response) => response.json())
      .then((data) => {
        this.data.loading = false;
        data.sort((l1, l2) => l1.configuration.port - l2.configuration.port);
        this.data.data = data;
        this.next(this.data);
      });
  }

  setEnabledUpstreams(listener, upstreamNames) {
    this.data.saving = true;
    const toUpdate = this.data.data.find((l) => l.configuration.port == listener.configuration.port);
    toUpdate.enabledUpstreams = upstreamNames;
    this.next(this.data);

    const options = {
      body: JSON.stringify(toUpdate.enabledUpstreams),
      headers: { 'Content-type': 'application/json'},
      method: 'PUT',
    };
    return fetch(`/listeners/${toUpdate.configuration.port}/upstreams`, options)
      .then((r) => r.json())
      .then((data) => {
        const newListener = data.listener;
        const indexOf = this.data.data.indexOf((l) => l.configuration.name == newListener.configuration.name);
        this.data.data[indexOf] = newListener;
        this.data.saving = false;
        this.next(this.data);
      })
  }
}

const listenersService = new ListenersService();
listenersService.fetch();
export default listenersService;
