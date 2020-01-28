import { BehaviorSubject } from 'rxjs';

class CustomDirectoriesService extends BehaviorSubject {
  constructor() {
    super();
    this.data = {
      data: null,
      loading: false,
      saving: false,
    };
    this.next(this.data);
  }

  add(toAdd) {
    this.data.saving = true;
    this.next(this.data);

    const options = {
      body: JSON.stringify(toAdd),
      headers: { 'Content-type': 'application/json'},
      method: 'POST',
    };

    return fetch('/api/custom-directories', options)
      .then(() => {
        this.data.saving = false;
        this.data.data.push(toAdd);
        this.next(this.data);
        this.fetch();
      });
  }

  fetch() {
    this.data.loading = true;
    this.data.data = null;
    this.next(this.data);
    return fetch('/api/custom-directories')
      .then((response) => response.json())
      .then((data) => {
        this.data.loading = false;
        data.sort((c1, c2) => c1.localeCompare(c2));
        this.data.data = data;
        this.next(this.data);
      });
  }

  remove(toRemove) {
    this.data.saving = true;
    this.next(this.data);

    const options = {
      body: JSON.stringify(toRemove),
      headers: { 'Content-type': 'application/json'},
      method: 'DELETE',
    };

    return fetch('/api/custom-directories', options)
      .then(() => {
        this.data.saving = false;
        const indexOf = this.data.data.indexOf(toRemove);
        this.data.data.splice(indexOf, 1);
        this.next(this.data);
        this.fetch();
      });
  }
}

const customDirectoriesService = new CustomDirectoriesService();
export default customDirectoriesService;
