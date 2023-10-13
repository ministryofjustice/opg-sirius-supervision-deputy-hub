describe("Manage an Assurance Visit", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.setCookie("user", "finance-user");
    });

    describe("Manage an assurance visit form", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/3/manage-assurance/35");
        });

        it("cancel button returns user to the assurance visit page", () => {
            cy.get(".govuk-button-group > .govuk-link")
                .should("contain", "Cancel")
                .click();
            cy.url().should(
                "not.contain",
                "/supervision/deputies/3/manage-assurance/35"
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

        it("form keeps data if validation error", () => {
            cy.setCookie("fail-route", "updateAssurance");
            cy.get("#f-commissioned-date").type("2021-02-01");
            cy.get('#visitor-allocated').select("John Johnson");
            cy.get("#f-report-due-date").type("2021-02-02");
            cy.get("#f-report-received-date").type("2021-02-03");
            cy.get('#visit-outcome-Successful').click();
            cy.get("#f-report-review-date").type("2021-02-04");
            cy.get('#visit-report-marked-as-Green').click();
            cy.get('#f-note').type("This is a test note");

            cy.get("#manage-assurance-form").submit();
            cy.get('.govuk-error-summary').should('be.visible');
            cy.get(".govuk-error-summary__title").should(
                "contain",
                "There is a problem"
            );
            cy.get(".govuk-error-summary__body").should(
                "contain",
                "Report due date must be in the future"
            );

            cy.get('#f-report-due-date.govuk-input--error').should("exist");
            cy.get('#manage-assurance-visit-form :nth-child(4).govuk-form-group--error').should("exist");
            cy.get('#manage-assurance-visit-form :nth-child(4) > #name-error')
                .should("contain", "Report due date must be in the future");

            cy.get("#f-commissioned-date").should("have.value", "2021-02-01");
            cy.get('#visitor-allocated').should("have.value", "John Johnson");
            cy.get("#f-report-due-date").should("have.value","2021-02-02");
            cy.get("#f-report-received-date").should("have.value","2021-02-03");
            cy.get('#visit-outcome-Successful').should("be.checked");
            cy.get("#f-report-review-date").should("have.value","2021-02-04");
            cy.get('#visit-report-marked-as-Green').should("be.checked");
            cy.get('#f-note').contains("This is a test note");
        });

        it("allows user to edit and submit the form", () => {
            cy.setCookie("success-route", "/assurances/2");
            cy.get("#f-commissioned-date").type("2021-02-01");
            cy.get("#manage-assurance-form").submit();
            cy.url().should("contain", "/supervision/deputies/3/assurances");
            cy.get(".moj-banner").should("contain", "Assurance visit updated");
        });
    });

    describe("Manage a PDR form", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/2/manage-assurance/36");
        });

        it("cancel button returns user to the assurance visit page", () => {
            cy.get(".govuk-button-group > .govuk-link")
                .should("contain", "Cancel")
                .click();
            cy.url().should(
                "not.contain",
                "/supervision/deputies/2/manage-assurance/36"
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
            cy.setCookie("success-route", "/assurances/2");
            cy.get("#f-reportDueDate").type("2021-02-01");
            cy.get("#manage-assurance-form").submit();
            cy.url().should("contain", "/supervision/deputies/2/assurances");
            cy.get(".moj-banner").should("contain", "PDR updated");
        });
    });
});
