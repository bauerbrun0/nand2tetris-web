<script lang="ts">
  import { currentProjectName, currentHdlFileName } from "../store";
  import ChevronIcon from "../../../components/icons/Chevron.svelte";
  import { t } from "../../../../utils/i18n/i18n";
</script>

<div class="flex h-[48px] w-full items-center py-[8px]">
  <a
    data-tooltip-target="tooltip-back-to-projects"
    data-tooltip-placement="bottom"
    href="/projects"
    class="hover:dark:bg-silver-900 hover:bg-silver-100 rounded"
  >
    <ChevronIcon classes="h-7 w-7" />
  </a>
  <div
    id="tooltip-back-to-projects"
    role="tooltip"
    class="tooltip dark:bg-silver-800 bg-silver-100 invisible absolute z-10 inline-block rounded-lg px-3 py-2 text-sm font-medium opacity-0 shadow-xs"
  >
    {t("hardware_simulator_page.back_to_projects")}
    <div class="tooltip-arrow" data-popper-arrow></div>
  </div>

  <h1 class="ml-4 font-bold">
    {$currentProjectName} / {$currentHdlFileName}.hdl
  </h1>
</div>
