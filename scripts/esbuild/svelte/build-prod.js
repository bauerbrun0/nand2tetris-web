import fs from "fs";
import esbuild from "esbuild";
import esbuildSvelte from "esbuild-svelte";
import { sveltePreprocess } from "svelte-preprocess";

if (!fs.existsSync("./ui/static/js/")) {
  fs.mkdirSync("./ui/static/js/");
}
esbuild
  .build({
    entryPoints: {
      "main-hardware-simulator": "./ui/js/entries/main-hardware-simulator.ts",
    },
    bundle: true,
    outdir: `./ui/static/js/`,
    mainFields: ["svelte", "browser", "module", "main"],
    conditions: ["svelte", "browser"],
    minify: true,
    sourcemap: false,
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
  })
  .catch((error, location) => {
    console.warn(`Errors: `, error, location);
    process.exit(1);
  });
