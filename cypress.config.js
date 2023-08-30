const {defineConfig} = require('cypress');

module.exports = defineConfig({
    e2e: {
        // We've imported your old cypress plugins here.
        // You may want to clean this up later by importing these.
        setupNodeEvents(on, config) {
            require('./cypress/plugins/index.js')(on, config);
            require('@cypress/grep/src/plugin')(config);
            return config;
        },
        baseUrl: 'http://localhost:8888',
    },
    env: {
        grepOmitFiltered: true,
        grepFilterSpecs: true,
    }
});
