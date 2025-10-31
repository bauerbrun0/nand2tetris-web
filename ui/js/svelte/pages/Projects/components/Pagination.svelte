<script lang="ts">
  import { currentPage, pageCount } from "../store";

  function handlePreviousClick() {
    currentPage.set($currentPage - 1);
  }

  function handleNextClick() {
    currentPage.set($currentPage + 1);
  }
</script>

<div class="flex items-center justify-center">
  <button
    onclick={handlePreviousClick}
    disabled={$currentPage === 1}
    class={`
        dark:bg-silver-900 dark:hover:bg-silver-800 order-white-300
        hover:bg-primary-400 bg-primary-500 me-3 flex h-10 cursor-pointer items-center justify-center
        rounded-lg px-4  text-base
        font-medium disabled:cursor-not-allowed
    `}
  >
    <svg
      class="me-2 h-3.5 w-3.5 rtl:rotate-180"
      aria-hidden="true"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 14 10"
    >
      <path
        stroke="currentColor"
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="2"
        d="M13 5H1m0 0 4 4M1 5l4-4"
      />
    </svg>
    Previous
  </button>
  <button
    onclick={handleNextClick}
    disabled={$currentPage === $pageCount}
    class={`
      dark:bg-silver-900 dark:hover:bg-silver-800 order-white-300
      hover:bg-primary-400 bg-primary-500 me-3 flex h-10 cursor-pointer items-center justify-center
      rounded-lg px-4  text-base
      font-medium disabled:cursor-not-allowed
    `}
  >
    Next
    <svg
      class="ms-2 h-3.5 w-3.5 rtl:rotate-180"
      aria-hidden="true"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 14 10"
    >
      <path
        stroke="currentColor"
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="2"
        d="M1 5h12m0 0L9 1m4 4L9 9"
      />
    </svg>
  </button>
</div>
