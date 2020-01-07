<script>
import  { onDestroy, onMount } from 'svelte';

import requestsService from '../services/requestsService';

let loadingRequests = false;
let requests = null;

const subscription = requestsService.subscribe(({ loading, data }) => {
  loadingRequests = loading;
  if (data == null) {
    return;
  }

  requests = data;
});

function getStatusClass(statusCode) {
  if (statusCode == 0) {
    return "status_loading";
  }

  if (statusCode < 300) {
    return "status_success";
  }

  if (statusCode < 400) {
    return "status_redirect";
  }

  if (statusCode < 500) {
    return "status_client_error";
  }

  return "status_server_error";
}

onDestroy(() => subscription.unsubscribe());

onMount(() => {
  requestsService.fetch();
});
</script>

{#if loadingRequests == true || requests == null}
  <p>Loading...</p>
{:else}
  <table class="ui celled table">
    <thead>
      <tr>
        <th>Status Code</th>
        <th>Path</th>
        <th>Timings</th>
      </tr>
    </thead>
    <tbody>
      {#each requests as request }
        <tr class="{getStatusClass(request.statusCode)}">
          <td>
            {#if request.timings.completed == 0}
              <div class="ui active inline loader tiny"></div>
            {:else}
              {request.statusCode}
            {/if}
          </td>
          <td>
            {#if request.url.length >= 50}
              {request.url.substring(0, 47)}...
            {:else}
              {request.url}
            {/if}
          <td>
            {#if request.timings.completed == 0}
              <div class="ui active inline loader tiny"></div>
            {:else}
              {Math.round((request.timings.completed - request.timings.started) / 1000000)}ms
            {/if}
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
{/if}
