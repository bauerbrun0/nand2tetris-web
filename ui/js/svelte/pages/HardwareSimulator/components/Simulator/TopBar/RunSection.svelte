<script lang="ts">
  import PlayIcon from "../../../../../components/icons/Play.svelte";
  import PauseIcon from "../../../../../components/icons/Pause.svelte";
  import HertzDropdown from "./HertzDropdown.svelte";
  import { t } from "../../../../../../utils/i18n/i18n.ts";
  import {
    automaticSimulationRunning,
    simulationRunning,
  } from "../../../store.ts";

  function handleRunClick(e: Event) {
    e.preventDefault();
    automaticSimulationRunning.set(true);
  }

  function handlePauseClick(e: Event) {
    e.preventDefault();
    automaticSimulationRunning.set(false);
  }
</script>

<div
  class="dark:bg-silver-900 bg-white-700 flex h-[44px] items-center gap-1 rounded-md p-2"
>
  <HertzDropdown />
  {#if !$automaticSimulationRunning}
    <button
      disabled={$simulationRunning}
      onclick={handleRunClick}
      class={`
        dark:bg-silver-800 dark:hover:bg-silver-700 hover:bg-white-400 disabled:dark:hover:bg-silver-800 disabled:hover:bg-white-400
        flex h-[32px] cursor-pointer items-center gap-2 rounded-md bg-white
        px-4 py-1 disabled:cursor-not-allowed
      `}
    >
      <PlayIcon classes="w-4 h-4 stroke-[1.5px]" />

      {t("hardware_simulator_page.run")}
    </button>
  {:else}
    <button
      onclick={handlePauseClick}
      class="dark:bg-silver-800 dark:hover:bg-silver-700 hover:bg-white-400 flex h-[32px] cursor-pointer items-center gap-2 rounded-md bg-white px-4 py-1"
    >
      <PauseIcon classes="w-4 h-4 stroke-[1.5px]" />

      {t("hardware_simulator_page.pause")}
    </button>
  {/if}
</div>
