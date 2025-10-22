<script lang="ts">
  import "prism-code-editor/layout.css";
  import { createEditor } from "prism-code-editor";
  import { onMount } from "svelte";
  import { editorError, hdl } from "../../store";
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
        value: $hdl,
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

    editorError.subscribe((error) => {
      highlightError(editor, error);
    });

    hdl.subscribe((value) => {
      editor.setOptions({ value });
    });

    return () => themeChangeObserver.disconnect();
  });
</script>

<div class="relative h-full">
  <div
    class={`
    ${$editorError ? "h-[calc(100%-50px)] max-h-[calc(100%-50px)]" : "h-full max-h-full"}
    overflow-auto
    `}
  >
    <style id="editor-style"></style>
    <div id="editor"></div>
  </div>
  {#if $editorError}
    <ErrorBox errorMessage={$editorError.message} />
  {/if}
</div>
