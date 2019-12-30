<script>
import { onDestroy } from 'svelte';

import Listener from './Listener.svelte';
import ListenerSelector from './ListenerSelector.svelte';

import listenersService from '../services/listenersService';
import upstreamsService from '../services/upstreamsService';

let listeners;
let loadingListeners = false;
let selectedListener;

const listenersSubscription = listenersService.subscribe(({ loading, data }) => {
  loadingListeners = loading;
  listeners = data;
});

onDestroy(() => {
  listenersService.unsubscribe(listenersSubscription);
});

function selectedListenerChanged(event) {
  selectedListener = event.detail;
}
</script>
<div class="ui menu">
  <div class="header item">
    go-proxy
  </div>
</div>

{#if loadingListeners || listeners == null }
  <p>Loading...</p>
{:else}
  <div class="header-justified">
    <Listener listener={selectedListener} />

    <div>
      <ListenerSelector listeners={listeners} on:changed={selectedListenerChanged} />
    </div>
  </div>
  <hr />
{/if}
