<script>
import { createEventDispatcher, onMount } from 'svelte';

const dispatch = createEventDispatcher();

export let listeners;
let selectedListener = null;

onMount(() => {
  selectedListener = listeners[0];
  dispatch('changed', selectedListener);
});

function setSelection(event) {
  selectedListener = listeners[event.target.value];
  dispatch('changed', selectedListener);
}
</script>

<div>
  <label>Listener: </label>
  {#if listeners.length == 1}
    {listeners[0].name}
  {:else}
    <select on:change={setSelection}>
      {#each listeners as listener, index}
        <option value={index}>{listener.name} ({listener.port})</option>
      {/each}
    </select>
  {/if}
</div>