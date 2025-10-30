<script>
  import RefreshIcon from "../../../../../components/icons/Refresh.svelte";
  import { t } from "../../../../../../utils/i18n/i18n.ts";
  import {
    simulationRunning,
    simulationLoopRunning,
    resetCycle,
  } from "../../../store.ts";

  function reset() {
    simulationRunning.set(true);
    window.WASM.HardwareSimulator.processHdls();
    resetCycle();
    simulationRunning.set(false);
  }
</script>

<button
  disabled={$simulationRunning || $simulationLoopRunning}
  onclick={reset}
  class={`
    dark:bg-silver-900 dark:hover:bg-silver-800 bg-white-700 hover:bg-white-900
    disabled:dark:hover:bg-silver-900 disabled:hover:bg-white-700 flex h-[44px] cursor-pointer items-center gap-2
    rounded-md p-2 disabled:cursor-not-allowed
  `}
>
  <RefreshIcon classes="w-5 h-5 stroke-[1.5px]" />
  {t("hardware_simulator_page.reset")}
</button>
