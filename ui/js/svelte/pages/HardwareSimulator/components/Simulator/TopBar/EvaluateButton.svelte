<script>
  import CalculatorIcon from "../../../../../components/icons/Calculator.svelte";
  import { t } from "../../../../../../utils/i18n/i18n.ts";
  import {
    simulationRunning,
    automaticSimulationRunning,
  } from "../../../store.ts";

  async function evaluate() {
    simulationRunning.set(true);
    window.WASM.HardwareSimulator.evaluate();
    simulationRunning.set(false);
  }
</script>

<button
  disabled={$simulationRunning || $automaticSimulationRunning}
  onclick={evaluate}
  class={`
    dark:bg-silver-900 dark:hover:bg-silver-800 bg-white-700 hover:bg-white-900
    disabled:dark:hover:bg-silver-900 disabled:hover:bg-white-700 flex h-[44px] cursor-pointer items-center gap-2
    rounded-md p-2 disabled:cursor-not-allowed
  `}
>
  <CalculatorIcon
    classes="w-5 h-5 stroke-text dark:stroke-text-dark stroke-[1.5px]"
  />
  {t("hardware_simulator_page.evaluate")}
</button>
