<script lang="ts">
  import { currentProjectName, currentHdlFileName } from "../../store";
</script>

<h1 class="ml-4 font-bold whitespace-nowrap">
  {$currentProjectName}
  {$currentHdlFileName ? ` / ${$currentHdlFileName}.hdl` : ""}
</h1>
