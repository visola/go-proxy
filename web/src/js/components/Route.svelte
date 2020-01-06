<script>
import { onDestroy } from 'svelte';
import routingService from '../services/routingService';

export let path;
let currentPath;

const subscription = routingService.subscribe((newPath) => {
  currentPath = newPath;
});

onDestroy(() => {
  subscription.unsubscribe();
});
</script>

{#if currentPath == path }
<slot></slot>
{/if}
