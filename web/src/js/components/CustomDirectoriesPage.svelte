<script>
import { onDestroy, onMount } from 'svelte';

import AddCustomDirectory from './AddCustomDirectory.svelte';
import customDirectoriesService from '../services/customDirectoriesService';

let customDirectories;

const customDirectoriesSubscription = customDirectoriesService.subscribe(({data, loading}) => {
  if (loading) {
    customDirectories = null;
    return;
  }

  customDirectories = data;
});

function removeCustomDirectory(toRemove) {
  customDirectoriesService.remove(toRemove);
}

onDestroy(() => customDirectoriesSubscription.unsubscribe());
onMount(() => customDirectoriesService.fetch());
</script>

{#if customDirectories == null}
  <p>Loading...</p>
{:else if customDirectories.length === 0}
  <AddCustomDirectory />
  <p>No custom directories configured.</p>
{:else}
  <AddCustomDirectory />
  <ul>
    {#each customDirectories as customDirectory}
      <li>
        {customDirectory}
        <i class="clickable ui icon alternate trash outline red" on:click={() => removeCustomDirectory(customDirectory)} />
      </li>
    {/each}
  </ul>
{/if}
