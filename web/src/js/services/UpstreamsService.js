import { upstreams, loading } from '../stores/upstreamsStore';

export default {
  fetch() {
    loading.set(true);
    return fetch('/upstreams')
      .then((response) => response.json())
      .then((data) => {
        loading.set(false);
        upstreams.set(data);
      });
  },
};