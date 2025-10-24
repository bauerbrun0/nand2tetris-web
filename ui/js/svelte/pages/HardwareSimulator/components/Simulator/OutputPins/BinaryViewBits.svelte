<script lang="ts">
  import type { Pin } from "../../../types";

  export let pin: Pin;
</script>

<div class="flex items-center gap-2">
  {#each pin.bits as bit, index (index)}
    <div
      class="dark:bg-silver-800 bg-white-800 flex h-7 items-center rounded px-2"
    >
      {bit ? "1" : "0"}
    </div>
  {/each}
</div>
