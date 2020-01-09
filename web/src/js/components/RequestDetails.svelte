<script>
import { afterUpdate } from 'svelte';

import BodyDetails from './BodyDetails.svelte';
import HeadersForm from './HeadersForm.svelte';

export let request;

let selectedTab = 0;
let requestId;

afterUpdate(() => {
  if (requestId != request.id) {
    requestId = request.id;
    selectedTab = 0
  }
}); 
</script>
<div class="request-detail">
  <form class="ui form">
    <h4 class="ui dividing header">Request</h4>
    <div class="three fields">
      <div class="field">
        <label>Status Code:</label>
        {request.statusCode}
      </div>
      <div class="field">
        <label>Method:</label>
        {request.method}
      </div>
      <div class="field">
        <label>Internal ID:</label>
        {request.id}
      </div>
    </div>
    <div class="field">
      <label>Requested URL:</label>
      {request.url}
    </div>
    <div class="field">
      <label>Executed URL:</label>
      {request.executedURL}
    </div>
    <div class="ui pointing menu">
      <a class="item" class:active={selectedTab == 0} href="#0" on:click|preventDefault={() => selectedTab = 0}>
        Request Headers
      </a>
      {#if request.request.body != ""}
        <a class="item" class:active={selectedTab == 5} href="#5" on:click|preventDefault={() => selectedTab = 5}>
          Request Body
        </a>
      {/if}
      <a class="item" class:active={selectedTab == 10} href="#10" on:click|preventDefault={() => selectedTab = 10}>
        Response Headers
      </a>
      {#if request.response.body != ""}
        <a class="item" class:active={selectedTab == 20} href="#20" on:click|preventDefault={() => selectedTab = 20}>
          Response Body
        </a>
      {/if}
    </div>
    <div class="ui segment">
      {#if selectedTab == 0}
        <HeadersForm headers={request.request.headers} />
      {:else if selectedTab == 5} 
        <BodyDetails body={request.request.body} headers={request.request.headers} />
      {:else if selectedTab == 10} 
        <HeadersForm headers={request.response.headers} />
      {:else if selectedTab == 20} 
        <BodyDetails body={request.response.body} headers={request.response.headers} />
      {/if}
    </div>
  </form>
  
</div>
