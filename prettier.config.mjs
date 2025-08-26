import pluginSvelte from "prettier-plugin-svelte";

export default {
  plugins: [
    pluginSvelte,
    "prettier-plugin-tailwindcss",
    "prettier-plugin-templ-script",
  ],
  overrides: [
    {
      files: "*.svelte",
      options: {
        parser: "svelte",
      },
    },
    {
      files: ["*.templ"],
      options: {
        parser: "templ",
      },
    },
  ],
};
