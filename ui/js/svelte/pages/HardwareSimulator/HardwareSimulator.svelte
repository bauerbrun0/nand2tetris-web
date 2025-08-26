<script lang="ts">
  let quotes: string[] = [
    "Simplicity is the soul of efficiency.",
    "Programs must be written for people to read.",
    "Premature optimization is the root of all evil.",
    "Talk is cheap. Show me the code.",
    "Code is like humor. When you have to explain it, it’s bad.",
  ];

  let current = quotes[Math.floor(Math.random() * quotes.length)];

  function shuffle(_a: number) {
    current = quotes[Math.floor(Math.random() * quotes.length)];
  }
</script>

<div class="rounded-xl p-4 shadow-md">
  <p class="mb-2 italic">“{current}”</p>
  <button
    class="rounded-lg bg-red-500 px-3 py-1 text-white"
    on:click={() => shuffle(2)}
  >
    New Quote
  </button>
</div>
