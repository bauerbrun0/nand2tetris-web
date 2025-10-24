<script lang="ts">
  import { inputPins } from "../../../store";
  import InputPin from "./InputPin.svelte";
</script>

<div class="p-2">
  <h2 class="my-2">Input pins</h2>
  <div class="ml-2 flex flex-col gap-3">
    {#each $inputPins as input (input.name)}
      <InputPin pin={input} />
    {/each}
  </div>
</div>
