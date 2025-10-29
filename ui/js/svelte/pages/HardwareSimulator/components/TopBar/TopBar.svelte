<script lang="ts">
  import BackToProjectsLink from "./BackToProjectsLink.svelte";
  import ChipName from "./ChipName.svelte";
  import ChipSyncStatus from "./ChipSyncStatus.svelte";

  let {
    chipSyncStatus,
  }: { chipSyncStatus: "synced" | "unsynced" | "syncing" } = $props();
</script>

<div class="flex h-[48px] w-full items-center justify-between py-[8px]">
  <div class="flex items-center">
    <BackToProjectsLink />
    <ChipName />
  </div>
  <div>
    <ChipSyncStatus {chipSyncStatus} />
  </div>
</div>
