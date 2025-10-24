<script>
  import ClockIcon from "../../../../../components/icons/Clock.svelte";
  import { t } from "../../../../../../utils/i18n/i18n.ts";
  import {
    simulationRunning,
    automaticSimulationRunning,
  } from "../../../store.ts";

  const clockCount = $state("43+");
</script>

<button
  disabled={$simulationRunning || $automaticSimulationRunning}
  onclick={() => console.log("clock clicked")}
  class={`
    dark:bg-silver-900 dark:hover:bg-silver-800 bg-white-700 hover:bg-white-900 disabled:dark:hover:bg-silver-900 disabled:hover:bg-white-700
    flex h-[44px] cursor-pointer items-center gap-2 rounded-md
    p-2 whitespace-nowrap disabled:cursor-not-allowed
  `}
>
  <ClockIcon classes="w-5 h-5 stroke-[1.5px]" />
  {t("hardware_simulator_page.clock")}: {clockCount}
</button>
