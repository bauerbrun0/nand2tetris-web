<script lang="ts">
  import "prism-code-editor/layout.css";
  import { createEditor } from "prism-code-editor";
  import { onMount } from "svelte";
  import { hardwareSimulatorError, hdl } from "../../store";
  import ErrorBox from "./ErrorBox.svelte";
  import {
    changeEditorTheme,
    extensions,
    highlightError,
    registerHDLCompletions,
    startThemeChangeObserver,
  } from "../../utils/prismEditor";

  onMount(() => {
    const editor = createEditor(
      "#editor",
      {
        language: "nand2tetris-hdl",
        value: $hdl || "",
        tabSize: 4,
        onUpdate: (newValue) => {
          hdl.set(newValue);
        },
      },
      ...extensions,
    );
    registerHDLCompletions();

    const isDark = document.documentElement.classList.contains("dark");
    changeEditorTheme(isDark, "editor-style");
    const themeChangeObserver = startThemeChangeObserver("editor-style");

    hardwareSimulatorError.subscribe((error) => {
      highlightError(editor, error);
    });

    hdl.subscribe((value) => {
      if (value === null) {
        return;
      }
      editor.setOptions({ value });
    });

    return () => themeChangeObserver.disconnect();
  });
</script>

<div class="relative h-full">
  <div
    class={`
        ${$hardwareSimulatorError ? "h-[calc(100%-50px)] max-h-[calc(100%-50px)]" : "h-full max-h-full"}
        overflow-auto
      `}
  >
    <style id="editor-style"></style>
    <div class:hidden={$hdl === null} id="editor"></div>
  </div>
  {#if $hardwareSimulatorError}
    <ErrorBox errorMessage={$hardwareSimulatorError.message} />
  {/if}
</div>
