<script>
import customDirectoriesService from '../services/customDirectoriesService';

const FILE_PARTS = /([\/\\][^\/\\]+)/g;
let canAdd = false;
let toAdd = "";

function addCustomDirectory() {
  if (!canAdd) {
    return;
  }

  customDirectoriesService.add(toAdd);
}

function handlePathChanged(e) {
  toAdd = e.target.value;
  canAdd = toAdd.length > 3 && toAdd.match(FILE_PARTS).length > 2;
  if (e.keyCode === 13) {
    addCustomDirectory();
  }
}
</script>

<div class="ui segment compact">
  <div class="ui top attached label">
    Add Custom Directory
  </div>
  <div class="ui action input">
    <input type="text" placeholder="/path/to/directory" on:keyup={handlePathChanged} />
    <button disabled={!canAdd} class="ui button" on:click|preventDefault={addCustomDirectory}>
      Add
    </button>
  </div>
</div>
