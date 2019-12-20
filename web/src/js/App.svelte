<script>
import { onMount } from 'svelte';

let loadingListeners = false;
let listeners = [];

onMount(() => {
  loadingListeners = true;
  fetch('/listeners')
    .then((response) => response.json())
    .then((newListeners) => {
      loadingListeners = false;
      listeners = Object.values(newListeners);
    });
});
</script>

{#if loadingListeners}
  <p>Loading...</p>
{:else}
  <h2>Listeners</h2>
  <ul>
  {#each listeners as listener}
    <li>{listener.configuration.port}: {listener.enabledUpstreams.join(',')}</li>
  {/each}
  </ul>
{/if}
