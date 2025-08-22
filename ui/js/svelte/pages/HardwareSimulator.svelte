<script lang="ts">
    let quotes: string[] = [
        "Simplicity is the soul of efficiency.",
        "Programs must be written for people to read.",
        "Premature optimization is the root of all evil.",
        "Talk is cheap. Show me the code.",
        "Code is like humor. When you have to explain it, it’s bad.",
    ];

    let current = quotes[Math.floor(Math.random() * quotes.length)];

    function shuffle() {
        current = quotes[Math.floor(Math.random() * quotes.length)];
    }
</script>

<div class="p-4 rounded-xl shadow-md">
    <p class="italic mb-2">“{current}”</p>
    <button
        class="px-3 py-1 bg-red-500 text-white rounded-lg"
        on:click={shuffle}
    >
        New Quote
    </button>
</div>
