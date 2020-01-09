<script>
import { createEventDispatcher } from 'svelte';
import { elasticOut } from 'svelte/easing';

export let requests;
export let short = false;

const dispatch = createEventDispatcher();

const MAX_LENGTH = 100;
export let selected = {};

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

function handleRequestSelected(id) {
  dispatch('requestClicked', id);
}

function trimToMax(text) {
  if (text.length > MAX_LENGTH) {
    return text.substring(0, MAX_LENGTH - 3) + "...";
  }
  return text;
}
</script>

{#if requests.length == 0}
  <p style="margin-left: 10px;">Nothing here. Execute some requests to see the data.</p>
{:else}
  <table class="ui celled table requests collapsing">
    <thead>
      <tr>
        <th>Status Code</th>
        <th>Method</th>
        <th>Path</th>
        {#if !short}
          <th>Executed Path</th>
          <th>Timings</th>
        {/if}
      </tr>
    </thead>
    <tbody>
      {#each requests as request (request.id)}
        <tr
          class:selected={selected != null && request.id == selected.id}
          class="{getStatusClass(request.statusCode)}"
          on:click|preventDefault={() => handleRequestSelected(request.id)}
        >
          <td>
            {#if request.timings.completed == 0}
              <div class="ui active inline loader tiny"></div>
            {:else}
              {request.statusCode}
            {/if}
          </td>
          <td>{request.method}</td>
          <td title={request.url}>
            {trimToMax(request.url)}
          </td>
          {#if !short}
            <td title={request.executedURL}>
              {trimToMax(request.executedURL)}
            </td>
            <td>
              {#if request.timings.completed == 0}
                <div class="ui active inline loader tiny"></div>
              {:else}
                {Math.round((request.timings.completed - request.timings.started) / 1000000)}ms
              {/if}
            </td>
          {/if}
        </tr>
      {/each}
    </tbody>
  </table>
{/if}
