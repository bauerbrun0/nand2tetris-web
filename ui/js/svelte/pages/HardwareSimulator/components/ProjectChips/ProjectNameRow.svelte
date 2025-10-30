<script lang="ts">
  import PlusIcon from "../../../../components/icons/Plus.svelte";
  import { t } from "../../../../../utils/i18n/i18n";
  import { currentProjectName, newChipNameInputOpen } from "../../store";
</script>

<div class="flex items-center justify-between px-4 py-1">
  <span>{$currentProjectName}</span>
  <button
    onclick={() => {
      newChipNameInputOpen.set(true);
    }}
    data-tooltip-target="tooltip-add-chip"
    data-tooltip-placement="bottom"
    class="dark:hover:bg-silver-800 flex cursor-pointer items-center justify-center rounded p-1"
  >
    <PlusIcon classes="h-4 w-4 stroke-[1.5px]" />
  </button>
  <div
    id="tooltip-add-chip"
    role="tooltip"
    class="tooltip dark:bg-silver-800 bg-silver-100 invisible absolute z-10 inline-block rounded-lg px-3 py-2 text-sm opacity-0 shadow-xs"
  >
    {t("hardware_simulator_page.add_new_chip")}
    <div class="tooltip-arrow" data-popper-arrow></div>
  </div>
</div>
