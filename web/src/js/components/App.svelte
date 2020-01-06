<script>
import ListenersPage from './ListenersPage.svelte';
import RequestsPage from './RequestsPage.svelte';
import Route from './Route.svelte';
import TopMenu from './TopMenu.svelte';
</script>

<TopMenu />
<Route path="/"><RequestsPage /></Route>
<Route path="/listeners"><ListenersPage /></Route>
<Route path="/requests"><RequestsPage /></Route>
