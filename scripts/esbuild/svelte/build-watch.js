import fs from "fs";
import * as esbuild from "esbuild";
import esbuildSvelte from "esbuild-svelte";
import { sveltePreprocess } from "svelte-preprocess";

if (!fs.existsSync("./ui/static/js/")) {
  fs.mkdirSync("./ui/static/js/");
}

let ctx = await esbuild.context({
  entryPoints: {
    "main-hardware-simulator": "./ui/js/entries/main-hardware-simulator.ts",
    "main-projects": "./ui/js/entries/main-projects.ts",
  },
  bundle: true,
  outdir: `./ui/static/js/`,
  mainFields: ["svelte", "browser", "module", "main"],
  conditions: ["svelte", "browser"],
  minify: false,
  sourcemap: true,
  splitting: true,
  write: true,
  format: `esm`,
  loader: {
    ".png": "file",
    ".jpg": "file",
    ".svg": "file",
  },
  plugins: [
    esbuildSvelte({
      preprocess: sveltePreprocess(),
    }),
  ],
});

try {
  await ctx.watch();
  console.log("watching for changes");
} catch (error) {
  console.warn(`Errors: `, error);
  process.exit(1);
}
