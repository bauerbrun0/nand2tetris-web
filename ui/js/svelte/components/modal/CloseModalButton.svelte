<script lang="ts">
  import XIcon from "../icons/X.svelte";
  let {
    modalId,
    buttonElement = $bindable(),
  }: { modalId: string; buttonElement?: HTMLButtonElement } = $props();
</script>

<button
  bind:this={buttonElement}
  type="button"
  class={`
    end-2.5 ms-auto inline-flex h-8 w-8 cursor-pointer
	items-center justify-center rounded-lg bg-transparent text-sm text-gray-400 hover:bg-gray-200 hover:text-gray-900
	dark:hover:bg-gray-600 dark:hover:text-white
  `}
  data-modal-hide={modalId}
>
  <XIcon classes="w-5 h-5" />
  <span class="sr-only">Close modal</span>
</button>
