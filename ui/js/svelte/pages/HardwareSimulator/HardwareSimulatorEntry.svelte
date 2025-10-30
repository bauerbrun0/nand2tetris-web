<script lang="ts">
  import HardwareSimulator from "./HardwareSimulator.svelte";
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";

  const queryClient = new QueryClient();
</script>

<QueryClientProvider client={queryClient}>
  <HardwareSimulator />
</QueryClientProvider>
