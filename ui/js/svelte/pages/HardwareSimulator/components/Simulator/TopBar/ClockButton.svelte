<script lang="ts">
  import ClockIcon from "../../../../../components/icons/Clock.svelte";
  import { t } from "../../../../../../utils/i18n/i18n.ts";
  import {
    simulationRunning,
    simulationLoopRunning,
    cycleCount,
    cycleStage,
    advanceCycle,
  } from "../../../store.ts";

  const clockCount = $derived.by(() => {
    const count = $cycleCount;
    const stage = $cycleStage;
    const result: string = count + (stage === "tock" ? "+" : "\u2007");
    return result;
  });

  function handleClick() {
    simulationRunning.set(true);
    if ($cycleStage === "tick") {
      window.WASM.HardwareSimulator.tick();
    } else {
      window.WASM.HardwareSimulator.tock();
    }
    advanceCycle();
    simulationRunning.set(false);
  }
</script>

<button
  disabled={$simulationRunning || $simulationLoopRunning}
  onclick={handleClick}
  class={`
    dark:bg-silver-900 dark:hover:bg-silver-800 bg-white-700 hover:bg-white-900 disabled:dark:hover:bg-silver-900 disabled:hover:bg-white-700
    flex h-[44px] min-w-[140px] cursor-pointer items-center
    gap-2 rounded-md p-2 whitespace-nowrap disabled:cursor-not-allowed
  `}
>
  <div class="flex items-center gap-2">
    <ClockIcon classes="w-5 h-5 stroke-[1.5px]" />
    {t("hardware_simulator_page.clock")}:
  </div>
  <span class="w-full text-center">
    {clockCount}
  </span>
</button>
