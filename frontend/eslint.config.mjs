import nextConfig from "eslint-config-next/core-web-vitals";

/** @type {import('eslint').Linter.Config[]} */
const eslintConfig = [
  ...nextConfig,
  {
    // Ignore e2e test files - Playwright's `use` function triggers react-hooks/rules-of-hooks
    ignores: ["e2e/**"],
  },
];

export default eslintConfig;
