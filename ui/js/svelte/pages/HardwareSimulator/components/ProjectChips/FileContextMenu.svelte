<script lang="ts">
  import EditIcon from "../../../../components/icons/Edit.svelte";
  import TrashIcon from "../../../../components/icons/Trash.svelte";
  import { t } from "../../../../../utils/i18n/i18n.ts";

  export let position: { x: number; y: number };
  export let clickedFileName: string;
</script>

<div
  class="bg-white-300 dark:bg-silver-900 flex h-[72px] w-[100px] flex-col"
  style="position: absolute; top:{position.y}px; left:{position.x}px"
>
  <button
    onclick={() => {
      console.log("Rename clicked ", clickedFileName);
    }}
    class={`
      dark:hover:bg-silver-800 stroke-text dark:stroke-text-dark hover:bg-white-800
      flex w-full cursor-pointer items-center justify-between p-2 text-left text-sm
    `}
  >
    <span>{t("hardware_simulator_page.rename_chip")}</span>
    <EditIcon classes="h-4 w-4 ml-2" />
  </button>
  <button
    onclick={() => {
      console.log("Delete clicked ", clickedFileName);
    }}
    class={`
      dark:hover:bg-silver-800 stroke-text dark:stroke-text-dark hover:bg-white-800
      flex w-full cursor-pointer items-center justify-between p-2 text-left text-sm hover:stroke-red-500
    `}
  >
    <span>{t("hardware_simulator_page.delete_chip")}</span>
    <TrashIcon classes="h-4 w-4 ml-2" />
  </button>
</div>
