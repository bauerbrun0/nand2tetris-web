<script lang="ts">
  let { message }: { message: string } = $props();
</script>

<div
  class="flex flex-col gap-5 rounded-lg border border-red-500 bg-red-500/10 p-5 dark:bg-red-500/20"
>
  <h2 class="text-lg font-bold">An error happened while fetching projects</h2>
  <p class=""><strong>Message:</strong> {message}</p>
</div>
