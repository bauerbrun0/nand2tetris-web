<script lang="ts">
  import { derived, writable } from "svelte/store";
  import { hdls, rightClickedChipName } from "../../store";
  import { onMount } from "svelte";
  import { calculateFileContextMenuPosition } from "../../utils/projectChips";
  import type { Dimensions, Position } from "../../utils/projectChips";
  import ProjectNameRow from "./ProjectNameRow.svelte";
  import ChipRow from "./ChipRow.svelte";
  import FileContextMenu from "./FileContextMenu.svelte";
  import NewChipInput from "./NewChipInput.svelte";

  let {
    createChip,
  }: {
    createChip: (name: string) => void;
  } = $props();

  let containerElement: HTMLDivElement;

  const sortedHdlFileNames = derived(hdls, ($hdls) =>
    Object.keys($hdls).sort((a, b) => a.localeCompare(b)),
  );

  const fileContextMenuPosition = writable<Position>({ x: 0, y: 0 });
  const fileContextMenuDimensions: Dimensions = {
    width: 110,
    height: 72,
  };
  const showFileContextMenu = writable<boolean>(false);
  const fileContextMenuClickedFileName = writable<string>("");

  function handleContextMenu(event: MouseEvent, chipFileName: string) {
    event.preventDefault();

    rightClickedChipName.set(chipFileName);

    const position = calculateFileContextMenuPosition(
      containerElement,
      event,
      fileContextMenuDimensions,
    );

    fileContextMenuPosition.set(position);
    fileContextMenuClickedFileName.set(chipFileName);
    showFileContextMenu.set(true);
  }

  onMount(() => {
    // just so there is no eslint warning about setting an onclick on a div
    containerElement.onclick = () => showFileContextMenu.set(false);
  });
</script>

<div
  bind:this={containerElement}
  class="relative h-full w-full overflow-auto font-mono"
>
  <ProjectNameRow />
  <NewChipInput {createChip} />
  {#each $sortedHdlFileNames as name (name)}
    <ChipRow {name} onContextMenu={handleContextMenu} />
  {/each}
  {#if $showFileContextMenu}
    <FileContextMenu position={$fileContextMenuPosition} />
  {/if}
</div>
