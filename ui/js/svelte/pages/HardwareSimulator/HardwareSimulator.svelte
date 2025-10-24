<script lang="ts">
  import { onMount } from "svelte";
  import { GoldenLayout } from "golden-layout";
  import "golden-layout/dist/css/goldenlayout-base.css";
  import TopBar from "./components/TopBar/TopBar.svelte";
  import ProjectChips from "./components/ProjectChips/ProjectChips.svelte";
  import Editor from "./components/Editor/Editor.svelte";
  import Simulator from "./components/Simulator/Simulator.svelte";
  import { loadHardwareSimulator } from "./utils/hardwareSimulator.ts";
  import {
    registerComponent,
    disableTooltips,
    getLayoutConfig,
  } from "./utils/goldenLayout.ts";
  import { t } from "../../../utils/i18n/i18n.ts";
  import { hardwareSimulatorError, hdl } from "./store.ts";

  let layoutContainer: HTMLElement;

  onMount(() => {
    loadHardwareSimulator().then(() => {
      hdl.subscribe((_hdl) => {
        hardwareSimulatorError.set(null);
        window.WASM.HardwareSimulator.processHdls();
      });
    });

    const layout = new GoldenLayout(layoutContainer);
    registerComponent(layout, "editor", Editor);
    registerComponent(layout, "simulator", Simulator);
    registerComponent(layout, "project-chips", ProjectChips);

    const layoutConfig = getLayoutConfig({
      editorTitle: t("hardware_simulator_page.editor_window_title"),
      simulatorTitle: t("hardware_simulator_page.simulator_window_title"),
      projectChipsTitle: t(
        "hardware_simulator_page.project_chips_window_title",
      ),
    });
    layout.loadLayout(layoutConfig);

    disableTooltips();

    return () => {
      layout.destroy();
    };
  });
</script>

<div>
  <TopBar />
  <div
    bind:this={layoutContainer}
    class="mb-[8px] flex h-[calc(100dvh-16px-var(--header-height)-40px)] flex-auto overflow-hidden"
  ></div>
</div>
