<script lang="ts">
  import { currentHdlFileName } from "../../store";

  export let name: string;
  export let onContextMenu: (event: MouseEvent, chipFileName: string) => void;
</script>

<div
  class="dark:hover:bg-silver-800 hover:bg-white-800 flex w-full items-center justify-between"
  class:dark:bg-silver-800={$currentHdlFileName === name}
  class:bg-white-800={$currentHdlFileName === name}
>
  <button
    oncontextmenu={(event) => onContextMenu(event, name)}
    onclick={() => currentHdlFileName.set(name)}
    class="w-full cursor-pointer py-1 pr-4 pl-6 text-left"
  >
    {name}.hdl
  </button>
</div>
