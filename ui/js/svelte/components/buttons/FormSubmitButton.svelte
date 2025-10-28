<script lang="ts">
  let {
    text = "Submit",
    loading = false,
    classes = "",
  }: { text?: string; loading?: boolean; classes?: string } = $props();
</script>

<button
  type="submit"
  disabled={loading}
  class={`
	bg-primary-700
	hover:bg-primary-600 focus:ring-primary-300 dark:hover:bg-primary-800 dark:focus:ring-primary-800
	disabled:bg-primary-400 w-full rounded-lg px-5
	py-2.5 font-medium text-white hover:cursor-pointer
	focus:ring-4
	focus:outline-none disabled:cursor-wait
	${classes}
  `}
>
  {text}
</button>
