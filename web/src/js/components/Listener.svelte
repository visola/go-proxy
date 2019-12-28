<script>
  import { selectedListener } from '../stores/listenersStore';


  function getUrl(listener) {
    if (listener.configuration.certificateFile == "" || listener.configuration.keyFile == "") {
      return `http://localhost:${listener.configuration.port}`;
    }
    return `https://localhost:${listener.configuration.port}`;
  }
</script>

<h3>{$selectedListener.configuration.name}</h3>

<a href={getUrl($selectedListener)}>Go to <i class="external alternate icon"></i></a><br />
<label>Port: </label>{$selectedListener.configuration.port}<br />
<label>Certificate: </label>{$selectedListener.configuration.certificateFile}<br />
<label>Key: </label>{$selectedListener.configuration.keyFile}<br />
