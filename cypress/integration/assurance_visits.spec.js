describe("Manage Assurance Visits", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/3/assurance-visits");
    });

    describe("Header", () => {
        it("shows title and button", () => {
            cy.get(".govuk-main-wrapper > header").contains("Assurance visits");
            cy.get(".govuk-button").contains("Add a visit");
        });
    });

    describe("Main content", () => {
        it("should display assurance visit main content", () => {
            cy.get(":nth-child(5) > .govuk-grid-column-one-half > #assurance-visit-details > :nth-child(1) > .govuk-summary-list__key").contains("Requested date");
            cy.get(':nth-child(5) > .govuk-grid-column-one-half > #assurance-visit-details > :nth-child(1) > .govuk-summary-list__value').contains("30/06/2022");
            cy.get(":nth-child(5) > .govuk-grid-column-one-half > #assurance-visit-details > :nth-child(2) > .govuk-summary-list__key").contains("Requested by");
            cy.get(":nth-child(5) > .govuk-grid-column-one-half > #assurance-visit-details > :nth-child(2) > .govuk-summary-list__value").contains("case manager");
        });
    });


    describe("Add a visit button", () => {
        it("should display assurance visit main content", () => {
            cy.get(".govuk-button").click();
            cy.url().should("contain", "/supervision/deputies/3/add-assurance-visit");
        });
    });
});
