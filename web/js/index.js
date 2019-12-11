let data = {
  listeners: [],
  loadingListeners: false,
};

let vm = new Vue({
  el: '#app',
  data
});

data.loadingListeners = true;
fetch('/listeners')
  .then((resp) => {
    data.loadingListeners = false;
    return resp.json();
  }).then((listeners) => data.listeners = listeners);