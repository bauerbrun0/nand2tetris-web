<script>
  import {
    hardwareSimulatorError,
    internalPins,
    outputPins,
  } from "../../store";
  import InputPins from "./InputPins/InputPins.svelte";
  import OutputPins from "./OutputPins/OutputPins.svelte";

  import TopBar from "./TopBar/TopBar.svelte";
</script>

<div class="h-full w-full overflow-auto p-2">
  {#if !$hardwareSimulatorError}
    <TopBar />
    <InputPins />
    <OutputPins pins={$outputPins} title="Output pins" />
    {#if $internalPins.length > 0}
      <OutputPins pins={$internalPins} title="Internal pins" />
    {/if}
  {:else}
    There is an error
  {/if}
</div>
