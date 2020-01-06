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

  const dataArray = Object.values(data);
  dataArray.sort((r1, r2) => {
    return r2.timings.started - r1.timings.started;
  });

  requests = dataArray;
});

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
        <tr>
          <td>{request.timings.completed == 0 ? '-' : request.statusCode}</td>
          <td>{request.url}</td>
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
