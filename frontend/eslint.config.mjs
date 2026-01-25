import nextConfig from "eslint-config-next/core-web-vitals";

/** @type {import('eslint').Linter.Config[]} */
const eslintConfig = [
  ...nextConfig,
  {
    // Ignore e2e test files - Playwright's `use` function triggers react-hooks/rules-of-hooks
    // Ignore Playwright report/results - generated minified JS triggers lint errors
    ignores: ["e2e/**", "playwright-report/**", "test-results/**"],
  },
];

export default eslintConfig;
