<script lang="ts">
  import Projects from "./Projects.svelte";
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";

  const queryClient = new QueryClient();
</script>

<QueryClientProvider client={queryClient}>
  <Projects />
</QueryClientProvider>
