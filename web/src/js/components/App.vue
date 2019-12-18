<template>
  <div v-if="loadingListeners">Loading...</div>  
  <div v-else>
    <h2>Listeners</h2>
    <ul>
      <li v-for="listener in listeners" v-bind:key="listener.configuration.port">
        {{ listener.configuration.port }}: {{ listener.enabledUpstreams }}
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  created() {
    this.loadListeners();
  },

  data() {
    return {
      loadingListeners: false,
      listeners: [],
    }
  },

  methods: {
    loadListeners: function () {
      this.loadingListeners = true;
      fetch('/listeners')
        .then((response) => response.json())
        .then((listeners) => {
          this.listeners = Object.values(listeners);
          this.loadingListeners = false;
        });
    }
  },
}
</script>