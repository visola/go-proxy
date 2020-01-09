<script>
import { afterUpdate, createEventDispatcher } from 'svelte';

import BodyDetails from './BodyDetails.svelte';
import HeadersForm from './HeadersForm.svelte';

export let request;

let selectedTab = 0;
let requestId;

const dispatch = createEventDispatcher();

function handleClose () {
  dispatch('close');
}

afterUpdate(() => {
  if (requestId != request.id) {
    requestId = request.id;
    selectedTab = 0
  }
});
</script>
<div class="request-detail">
  <form class="ui form">
    <h4 class="ui dividing header">Request <i class="button close icon fitted" on:click|preventDefault={handleClose}></i></h4>
    <div class="three fields">
      <div class="field">
        <label>Status Code:</label>
        <input type="text" readonly value={request.statusCode} />
      </div>
      <div class="field">
        <label>Method:</label>
        <input type="text" read value={request.method} />
      </div>
      <div class="field">
        <label>Internal ID:</label>
        <input type="text" readonly value={request.id} />
      </div>
    </div>
    <div class="field">
      <label>Requested URL:</label>
      <input type="text" readonly value={request.url} />
    </div>
    <div class="field">
      <label>Executed URL:</label>
      <input type="text" readonly value={request.executedURL} />
    </div>

    <div class="ui pointing menu">
      <a class="item" class:active={selectedTab == 0} href="#0" on:click|preventDefault={() => selectedTab = 0}>
        Response Headers
      </a>
      {#if request.response.body != ""}
        <a class="item" class:active={selectedTab == 1} href="#1" on:click|preventDefault={() => selectedTab = 1}>
          Response Body
        </a>
      {/if}
      <a class="item" class:active={selectedTab == 2} href="#2" on:click|preventDefault={() => selectedTab = 2}>
        Request Headers
      </a>
      {#if request.request.body != ""}
        <a class="item" class:active={selectedTab == 3} href="#3" on:click|preventDefault={() => selectedTab = 3}>
          Request Body
        </a>
      {/if}
    </div>

    <div class="ui segment">
      {#if selectedTab == 0}
        <HeadersForm headers={request.response.headers} />
      {:else if selectedTab == 1} 
        <BodyDetails body={request.response.body} headers={request.response.headers} />
      {:else if selectedTab == 2} 
        <HeadersForm headers={request.request.headers} />
      {:else if selectedTab == 3} 
        <BodyDetails body={request.request.body} headers={request.request.headers} />
      {/if}
    </div>
  </form>
  
</div>
