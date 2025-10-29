<script lang="ts">
  import CloudIcon from "../../../../components/icons/Cloud.svelte";
  import CloudOffIcon from "../../../../components/icons/CloudOff.svelte";
  import CloudUploadIcon from "../../../../components/icons/CloudUpload.svelte";
  let {
    chipSyncStatus,
  }: { chipSyncStatus: "synced" | "unsynced" | "syncing" } = $props();
</script>

<div class="">
  {#if chipSyncStatus === "synced"}
    <div
      class="flex items-center gap-1 text-green-500"
      title="HDL code is saved"
    >
      <CloudIcon classes="h-4 w-4" />
      <span class="hidden text-sm md:block">Synced</span>
    </div>
  {:else if chipSyncStatus === "unsynced"}
    <div
      class="flex items-center gap-1 text-red-500"
      title="HDL code is not saved"
    >
      <CloudOffIcon classes="h-4 w-4" />
      <span class="hidden text-sm md:block">Unsynced</span>
    </div>
  {:else if chipSyncStatus === "syncing"}
    <div
      class="text-primary-500 flex items-center gap-1"
      title="Saving HDL code..."
    >
      <CloudUploadIcon classes="h-4 w-4 animate-pulse" />
      <span class="hidden text-sm md:block">Syncing</span>
    </div>
  {/if}
</div>
