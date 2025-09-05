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

  import { createEditor } from "prism-code-editor";

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

  const initialValue = `// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/1/And.hdl
/**
* And gate:
* if (a and b) out = 1, else out = 0
*/
CHIP And {
    IN a, b;
    OUT out;

    PARTS:
    Nand(a = a, b = b, out = aNandB);
    Not(in = aNandB, out = out);
}
`;

  onMount(() => {
    const HDL_KEYWORDS = [
      "CHIP",
      "IN",
      "OUT",
      "PARTS:",
      "BUILTIN",
      "CLOCKED",
    ] as const;

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

    const editor = createEditor(
      "#editor",
      {
        language: "nand2tetris-hdl",
        value: initialValue,
        tabSize: 4,
        class: "h-[calc(100dvh-var(--header-height))]",
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

    return () => observer.disconnect();
  });
</script>

<div class="w-1/2">
  <style id="editor-style"></style>
  <div id="editor"></div>
</div>
