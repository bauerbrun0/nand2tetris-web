<script lang="ts">
  import { onMount, mount } from "svelte";
  import { progressWASM } from "./store.ts";
  import Editor from "./components/Editor/Editor.svelte";
  import TaskDescription from "./components/TaskDescription/TaskDescription.svelte";
  import { GoldenLayout, ItemType } from "golden-layout";
  import "golden-layout/dist/css/goldenlayout-base.css";

  let layoutContainer: HTMLElement;

  onMount(() => {
    window.WASM = {} as typeof window.WASM;
    window.WASM.HardwareSimulator = {} as typeof window.WASM.HardwareSimulator;
    window.WASM.HardwareSimulator.setProgressWASM = (str) => {
      progressWASM.set(str);
    };

    const go = new Go();
    WebAssembly.instantiateStreaming(
      fetch("/static/wasm/hardware_simulator.wasm"),
      go.importObject,
    ).then((result) => {
      go.run(result.instance);
    });

    const layout = new GoldenLayout(layoutContainer);

    layout.registerComponentFactoryFunction("editor", (container) => {
      mount(Editor, { target: container.element });
    });
    layout.registerComponentFactoryFunction("task", (container) => {
      mount(TaskDescription, { target: container.element });
    });

    layout.loadLayout({
      settings: {
        showPopoutIcon: false,
      },
      dimensions: {
        headerHeight: 36,
      },
      root: {
        type: ItemType.row,
        content: [
          {
            type: ItemType.component,
            componentType: "editor",
            title: "💻 Editor",
            width: 70,
          },
          {
            type: ItemType.component,
            componentType: "task",
            title: "📄 Task",
            // width: 50,
          },
          {
            type: ItemType.component,
            componentType: "task",
            title: "📄 Task",
            // width: 50,
          },
        ],
      },
    });

    // change the editor's height if the pane's height changes
    const editorContent = document.querySelector(
      "#svelte-app > div > div > div > div:nth-child(1) > section.lm_items > div > div",
    );
    const editor = document.getElementsByClassName("prism-code-editor");
    const resizeObserver = new ResizeObserver((entries) => {
      for (const entry of entries) {
        const newHeight = entry.contentRect.height;
        if (editor[0]) {
          (editor[0] as HTMLElement).style.height = `${newHeight}px`;
        }
      }
    });
    resizeObserver.observe(editorContent as HTMLElement);

    // remove all title attributes to disable tooltips
    document.querySelectorAll(".lm_tab").forEach((tab) => {
      tab.removeAttribute("title");
    });

    // start an observer to dynamically disable tooltips on new
    // lm_tabs
    const mutationObserver = new MutationObserver((mutations) => {
      for (const mutation of mutations) {
        mutation.addedNodes.forEach((node) => {
          if (node.nodeType === Node.ELEMENT_NODE) {
            const el = node as HTMLElement;

            // check if added node is an .lm_tab
            if (el.matches(".lm_tab[title]")) {
              el.removeAttribute("title");
            }
          }
        });
      }
    });
    mutationObserver.observe(
      document.querySelector(".lm_goldenlayout") as HTMLElement,
      {
        childList: true,
        subtree: true,
      },
    );
  });
</script>

<div
  bind:this={layoutContainer}
  class="my-[8px] flex h-[calc(100dvh-16px-var(--header-height))] flex-auto overflow-hidden"
></div>

<style>
</style>
