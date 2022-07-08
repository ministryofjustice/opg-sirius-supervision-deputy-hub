describe("Assurance Visits", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe("Navigation", () => {
        it("should navigate to Assurance Visits tab", () => {
            cy.visit("/supervision/deputies/3");
            cy.get(".moj-sub-navigation__list").contains("Assurance visits").click();

            cy.url().should("include", "/supervision/deputies/3/assurance-visits");
            cy.get(".govuk-heading-l").contains("Assurance visits");
        });

        it("should navigate to and from Add Assurance Visit", () => {
            cy.visit("/supervision/deputies/3/assurance-visits");

            cy.contains(".govuk-button", "Add a visit").click();
            cy.url().should("include","/supervision/deputies/3/add-assurance-visit");

            cy.get("#f-back-button").click();
            cy.get(".govuk-heading-l").contains("Assurance visits");
        });

        it("should navigate to and from Manage Assurance Visit", () => {
            cy.visit("/supervision/deputies/3/assurance-visits");

            cy.get(".govuk-button").contains("Manage assurance visit").click();
            cy.url().should("include","/supervision/deputies/3/manage-assurance-visit/1");

            cy.get("#f-back-button").click();
            cy.get(".govuk-heading-l").contains("Assurance visits");
        });
    });

    describe("Content", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/2/assurance-visits");
        });

        it("shows header title and button", () => {
            cy.get(".govuk-main-wrapper > header").contains("Assurance visits");
            cy.get(".govuk-button").contains("Add a visit");
        });

        it("should display assurance visit main content", () => {
            cy.get("#assurance-visit-details > :nth-child(1) > .govuk-summary-list__key").contains("Requested date");
            cy.get('#assurance-visit-details > :nth-child(1) > .govuk-summary-list__value').contains("30/06/2022");
            cy.get("#assurance-visit-details > :nth-child(2) > .govuk-summary-list__key").contains("Requested by");
            cy.get("#assurance-visit-details > :nth-child(2) > .govuk-summary-list__value").contains("case manager");
        });
    });

    describe("Add a visit button", () => {
        it("is enabled when latest visit is marked as reviewed", () => {
            cy.visit("/supervision/deputies/3/assurance-visits");
            cy.contains(".govuk-button", "Add a visit").click();
            cy.url().should("contain", "/supervision/deputies/3/add-assurance-visit");
        });

        it("displays warning and does not navigate when latest visit is not marked as reviewed", () => {
            cy.visit("/supervision/deputies/2/assurance-visits");
            cy.contains(".govuk-button", "Add a visit").click();
            cy.url().should("contain", "/supervision/deputies/2");
            cy.get("#f-button-disabled-warning")
                .should("be.visible")
                .should("contain", "You cannot add anything until the current assurance process has a review date and RAG status");
        });
    });
});
