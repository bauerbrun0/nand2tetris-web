<script>
  import ClockButton from "./ClockButton.svelte";
  import EvaluateButton from "./EvaluateButton.svelte";
  import ResetButton from "./ResetButton.svelte";
  import RunSection from "./RunSection.svelte";
</script>

<div class="flex items-center gap-2">
  <EvaluateButton />
  <ClockButton />
  <RunSection />
  <ResetButton />
</div>
