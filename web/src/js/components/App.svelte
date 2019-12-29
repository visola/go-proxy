<script>
import { onDestroy } from 'svelte';

import { listeners, loading as loadingListeners, selectedListener } from '../stores/listenersStore';
import ListenersService from '../services/ListenersService';

import Listener from './Listener.svelte';

ListenersService.fetch();

function setSelection(event) {
  ListenersService.setSelected(event.target.value);
}
</script>

<div class="ui menu">
  <div class="header item">
    go-proxy
  </div>
</div>

{#if $loadingListeners || $listeners == null || $selectedListener == null}
  <p>Loading...</p>
{:else}
  {#if $listeners.length == 1}
    <label>Listener: </label>{$listeners[0].configuration.name}
  {:else}
    <label>Listener: </label>
    <select on:change={setSelection}>
      {#each $listeners as listener, index}
        <option value={index}>{listener.configuration.name} ({listener.configuration.port})</option>
      {/each}
    </select>
  {/if}

  <hr />
  <Listener listener={$selectedListener} />
{/if}
