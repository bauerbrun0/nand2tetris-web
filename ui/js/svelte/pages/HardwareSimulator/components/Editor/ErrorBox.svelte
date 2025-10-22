<script lang="ts">
  export let errorMessage: string;
</script>

<div
  class="absolute bottom-0 flex h-[50px] w-full items-center overflow-auto rounded-lg border border-red-500 bg-red-500/10 p-2.5 dark:bg-red-500/20"
>
  <div class="flex items-center">
    {errorMessage}
  </div>
</div>
