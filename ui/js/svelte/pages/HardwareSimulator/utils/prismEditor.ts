import type { EditorExtension, PrismEditor } from "prism-code-editor";
import "prism-code-editor/prism/languages/nand2tetris-hdl";
import { loadTheme } from "prism-code-editor/themes";
import { getClosestToken } from "prism-code-editor/utils";

// autocomplete
import {
  autoComplete,
  fuzzyFilter,
  registerCompletions,
} from "prism-code-editor/autocomplete";
import type {
  Completion,
  CompletionSource,
} from "prism-code-editor/autocomplete";

// extensions
import { defaultCommands, editHistory } from "prism-code-editor/commands";
import { cursorPosition } from "prism-code-editor/cursor";
import { indentGuides } from "prism-code-editor/guides";
import { highlightBracketPairs } from "prism-code-editor/highlight-brackets";
import { matchBrackets } from "prism-code-editor/match-brackets";
import { matchTags } from "prism-code-editor/match-tags";
import {
  highlightSelectionMatches,
  searchWidget,
} from "prism-code-editor/search";
// styles required for extensions
import "prism-code-editor/search.css";
import "prism-code-editor/autocomplete.css";
import "prism-code-editor/autocomplete-icons.css";
import "prism-code-editor/copy-button.css";
import "prism-code-editor/guides.css";

import type { EditorError } from "../types";

const HDL_KEYWORDS = ["CHIP", "IN", "OUT", "PARTS:"] as const;

const options: Completion[] = HDL_KEYWORDS.map((label) => ({
  label,
  icon: "keyword",
}));

const hdlSource: CompletionSource = (context, editor) => {
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

export const extensions: EditorExtension[] = [
  matchBrackets(),
  highlightSelectionMatches(),
  searchWidget(),
  defaultCommands(),
  matchTags(),
  highlightBracketPairs(),
  cursorPosition(),
  editHistory(),
  indentGuides(),
  autoComplete({
    filter: fuzzyFilter,
  }),
];

export function registerHDLCompletions() {
  registerCompletions(["nand2tetris-hdl"], {
    sources: [hdlSource],
  });
}

export function changeEditorTheme(isDark: boolean, styleTagId: string) {
  const editorStyleTag = document.querySelector(`#${styleTagId}`);
  if (!editorStyleTag) {
    return;
  }

  loadTheme(isDark ? "vs-code-dark" : "vs-code-light").then((theme) => {
    editorStyleTag.textContent = theme as string;
  });
}

export function startThemeChangeObserver(styleTagId: string): MutationObserver {
  const observer = new MutationObserver(() => {
    const isDark = document.documentElement.classList.contains("dark");
    changeEditorTheme(isDark, styleTagId);
  });

  observer.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ["class"],
  });

  return observer;
}

export function highlightError(editor: PrismEditor, error: EditorError | null) {
  clearErrorHighlight(editor);

  if (!error || !error.line) {
    return;
  }

  const lines = editor.lines;
  const line = lines[error.line];
  if (!line) {
    return;
  }
  line.style.textDecoration = "red wavy underline";
}

export function clearErrorHighlight(editor: PrismEditor) {
  const lines = editor.lines;

  for (let i = 1; i <= lines.length; i++) {
    const line = lines[i];
    if (line) {
      line.style.textDecoration = "";
    }
  }
}
