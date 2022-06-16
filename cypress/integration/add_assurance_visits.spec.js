describe("Add Assurance Visit", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/3/add-assurance-visit");
    });

    describe("Header", () => {
        it("shows content", () => {
            cy.get(".govuk-main-wrapper > header").contains("Add assurance visit");
            cy.get(".govuk-label").contains("Requested date");
            cy.get(".govuk-button").contains("Save assurance visit");
            cy.get(".govuk-button-group > .govuk-link").contains("Cancel");


        });
    });

    describe("Successfully submitting assurance visit form", () => {
        it("should allow me to submit the form", () => {
            cy.get("#f-requested-date").type("2021-02-01");
            cy.get("#add-assurance-visit-form").submit();
            cy.url().should("contain", "/supervision/deputies/3/assurance-visits");
            cy.get(".moj-banner").should("contain", "Assurance process updated");
        });
    });

});
