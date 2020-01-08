<script>
import { afterUpdate } from 'svelte';

const TEXT_CONTENT_TYPES = [
  'text/css',
  'text/html',
  'text/plain',
  'application/x-javascript',
  'application/javascript',
  'application/json'
];

export let body;
export let headers;

let decodedBody;
let forceDecodedBody;
let type;

function forceDecodeBody() {
  forceDecodedBody = atob(body);
}

afterUpdate(() => {
  decodedBody = null;
  type = "unknown";
  if (headers == null) {
    return;
  }

  const contentTypeHeader = Object.keys(headers)
    .find((n) => n.toLowerCase() == 'content-type');

  if (contentTypeHeader) {
    type = headers[contentTypeHeader][0];
    const mime = type.split(";")[0];
    if (TEXT_CONTENT_TYPES.indexOf(mime) >= 0) {
      decodedBody = atob(body);
    }
  }
});
</script>

<form class="ui form">
  <div class="field">
    <label>Content Type:</label>
    {type}
  </div>
  <div class="field">
    {#if decodedBody}
      <label>Content:</label>
      <textarea>{decodedBody}</textarea>
    {:else}
      <label>
        Base64 Encoded Bytes:
        <button on:click|preventDefault={forceDecodeBody}>Decode</button>
      </label>
      {#if forceDecodedBody == null}
        <textarea>{body}</textarea>
      {:else}
        <textarea>{forceDecodedBody}</textarea>
      {/if}
    {/if}
  </div>
</form>
