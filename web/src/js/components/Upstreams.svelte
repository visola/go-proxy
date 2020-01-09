<script>
import  { afterUpdate, onDestroy } from 'svelte';

import listenersService from '../services/listenersService';
import upstreamsService from '../services/upstreamsService';

export let listener;
let endpoints;
let upstreams;
let loadingUpstreams = false;
let savingListeners = false;

const listenersSubscription = listenersService.subscribe(({ saving }) => {
  savingListeners = saving;
});

const upstreamsSubscription = upstreamsService.subscribe(({ loading, data }) => {
  loadingUpstreams = loading;
  upstreams = data;
  updateEndpoints();
});  

afterUpdate(() => {
  updateEndpoints();
});

onDestroy(() => {
  listenersSubscription.unsubscribe();
  upstreamsSubscription.unsubscribe();
});

function updateEndpoints() {
  if (upstreams == null) {
    return;
  }

  endpoints = [];
  const addToEndpoints = (e) => {
    endpoints.push(e);
  };

  upstreams.forEach((u) => {
    u.proxyEndpoints.forEach(addToEndpoints);
    u.staticEndpoints.forEach(addToEndpoints);
  });

  endpoints.sort((e1, e2) => {
    const path1 = e1.from == "" ? e1.regexp : e1.from;
    const path2 = e2.from == "" ? e2.regexp : e2.from;

    const partCountDiff = path2.split("/").length - path1.split("/").length;
    if (partCountDiff != 0) {
      return partCountDiff;
    }

    const pathLengthDiff = path2.length - path1.length;
    if (pathLengthDiff != 0) {
      return pathLengthDiff;
    }

    return path1.localeCompare(path2);
  });
}

function upstreamSelected(checked, name) {
  if (checked) { // add to list
    listener.enabledUpstreams.push(name);
  } else { // remove from list
    const indexOf = listener.enabledUpstreams.indexOf(name);
    listener.enabledUpstreams.splice(indexOf, 1);
  }
  listenersService.setEnabledUpstreams(listener, listener.enabledUpstreams);
}
</script>

{#if listener != null && upstreams != null}
  <div class="ui segment">
    <table class="ui celled table">
      <thead>
        <tr>
          <th>Listener</th>
          <th>From or Regexp</th>
          <th>To</th>
        </tr>
      </thead>
      <tbody>
        {#each endpoints as endpoint}
          <tr>
            <td>
              <input
                checked={listener.enabledUpstreams.indexOf(endpoint.upstreamName) >= 0}
                on:change={(e) => upstreamSelected(e.target.checked, endpoint.upstreamName)}
                type="checkbox"
              />
              {endpoint.upstreamName}
            </td>
            <td>{endpoint.from == "" ? endpoint.regexp : endpoint.from}</td>
            <td>{endpoint.to}</td>
          </tr>
        {/each}
      </tbody>
    </table>
    <div class:active={savingListeners} class="ui dimmer">
      <div class="ui text loader">Saving...</div>
    </div>
  </div>
{/if}
