<script>
import { onDestroy } from 'svelte';
import Link from './Link.svelte';

import routingService from '../services/routingService';

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
  <Link class="item {currentPath == "/" || currentPath == "/requests" ? 'active' : ''}" href="/requests">Requests</Link>
  <Link class="item {currentPath == "/listeners" ? 'active' : ''}" href="/listeners">Listeners</Link>
</div>
