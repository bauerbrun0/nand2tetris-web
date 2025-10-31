export default {
  extends: ["@commitlint/config-conventional"],
  rules: {
    // Disable max line length for body
    "body-max-line-length": [0, "always"],
  },
};
