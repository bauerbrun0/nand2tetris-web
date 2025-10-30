<script lang="ts">
  import type { Pin } from "../../../types";
  import OutputPin from "./OutputPin.svelte";

  export let pins: Pin[];
  export let title: string;
</script>

<div class="p-2">
  <h2 class="my-2">{title}</h2>
  <div class="ml-2 flex flex-col gap-3">
    {#each pins.sort( (a, b) => a.name.localeCompare(b.name), ) as output (output.name)}
      <OutputPin pin={output} />
    {/each}
  </div>
</div>
