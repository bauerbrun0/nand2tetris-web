<script lang="ts">
  import { simulationSpeeds } from "../../../utils/simulation";
  import { simulationSpeed } from "../../../store";
  import { Dropdown } from "flowbite";
  import { onMount } from "svelte";

  let targetElement: HTMLElement;
  let triggerElement: HTMLElement;
  let dropdown: Dropdown;
  onMount(() => {
    dropdown = new Dropdown(targetElement, triggerElement);
  });
</script>

<button
  bind:this={triggerElement}
  class={`
    dark:bg-silver-800 dark:hover:bg-silver-700 hover:bg-white-400 inline-flex h-[32px] cursor-pointer
    items-center rounded-lg bg-white px-4 py-1 text-center text-sm font-medium whitespace-nowrap
  `}
  type="button"
  >{$simulationSpeed.text}<svg
    class="ms-3 h-2.5 w-2.5"
    aria-hidden="true"
    xmlns="http://www.w3.org/2000/svg"
    fill="none"
    viewBox="0 0 10 6"
  >
    <path
      stroke="currentColor"
      stroke-linecap="round"
      stroke-linejoin="round"
      stroke-width="1.5"
      d="m1 1 4 4 4-4"
    />
  </svg>
</button>

<div
  bind:this={targetElement}
  class="dark:bg-silver-900 bg-white-500 z-10 hidden w-44 divide-y divide-gray-100 rounded-lg"
>
  <ul class="py-2 text-sm" aria-labelledby="speed-dropdown-button">
    {#each simulationSpeeds as speed (speed.delayMs)}
      <li>
        <button
          on:click={() => {
            simulationSpeed.set(speed);
            dropdown.hide();
          }}
          class="dark:hover:bg-silver-800 hover:bg-white-800 block w-full cursor-pointer px-4 py-2 text-left whitespace-nowrap"
          >{speed.text}</button
        >
      </li>
    {/each}
  </ul>
</div>
