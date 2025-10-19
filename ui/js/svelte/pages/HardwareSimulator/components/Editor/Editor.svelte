<script lang="ts">
  // layout and theme css
  import "prism-code-editor/layout.css";
  // hdl language
  import "prism-code-editor/prism/languages/nand2tetris-hdl";
  // required by autoComplete()
  import "prism-code-editor/autocomplete.css";
  import "prism-code-editor/autocomplete-icons.css";
  // required by searchWidget()
  import "prism-code-editor/search.css";
  // required by copyButton()
  import "prism-code-editor/copy-button.css";
  // required by indentGuides()
  import "prism-code-editor/guides.css";

  import { createEditor, type PrismEditor } from "prism-code-editor";

  // extensions
  import { matchBrackets } from "prism-code-editor/match-brackets";
  import {
    searchWidget,
    highlightSelectionMatches,
  } from "prism-code-editor/search";
  import { defaultCommands, editHistory } from "prism-code-editor/commands";
  import { cursorPosition } from "prism-code-editor/cursor";
  import { matchTags } from "prism-code-editor/match-tags";
  import { highlightBracketPairs } from "prism-code-editor/highlight-brackets";
  import { indentGuides } from "prism-code-editor/guides";
  import {
    fuzzyFilter,
    autoComplete,
    type Completion,
    registerCompletions,
    type CompletionSource,
  } from "prism-code-editor/autocomplete";

  import { getClosestToken } from "prism-code-editor/utils";
  import { loadTheme } from "prism-code-editor/themes";

  import { onMount } from "svelte";
  import { editorErrors, hdl } from "../../store";

  let editor: PrismEditor;

  onMount(() => {
    const HDL_KEYWORDS = ["CHIP", "IN", "OUT", "PARTS:"] as const;

    const options: Completion[] = HDL_KEYWORDS.map((label) => ({
      label,
      icon: "keyword",
    }));

    const mySource: CompletionSource = (context, editor) => {
      if (getClosestToken(editor, ".string, .comment", 0, 0, context.pos)) {
        return; // Disable autocomplete in comments and strings
      }
      const wordBefore = /\w*$/.exec(context.lineBefore)![0];

      if (wordBefore || context.explicit) {
        return {
          from: context.pos - wordBefore.length,
          options: options,
        };
      }
    };

    editor = createEditor(
      "#editor",
      {
        language: "nand2tetris-hdl",
        value: $hdl,
        tabSize: 4,
        // class: "h-[calc(100dvh-16px-var(--header-height))]",
        onUpdate: (newValue) => {
          hdl.set(newValue);
        },
      },
      matchBrackets(),
      highlightSelectionMatches(),
      searchWidget(),
      defaultCommands(),
      matchTags(),
      highlightBracketPairs(),
      cursorPosition(),
      editHistory(),
      indentGuides(),
    );

    registerCompletions(["nand2tetris-hdl"], {
      sources: [mySource],
    });

    editor.addExtensions(
      autoComplete({
        filter: fuzzyFilter,
      }),
    );

    const isDark = document.documentElement.classList.contains("dark");
    changeTheme(isDark);

    function changeTheme(isDark: boolean) {
      const editorStyle = document.querySelector("#editor-style");
      if (!editorStyle) {
        return;
      }

      loadTheme(isDark ? "vs-code-dark" : "vs-code-light").then((theme) => {
        editorStyle.textContent = theme as string;
      });
    }

    const observer = new MutationObserver(() => {
      const isDark = document.documentElement.classList.contains("dark");
      changeTheme(isDark);
    });

    observer.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ["class"],
    });

    editorErrors.subscribe((errors) => {
      const lines = editor.lines;
      // clear every error
      for (let i = 1; i <= lines.length; i++) {
        const line = lines[i];
        if (line) {
          line.style.textDecoration = "";
        }
      }

      errors.forEach((err) => {
        const line = lines[err.line];
        if (line) {
          line.style.textDecoration = "red wavy underline";
        }
      });
    });

    hdl.subscribe((value) => {
      editor.setOptions({ value });
    });

    return () => observer.disconnect();
  });
</script>

<div class="">
  <style id="editor-style"></style>
  <div id="editor"></div>
  {#if $editorErrors.length != 0}
    <div
      class="h-[50px] overflow-auto rounded-lg border border-red-500 bg-red-500/10 p-2.5 dark:bg-red-500/20"
    >
      {#each $editorErrors as err (err.line)}
        <div class="flex items-center">
          Error at line {err.line}: {err.message}
        </div>
      {/each}
    </div>
  {/if}
</div>
