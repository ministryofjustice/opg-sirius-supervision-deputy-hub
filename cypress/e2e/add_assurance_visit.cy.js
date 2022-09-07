describe("Add Assurance Visit", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/3/add-assurance-visit");
    });

    describe("Header", () => {
        it("shows content", () => {
            cy.get(".govuk-main-wrapper > header").contains("Add assurance visit");
            cy.get(".govuk-fieldset__legend").contains("Assurance type");
            cy.get(".govuk-label").contains("Requested date");
            cy.get(".govuk-button").contains("Save assurance visit");
            cy.get(".govuk-button-group > .govuk-link").contains("Cancel");
        });
    });

    describe("Success submitting assurance visit form", () => {
        it("should allow me to submit the form", () => {
            cy.setCookie("success-route", "addAssuranceVisit");
            cy.get("#assurance-pdr").check();
            cy.get("#f-requested-date").type("2021-02-01");
            cy.get("#add-assurance-visit-form").submit();
            cy.url().should("contain", "/supervision/deputies/3/assurance-visits");
            cy.get(".moj-banner").should("contain", "Assurance process updated");
        });
    });

    describe("Error submitting assurance visit form", () => {
        it("shows error message when submitting invalid data", () => {
        cy.setCookie("fail-route", "addAssuranceVisit");
        cy.get("#add-assurance-visit-form").submit();
        cy.get(".govuk-error-summary__body").should(
            "contain",
            "Enter a requested date"
            );
        });
    });
});
