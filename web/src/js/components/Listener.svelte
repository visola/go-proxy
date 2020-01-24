<script>
export let listener;

function getUrl(listener) {
  if (isHttps(listener)) {
    return `https://localhost:${listener.port}`;
  }
  return `http://localhost:${listener.port}`;
}

function isHttps(listener) {
  return listener.certificateFile != "" || listener.keyFile != "";
}
</script>

{#if listener == null}
  <!-- Nothing -->
{:else}
  <div>
    {listener.name}
    <a href={getUrl(listener)} target="_blank"><i class="external alternate icon"></i></a>
    <div class="ui label large">Port<div class="detail">{listener.port}</div></div>
    <div class="ui label large {isHttps(listener) ? 'green' : 'orange'}">
      HTTPS
      <div class="detail"><i class="{isHttps(listener) ? 'lock' : 'unlock'} alternate icon"></i></div>
    </div>
  </div>
{/if}