module.exports = {
  'frontend/**/*.{ts,tsx}': ['eslint --fix'],
  'server/**/*.go': ['gofmt -w'],
};