import { listeners, loading, selectedListener } from '../stores/listenersStore';

let allListeners;
listeners.subscribe(v => allListeners = v);

export default {
  fetch() {
    loading.set(true);
    return fetch('/listeners')
      .then((response) => response.json())
      .then((data) => {
        loading.set(false);
        listeners.set(data);
        selectedListener.set(data[0]);
      });
  },

  setSelected(index) {
    if (index >= 0) {
      selectedListener.set(allListeners[index]);
    }
  }
}
