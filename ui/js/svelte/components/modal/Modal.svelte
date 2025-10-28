<script lang="ts">
  import type { Snippet } from "svelte";

  let {
    id,
    classes,
    children,
  }: {
    id: string;
    classes: string;
    children: Snippet;
  } = $props();
</script>

<div
  {id}
  tabindex="-1"
  aria-hidden="true"
  class={`
	fixed top-0 right-0 left-0 z-50 hidden h-[calc(100%)] max-h-full w-full
	items-center justify-center overflow-x-hidden overflow-y-auto md:inset-0
  `}
>
  <div class={`relative mx-2 w-full ${classes}`}>
    {@render children?.()}
  </div>
</div>
