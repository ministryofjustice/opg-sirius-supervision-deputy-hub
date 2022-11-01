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
            cy.url().should("include","/supervision/deputies/3/manage-assurance-visit/35");

            cy.get("#f-back-button").click();
            cy.get(".govuk-heading-l").contains("Assurance visits");
        });

        it("should navigate to and from Manage PDR", () => {
            cy.visit("/supervision/deputies/2/assurance-visits");

            cy.get(".govuk-button").contains("Manage PDR").click();
            cy.url().should("include","/supervision/deputies/2/manage-assurance-visit/36");

            cy.get("#f-back-button").click();
            cy.get(".govuk-heading-l").contains("Assurance visits");
        });
    });

    describe("Content", () => {
        it("shows header title and button", () => {
            cy.visit("/supervision/deputies/3/assurance-visits");
            cy.get(".govuk-main-wrapper > header").contains("Assurance visits");
            cy.get(".govuk-button").contains("Add a visit");
        });

        it("should display assurance visit main content", () => {
            cy.visit("/supervision/deputies/3/assurance-visits");
            cy.get("#assurance-visit-details > :nth-child(1) > .govuk-summary-list__key").contains("Assurance type");
            cy.get('#assurance-visit-details > :nth-child(1) > .govuk-summary-list__value').contains("VISIT");
            cy.get("#assurance-visit-details > :nth-child(2) > .govuk-summary-list__key").contains("Requested date");
            cy.get('#assurance-visit-details > :nth-child(2) > .govuk-summary-list__value').contains("22/07/2021");
            cy.get("#assurance-visit-details > :nth-child(3) > .govuk-summary-list__key").contains("Requested by");
            cy.get("#assurance-visit-details > :nth-child(3) > .govuk-summary-list__value").contains("nice user");
            cy.get("#assurance-visit-details > :nth-child(4) > .govuk-summary-list__key").contains("Commissioned date");
            cy.get("#assurance-visit-details > :nth-child(5) > .govuk-summary-list__key").contains("Visitor");
            cy.get("#assurance-visit-details > :nth-child(6) > .govuk-summary-list__key").contains("Report due date");
            cy.get("#assurance-visit-details > :nth-child(7) > .govuk-summary-list__key").contains("Report received date");
            cy.get("#assurance-visit-details > :nth-child(8) > .govuk-summary-list__key").contains("Outcome");
            cy.get("#assurance-visit-details > :nth-child(9) > .govuk-summary-list__key").contains("Report reviewed date");
            cy.get("#assurance-visit-details > :nth-child(10) > .govuk-summary-list__key").contains("Reviewed by");
            cy.get("#assurance-visit-details > :nth-child(11) > .govuk-summary-list__key").contains("Report marked as");
        });

        it("should display PDR main content", () => {
            cy.visit("/supervision/deputies/2/assurance-visits");
            cy.get("#assurance-visit-details > :nth-child(1) > .govuk-summary-list__key").contains("Assurance type");
            cy.get('#assurance-visit-details > :nth-child(1) > .govuk-summary-list__value').contains("PDR");
            cy.get("#assurance-visit-details > :nth-child(2) > .govuk-summary-list__key").contains("Requested date");
            cy.get('#assurance-visit-details > :nth-child(2) > .govuk-summary-list__value').contains("30/06/2022");
            cy.get("#assurance-visit-details > :nth-child(3) > .govuk-summary-list__key").contains("Requested by");
            cy.get("#assurance-visit-details > :nth-child(3) > .govuk-summary-list__value").contains("case manager");
            cy.get("#assurance-visit-details > :nth-child(4) > .govuk-summary-list__key").contains("PDR due date");
            cy.get("#assurance-visit-details > :nth-child(5) > .govuk-summary-list__key").contains("PDR received date");
            cy.get("#assurance-visit-details > :nth-child(6) > .govuk-summary-list__key").contains("Outcome");
            cy.get('#assurance-visit-details > :nth-child(6) > .govuk-summary-list__value').contains("Received");
            cy.get("#assurance-visit-details > :nth-child(7) > .govuk-summary-list__key").contains("PDR reviewed date");
            cy.get("#assurance-visit-details > :nth-child(8) > .govuk-summary-list__key").contains("Reviewed by");
            cy.get("#assurance-visit-details > :nth-child(9) > .govuk-summary-list__key").contains("PDR marked as");
        });
    });

    describe("Add a visit button", () => {
        it("is enabled when latest visit is marked as reviewed", () => {
            cy.visit("/supervision/deputies/3/assurance-visits");
            cy.contains(".govuk-button", "Add a visit").click();
            cy.url().should("contain", "/supervision/deputies/3/add-assurance-visit");
        });

        it("displays warning and does not navigate when latest visit is not marked as reviewed or cancelled", () => {
            cy.visit("/supervision/deputies/2/assurance-visits");
            cy.contains(".govuk-button", "Add a visit").click();
            cy.url().should("contain", "/supervision/deputies/2");
            cy.get("#f-button-disabled-warning")
                .should("be.visible")
                .should("contain", "You cannot add anything until the current assurance process has a review date or is Not received");
        });

        it("is enabled when latest visit is cancelled", () => {
            cy.visit("/supervision/deputies/4/assurance-visits");
            cy.contains(".govuk-button", "Add a visit").click();
            cy.url().should("contain", "/supervision/deputies/4/add-assurance-visit");
        });
    });
});
