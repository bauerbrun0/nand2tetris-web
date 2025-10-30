<script lang="ts">
  import { onMount, mount } from "svelte";
  import { GoldenLayout } from "golden-layout";
  import "golden-layout/dist/css/goldenlayout-base.css";
  import TopBar from "./components/TopBar/TopBar.svelte";
  import ProjectChips from "./components/ProjectChips/ProjectChips.svelte";
  import Editor from "./components/Editor/Editor.svelte";
  import Simulator from "./components/Simulator/Simulator.svelte";
  import RenameChipModal from "./components/RenameChipModal.svelte";
  import BuiltInChipsModal from "./components/BuiltInChipsModal.svelte";
  import { loadHardwareSimulator } from "./utils/hardwareSimulator.ts";
  import {
    registerComponent,
    disableTooltips,
    getLayoutConfig,
  } from "./utils/goldenLayout.ts";
  import { t } from "../../../utils/i18n/i18n.ts";
  import {
    currentHdlFileName,
    currentProjectName,
    hardwareSimulatorError,
    hdl,
    hdls,
    rightClickedChipName,
  } from "./store.ts";
  import { getProjectSlug } from "./utils/routeParam.ts";
  import { createMutation, createQuery } from "@tanstack/svelte-query";
  import {
    fetchProjectBySlug,
    fetchProjectChips,
    createChip,
    updateChipHdl,
    updateChipName,
    deleteChipRequest,
  } from "./requests.ts";
  import type { Chip } from "../../../types/chips.ts";
  import { showToast } from "../../../utils/toast.ts";
  import DeleteChipModal from "./components/DeleteChipModal.svelte";

  let layoutContainer: HTMLElement;

  let newlyCreatingChipName: string = $state("");
  let previousHdlFileName: string = "";
  let chipSyncStatus = $state<"synced" | "unsynced" | "syncing">("synced");
  let rightClickedChip = $state<Chip | null>(null);

  const projectSlug = getProjectSlug();

  const projectQuery = createQuery(() => ({
    queryKey: ["project"],
    queryFn: () => fetchProjectBySlug(projectSlug),
  }));

  const chipsQuery = createQuery(() => ({
    queryKey: ["projectChips"],
    queryFn: () => fetchProjectChips(projectQuery.data?.id as number),
    enabled: projectQuery.data !== undefined,
  }));

  const createChipMutation = createMutation(() => ({
    mutationFn: (data: { name: string }) =>
      createChip(data.name, projectQuery.data?.id as number),
    onSuccess: async () => {
      await chipsQuery.refetch();
      if (newlyCreatingChipName !== "") {
        currentHdlFileName.set(newlyCreatingChipName);
        newlyCreatingChipName = "";
      }
    },
  }));

  const { mutate: mutateChipHdl } = createMutation(() => ({
    mutationFn: (data: { id: number; hdl: string }) => {
      chipSyncStatus = "syncing";
      return updateChipHdl(projectQuery.data?.id as number, data.id, data.hdl);
    },
    onError: (error: unknown) => {
      chipSyncStatus = "unsynced";
      showToast({
        duration: 3000,
        message: (error as Error).message,
        variant: "error",
      });
    },
    onSuccess: () => {
      chipSyncStatus = "synced";
    },
  }));

  const { mutateAsync: mutateChipName } = createMutation(() => ({
    mutationFn: (data: { id: number; name: string }) => {
      return updateChipName(
        projectQuery.data?.id as number,
        data.id,
        data.name,
      );
    },
    onError: (error: unknown) => {
      showToast({
        duration: 3000,
        message: (error as Error).message,
        variant: "error",
      });
    },
    onSuccess: async (chip) => {
      await chipsQuery.refetch();
      currentHdlFileName.set(chip.name);
    },
  }));

  const deleteChipMutation = createMutation(() => ({
    mutationFn: (data: { id: number }) => {
      return deleteChipRequest(projectQuery.data?.id as number, data.id);
    },
    onError: (error: unknown) => {
      showToast({
        duration: 3000,
        message: (error as Error).message,
        variant: "error",
      });
    },
    onSuccess: () => {
      chipsQuery.refetch();
    },
  }));

  $effect(() => {
    if (projectQuery.data) {
      currentProjectName.set(projectQuery.data.title);
    }
  });

  $effect(() => {
    if ($currentHdlFileName === null || $hdl === null) {
      return;
    }

    if (previousHdlFileName !== $currentHdlFileName) {
      previousHdlFileName = $currentHdlFileName;
      return;
    }

    const chipId = getChipId($currentHdlFileName);
    if (chipId === null) {
      return;
    }
    mutateChipHdl({ id: chipId, hdl: $hdl });
  });

  $effect(() => {
    if (chipsQuery.data) {
      let newHdls: Record<string, string> = {};
      for (const chip of chipsQuery.data) {
        newHdls[chip.name] = chip.hdl;
      }
      hdls.set(newHdls);

      if (chipsQuery.data.length > 0) {
        currentHdlFileName.set((chipsQuery.data[0] as Chip).name);
      } else {
        currentHdlFileName.set(null);
        hardwareSimulatorError.set(null);
      }
    }
  });

  async function createNewChip(name: string) {
    try {
      newlyCreatingChipName = name;
      await createChipMutation.mutateAsync({ name });
    } catch (error: unknown) {
      showToast({
        duration: 3000,
        message: (error as Error).message,
        variant: "error",
      });
    }
  }

  function getChipId(name: string): number | null {
    const chip = chipsQuery.data?.find((chip) => chip.name === name);
    return chip ? chip.id : null;
  }

  rightClickedChipName.subscribe((name) => {
    if (chipsQuery.data) {
      const chip = chipsQuery.data.find((chip) => chip.name === name);
      rightClickedChip = chip ? chip : null;
    }
  });

  async function renameChip(id: number, name: string): Promise<boolean> {
    try {
      await mutateChipName({ id, name });
      return true;
    } catch {
      return false;
    }
  }

  async function deleteChip(id: number): Promise<boolean> {
    try {
      await deleteChipMutation.mutateAsync({ id });
      return true;
    } catch {
      return false;
    }
  }

  onMount(() => {
    loadHardwareSimulator().then(() => {
      hdl.subscribe((hdl) => {
        if (hdl === null) {
          return;
        }
        hardwareSimulatorError.set(null);
        window.WASM.HardwareSimulator.processHdls();
      });
    });

    const layout = new GoldenLayout(layoutContainer);
    registerComponent(layout, "editor", Editor);
    registerComponent(layout, "simulator", Simulator);
    layout.registerComponentFactoryFunction(
      "project-chips",
      (layoutContainer) => {
        mount(ProjectChips, {
          target: layoutContainer.element,
          props: {
            createChip: createNewChip,
          },
        });
      },
    );

    const layoutConfig = getLayoutConfig({
      editorTitle: t("hardware_simulator_page.editor_window_title"),
      simulatorTitle: t("hardware_simulator_page.simulator_window_title"),
      projectChipsTitle: t(
        "hardware_simulator_page.project_chips_window_title",
      ),
      windowWidth: window.innerWidth,
    });
    layout.loadLayout(layoutConfig);

    disableTooltips();

    return () => {
      layout.destroy();
    };
  });
</script>

<div>
  <TopBar {chipSyncStatus} />
  <div
    bind:this={layoutContainer}
    class="mb-[8px] flex h-[calc(100dvh-16px-var(--header-height)-40px)] flex-auto overflow-hidden"
  ></div>
  <RenameChipModal
    id="rename-chip-modal"
    {renameChip}
    chip={rightClickedChip}
  />
  <DeleteChipModal
    id="delete-chip-modal"
    {deleteChip}
    chip={rightClickedChip}
  />
  <BuiltInChipsModal id="built-in-chips-modal" />
</div>
