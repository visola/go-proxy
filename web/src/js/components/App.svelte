<script>
import CustomDirectoriesPage from './CustomDirectoriesPage.svelte';
import ListenersPage from './ListenersPage.svelte';
import RequestsPage from './RequestsPage.svelte';
import Route from './Route.svelte';
import TopMenu from './TopMenu.svelte';

const routes = [{
  component: RequestsPage,
  label: "Requests",
  paths: ["/", "/requests"],
},{
  component: ListenersPage,
  label: "Listeners",
  paths: ["/listeners"],
},{
  component: CustomDirectoriesPage,
  label: "Custom Directories",
  paths: ["/custom-directories"],
}];
</script>

<TopMenu {routes} />

{#each routes as route}
  {#each route.paths as path}
    <Route path={path}><svelte:component this={route.component} /></Route>
  {/each}
{/each}
