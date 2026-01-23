module.exports = {
  'frontend/**/*.{js,jsx,ts,tsx}': () => 'cd frontend && yarn lint:fix',
  'server/**/*.go': ['gofmt -w'],
};
