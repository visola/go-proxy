<script>
import { onDestroy } from 'svelte';

import Listener from './Listener.svelte';
import ListenerSelector from './ListenerSelector.svelte';
import Upstreams from './Upstreams.svelte';

import listenersService from '../services/listenersService';

let listeners;
let loadingListeners = false;
let selectedListener;

const listenersSubscription = listenersService.subscribe(({ loading, data }) => {
  loadingListeners = loading;
  listeners = data;
});

onDestroy(() => {
  listenersSubscription.unsubscribe();
});

function selectedListenerChanged(event) {
  selectedListener = event.detail;
}
</script>

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
  <Upstreams listener={selectedListener} />
{/if}
