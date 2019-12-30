<script>
export let listener;

function getUrl(listener) {
  if (isHttps(listener)) {
    return `https://localhost:${listener.configuration.port}`;
  }
  return `http://localhost:${listener.configuration.port}`;
}

function isHttps(listener) {
  return listener.configuration.certificateFile != "" || listener.configuration.keyFile != "";
}
</script>

{#if listener == null}
  <!-- Nothing -->
{:else}
  <div>
    {listener.configuration.name}
    <a href={getUrl(listener)} target="_blank"><i class="external alternate icon"></i></a>
    <div class="ui label large">Port<div class="detail">{listener.configuration.port}</div></div>
    <div class="ui label large {isHttps(listener) ? 'green' : 'orange'}">
      HTTPS
      <div class="detail"><i class="{isHttps(listener) ? 'lock' : 'unlock'} alternate icon"></i></div>
    </div>
  </div>
{/if}