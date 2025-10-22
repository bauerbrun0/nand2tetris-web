<script lang="ts">
  import { derived, writable } from "svelte/store";
  import { hdls } from "../../store";
  import { onMount } from "svelte";
  import { calculateFileContextMenuPosition } from "../../utils/projectChips";
  import type { Dimensions, Position } from "../../utils/projectChips";
  import ProjectNameRow from "./ProjectNameRow.svelte";
  import ChipRow from "./ChipRow.svelte";
  import FileContextMenu from "./FileContextMenu.svelte";

  let containerElement: HTMLDivElement;

  const sortedHdlFileNames = derived(hdls, ($hdls) =>
    Object.keys($hdls).sort((a, b) => a.localeCompare(b)),
  );

  const fileContextMenuPosition = writable<Position>({ x: 0, y: 0 });
  const fileContextMenuDimensions: Dimensions = {
    width: 100,
    height: 72,
  };
  const showFileContextMenu = writable<boolean>(false);
  const fileContextMenuClickedFileName = writable<string>("");

  function handleContextMenu(event: MouseEvent, chipFileName: string) {
    event.preventDefault();

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

<div bind:this={containerElement} class="relative h-full w-full overflow-auto">
  <ProjectNameRow />
  {#each $sortedHdlFileNames as name (name)}
    <ChipRow {name} onContextMenu={handleContextMenu} />
  {/each}
  {#if $showFileContextMenu}
    <FileContextMenu
      position={$fileContextMenuPosition}
      clickedFileName={$fileContextMenuClickedFileName}
    />
  {/if}
</div>
