<script>
export let headers = {};
</script>

{#each Object.entries(headers || {}) as [name, values]}
  <div class="field">
    <label>{name}</label>
    {values.join(",")}
  </div>
{:else}
  <label>No headers</label>
{/each}