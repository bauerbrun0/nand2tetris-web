<script lang="ts">
  import { derived, readable, writable } from "svelte/store";
  import { hdls, currentHdlFileName, currentProjectName } from "../../store";
  import EditIcon from "../../../../components/icons/Edit.svelte";
  import TrashIcon from "../../../../components/icons/Trash.svelte";
  import PlusIcon from "../../../../components/icons/Plus.svelte";
  import { onMount } from "svelte";
  import { t } from "../../../../../utils/i18n/i18n";

  const sortedHdlFileNames = derived(hdls, ($hdls) =>
    Object.keys($hdls).sort((a, b) => a.localeCompare(b)),
  );

  let containerElement: HTMLDivElement;

  onMount(() => {
    containerElement.onclick = closeContextMenu;
  });

  type Dimensions = {
    width: number;
    height: number;
  };

  type Position = {
    x: number;
    y: number;
  };

  const containerDimensions = writable<Dimensions>({ width: 0, height: 0 });
  const fileContextMenuPosition = writable<Position>({ x: 0, y: 0 });
  const fileContextMenuDimensions = readable<Dimensions>({
    width: 100,
    height: 72,
  });
  const showFileContextMenu = writable<boolean>(false);
  const fileContextMenuClickedFileName = writable<string>("");

  function closeContextMenu() {
    showFileContextMenu.set(false);
  }

  function handleContextMenu(event: MouseEvent, chipFileName: string) {
    event.preventDefault();

    containerDimensions.set({
      width: containerElement.clientWidth,
      height: containerElement.clientHeight,
    });

    const absoluteClickPosition: Position = {
      x: event.clientX,
      y: event.clientY,
    };

    // Calculate position relative to container
    const relativeClickPosition: Position = {
      x:
        absoluteClickPosition.x - containerElement.getBoundingClientRect().left,
      y: absoluteClickPosition.y - containerElement.getBoundingClientRect().top,
    };

    // Calculate context menu position
    let fileContextMenuX: number = relativeClickPosition.x;
    let fileContextMenuY: number = relativeClickPosition.y;

    if (
      $containerDimensions.height - fileContextMenuY <
      $fileContextMenuDimensions.height
    ) {
      fileContextMenuY = fileContextMenuY - $fileContextMenuDimensions.height;
    }
    if (
      $containerDimensions.width - fileContextMenuX <
      $fileContextMenuDimensions.width
    ) {
      fileContextMenuX = fileContextMenuX - $fileContextMenuDimensions.width;
    }

    fileContextMenuPosition.set({
      x: fileContextMenuX,
      y: fileContextMenuY,
    });
    showFileContextMenu.set(true);
    fileContextMenuClickedFileName.set(chipFileName);
  }
</script>

<div class="relative h-full w-full overflow-auto" bind:this={containerElement}>
  <div class="flex items-center justify-between px-4 py-1">
    <span>{$currentProjectName}</span>
    <button
      onclick={() => {
        console.log("Add new chip clicked");
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
  {#each $sortedHdlFileNames as name (name)}
    <div
      class="dark:hover:bg-silver-800 hover:bg-white-800 flex w-full items-center justify-between"
      class:dark:bg-silver-800={$currentHdlFileName === name}
      class:bg-white-800={$currentHdlFileName === name}
    >
      <button
        oncontextmenu={(event) => handleContextMenu(event, name)}
        onclick={() => currentHdlFileName.set(name)}
        class="w-full cursor-pointer py-1 pr-4 pl-6 text-left"
      >
        {name}.hdl
      </button>
    </div>
  {/each}
  {#if $showFileContextMenu}
    <div
      class="bg-white-300 dark:bg-silver-900 flex h-[72px] w-[100px] flex-col"
      style="position: absolute; top:{$fileContextMenuPosition.y}px; left:{$fileContextMenuPosition.x}px"
    >
      <button
        onclick={() => {
          console.log("Rename clicked ", $fileContextMenuClickedFileName);
        }}
        class="dark:hover:bg-silver-800 stroke-text dark:stroke-text-dark hover:bg-white-800 flex w-full cursor-pointer items-center justify-between p-2 text-left text-sm"
      >
        <span>Rename</span>
        <EditIcon classes="h-4 w-4 ml-2" />
      </button>
      <button
        onclick={() => {
          console.log("Delete clicked ", $fileContextMenuClickedFileName);
        }}
        class="dark:hover:bg-silver-800 stroke-text dark:stroke-text-dark hover:bg-white-800 flex w-full cursor-pointer items-center justify-between p-2 text-left text-sm hover:stroke-red-500"
      >
        <span>Delete</span>
        <TrashIcon classes="h-4 w-4 ml-2" />
      </button>
    </div>
  {/if}
</div>
