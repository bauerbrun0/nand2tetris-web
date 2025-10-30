<script lang="ts">
  import type { Snippet } from "svelte";
  import type { Writable } from "svelte/store";

  let {
    label = "",
    name = "",
    placeholder = "",
    type = "text",
    icon = null,
    error = "",
    required = false,
    value,
  }: {
    label?: string;
    name?: string;
    placeholder?: string;
    type?: string;
    icon?: Snippet | null;
    error?: string;
    required?: boolean;
    value: Writable<string>;
  } = $props();

  const classes = `
    block w-full rounded-lg
	bg-white-500 focus:bg-white-200 dark:bg-silver-800 dark:focus:bg-silver-900
	text-silver-800 dark:text-silver-100
	border border-silver-800 dark:border-silver-700
	focus:border-primary-500 dark:focus:border-primary-500
	focus:ring-primary-500 dark:focus:ring-primary-500
	${error !== "" ? "!border-red-500" : ""}
	${icon !== null ? "ps-10 p-2.5" : "px-3 py-2"}
`;
</script>

{#if error === "" && label !== ""}
  <label
    for={name}
    class="text-text-light dark:text-text-dark mb-2 block font-medium"
  >
    {label}
  </label>
{/if}
{#if error !== "" && label !== ""}
  <label for={name} class="mb-2 block font-medium text-red-500">
    {label} - {error}
  </label>
{/if}
{#if error !== "" && label === ""}
  <label for={name} class="mb-2 block font-medium text-red-500">
    {error}
  </label>
{/if}
<div class="relative mb-6">
  {#if icon !== null}
    <div
      class="pointer-events-none absolute inset-y-0 start-0 flex items-center ps-3.5"
    >
      {@render icon()}
    </div>
  {/if}
  {#if type === "textarea"}
    <textarea
      {required}
      bind:value={$value}
      onchange={(e) => {
        value.set((e.target as HTMLTextAreaElement).value);
      }}
      id={name}
      {name}
      class={classes}
      {placeholder}
    ></textarea>
  {:else}
    <input
      {required}
      bind:value={$value}
      onchange={(e) => {
        value.set((e.target as HTMLInputElement).value);
      }}
      {type}
      id={name}
      {name}
      class={classes}
      {placeholder}
    />
  {/if}
</div>
