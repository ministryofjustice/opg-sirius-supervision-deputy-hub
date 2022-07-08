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
});
