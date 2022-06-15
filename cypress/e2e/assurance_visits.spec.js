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
            cy.setCookie("success-route", "addAssuranceVisit");
            cy.get("#f-requested-date").type("2021-02-01");
            cy.get("#add-assurance-visit-form").submit();
            cy.url().should("contain", "/supervision/deputies/3/assurance-visits");
        });
    });

});
