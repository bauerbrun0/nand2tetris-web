<script lang="ts">
  import type { Pin } from "../../../types";
  import BinaryViewBits from "./BinaryViewBits.svelte";

  export let pin: Pin;
</script>

<div class="flex items-center gap-3">
  <h3>{pin.name}</h3>
  <div>
    <BinaryViewBits {pin} />
  </div>
</div>
