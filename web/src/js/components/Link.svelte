<script>
import routingService from '../services/routingService';

export let href = "";

let className = "";
export { className as class };

function onClick() {
  routingService.goTo(href);
}
</script>

<a class={className} {href} on:click|preventDefault={onClick}>
  <slot></slot>
</a>
