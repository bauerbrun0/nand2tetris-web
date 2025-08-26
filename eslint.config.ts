import js from "@eslint/js";
import globals from "globals";
import tseslint from "typescript-eslint";
import svelte from "eslint-plugin-svelte";
import svelteParser from "svelte-eslint-parser";
import eslintConfigPrettier from "eslint-config-prettier/flat";
import type { Linter } from "eslint";
import { defineConfig, globalIgnores } from "eslint/config";

const commonRules: Partial<Linter.RulesRecord> = {
  "@typescript-eslint/no-unused-vars": [
    "error",
    { argsIgnorePattern: "^_", varsIgnorePattern: "^_" },
  ],
};

const baseConfig: Linter.Config[] = [
  js.configs.recommended,
  ...tseslint.configs.recommended,
  {
    files: ["**/*.{js,mjs,cjs,ts,mts,cts}"],
    plugins: { js },
    languageOptions: { globals: globals.browser },
    rules: {
      ...commonRules,
    },
  },
];

const svelteConfig: Linter.Config[] = [
  ...svelte.configs.recommended,
  {
    files: ["**/*.svelte"],
    languageOptions: {
      parser: svelteParser,
      parserOptions: {
        parser: tseslint.parser, // enable TS in <script lang="ts">
        extraFileExtensions: [".svelte"],
        svelteFeatures: { runes: true }, // enable Svelte 5 runes
      },
      globals: {
        ...globals.browser,
      },
    },
    rules: {
      ...commonRules,
    },
  },
];

export default defineConfig([
  globalIgnores(["node_modules/", "ui/static/"]),
  ...baseConfig,
  ...svelteConfig,
  // node environment for esbuild scripts
  {
    files: ["scripts/esbuild/svelte/*.{ts,js}"],
    languageOptions: { globals: globals.node },
  },
  // prettier eslint config which "Turns off all rules that are unnecessary or might conflict with Prettier."
  eslintConfigPrettier,
]);
