<script lang="ts">
  import type { Pin } from "../../../types";
  import { inputPins } from "../../../store";

  export let pin: Pin;
</script>

<div class="flex items-center gap-2">
  {#each pin.bits as bit, index (index)}
    <button
      onclick={() => {
        inputPins.update((pins) => {
          pins.forEach((p) => {
            if (p.name === pin.name) {
              p.bits[index] = p.bits[index] ? false : true;
            }
          });
          return pins;
        });
      }}
      class="bg-primary hover:bg-primary-800 h-7 cursor-pointer rounded px-2"
    >
      {bit ? "1" : "0"}
    </button>
  {/each}
</div>
