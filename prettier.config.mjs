import pluginSvelte from "prettier-plugin-svelte";

export default {
  plugins: [pluginSvelte, "prettier-plugin-tailwindcss"],
  overrides: [
    {
      files: "*.svelte",
      options: {
        parser: "svelte",
      },
    },
  ],
};
