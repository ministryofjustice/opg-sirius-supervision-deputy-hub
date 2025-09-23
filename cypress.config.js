const { defineConfig } = require("cypress");

module.exports = defineConfig({
    e2e: {
        // We've imported your old cypress plugins here.
        // You may want to clean this up later by importing these.
        setupNodeEvents(on, config) {
            require("./cypress/plugins/index.js")(on, config);
            const { plugin: cypressGrepPlugin } = require('@cypress/grep/plugin')
            cypressGrepPlugin(config)
            return config;
        },
        baseUrl: "http://localhost:8888",
    },
    env: {
        grepOmitFiltered: true,
        grepFilterSpecs: true,
    },
});
