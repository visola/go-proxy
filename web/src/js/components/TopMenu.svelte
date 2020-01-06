<script>
import { onDestroy } from 'svelte';
import Link from './Link.svelte';

import routingService from '../services/routingService';

export let routes;

let currentPath;
const subscription = routingService.subscribe((newPath) => {
  currentPath = newPath;
});

onDestroy(() => {
  subscription.unsubscribe();
});
</script>

<div class="ui menu">
  <div class="header item">go-proxy</div>
  {#each routes as route}
    <Link class="item {route.paths.indexOf(currentPath) >= 0 ? 'active' : ''}" href={route.paths[0]} >{route.label}</Link>
  {/each}
</div>
