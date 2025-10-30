<script>
  import FileTextIcon from "../../../../components/icons/FileText.svelte";
</script>

<button
  data-modal-target="built-in-chips-modal"
  data-modal-toggle="built-in-chips-modal"
  title="View built-in chips"
  class={`
  dark:hover:bg-silver-800
   hover:bg-white-900 flex
   cursor-pointer items-center gap-2 rounded-md
   p-2 text-sm whitespace-nowrap
  `}
>
  <FileTextIcon classes="h-4 w-4" />
  <span class="hidden md:block">Built-In Chips</span>
</button>
