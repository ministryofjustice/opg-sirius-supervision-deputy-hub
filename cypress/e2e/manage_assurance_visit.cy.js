describe("Manage an Assurance Visit", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.setCookie("user", "finance-user");
    });

    describe("Manage an assurance visit form", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/3/manage-assurance-visit/35");
        });

        it("cancel button returns user to the assurance visit page", () => {
            cy.get(".govuk-button-group > .govuk-link")
                .should("contain", "Cancel")
                .click();
            cy.url().should(
                "not.contain",
                "/supervision/deputies/3/manage-assurance-visit/35"
            );
            cy.get(".govuk-main-wrapper > header").contains("Assurance visits");
            cy.get(".govuk-button").contains("Add a visit");
        });

        it("shows relevant heading, field labels and submit button", () => {
            cy.get(".govuk-heading-l").contains("Manage assurance visit");
            cy.get(".govuk-label").contains("Commissioned date");
            cy.get(".govuk-label").contains("Select a visitor");
            cy.get(".govuk-label").contains("Report due date");
            cy.get(".govuk-label").contains("Report received date");
            cy.get(".govuk-fieldset__legend").contains("Outcome");
            cy.get(".govuk-label").contains("Report review date");
            cy.get(".govuk-fieldset__legend").contains("Report marked as");
            cy.get(".govuk-button").contains("Save assurance visit");
        });

        it("form autofills in existing data", () => {
            cy.get("#visit-report-marked-as-Red").should("be.checked");
        });

        it("allows user to edit and submit the form", () => {
            cy.setCookie("success-route", "manageAssuranceVisit");
            cy.get("#f-commissioned-date").type("2021-02-01");
            cy.get("#manage-assurance-visit-form").submit();
            cy.url().should("contain", "/supervision/deputies/3/assurance-visits");
            cy.get(".moj-banner").should("contain", "Assurance visit updated");
        });
    });

    describe("Manage a PDR form", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/2/manage-assurance-visit/36");
        });

        it("cancel button returns user to the assurance visit page", () => {
            cy.get(".govuk-button-group > .govuk-link")
                .should("contain", "Cancel")
                .click();
            cy.url().should(
                "not.contain",
                "/supervision/deputies/2/manage-assurance-visit/36"
            );
            cy.get(".govuk-main-wrapper > header").contains("Assurance visits");
            cy.get(".govuk-button").contains("Add a visit");
        });

        it("shows relevant heading, field labels and submit button", () => {
            cy.get(".govuk-heading-l").contains("Manage PDR");
            cy.get(".govuk-label").contains("PDR due date");
            cy.get(".govuk-label").contains("PDR received date");
            cy.get(".govuk-fieldset__legend").contains("Outcome");
            cy.get(".govuk-label").contains("PDR review date");
            cy.get(".govuk-fieldset__legend").contains("PDR marked as");
            cy.get(".govuk-label").contains("Note");
            cy.get(".govuk-button").contains("Save PDR");
        });

        it("form autofills in existing data", () => {
            cy.get("#visit-report-marked-as-Red").should("be.checked");
        });

        it("allows user to edit and submit the form", () => {
            cy.setCookie("success-route", "manageAssuranceVisit");
            cy.get("#f-report-due-date").type("2021-02-01");
            cy.get("#manage-assurance-visit-form").submit();
            cy.url().should("contain", "/supervision/deputies/2/assurance-visits");
            cy.get(".moj-banner").should("contain", "PDR updated");
        });
    });
});
