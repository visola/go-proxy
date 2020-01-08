<script>
import  { onDestroy, onMount } from 'svelte';

import RequestDetails from './RequestDetails.svelte';
import Requests from './Requests.svelte';
import requestsService from '../services/requestsService';

let loadingRequests = false;
let requests = null;
let selectedRequest = null;

const subscription = requestsService.subscribe(({ loading, data }) => {
  loadingRequests = loading;
  if (data == null) {
    return;
  }

  requests = data;

  const selectedId = window.location.hash.substring(1);
  selectedRequest = requestsService.findById(selectedId);
});

function setSelectedRequest(id) {
  if (selectedRequest != null && selectedRequest.id == id) {
    history.pushState("", document.title, window.location.pathname + window.location.search);
    selectedRequest = null;
    return;  
  }

  window.location.hash = id;
  selectedRequest = requestsService.findById(id);
}

onDestroy(() => subscription.unsubscribe());

onMount(() => {
  requestsService.fetch();
});
</script>

{#if loadingRequests == true || requests == null}
  <p>Loading...</p>
{:else}
  <div class="requests-page">
    <Requests
      {requests}
      on:requestClicked={(e) => setSelectedRequest(e.detail)}
      selected={selectedRequest}
      short={selectedRequest != null}
    />
    {#if selectedRequest != null}
      <RequestDetails request={selectedRequest} />
    {/if}
  </div>
{/if}
