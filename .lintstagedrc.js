const path = require('path');

module.exports = {
  'frontend/**/*.{ts,tsx}': (filenames) => {
    // Run eslint from the frontend directory so it picks up the config
    // We pass absolute paths, which eslint handles fine
    return `cd frontend && yarn lint:fix ${filenames.join(' ')}`;
  },
  'server/**/*.go': ['gofmt -w'],
};
