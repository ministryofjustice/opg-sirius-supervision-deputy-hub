describe("Manage important Information", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe("Navigation", () => {

        it("shows a manage important information button on the dashboard which takes me to the page", () => {
            cy.visit("/supervision/deputies/3");
            cy.get(":nth-child(2) > .govuk-button").should(
                "contain",
                "Manage important information"
            ).click();
            cy.url().should(
                "contain",
                "/supervision/deputies/3/manage-important-information"
            );
            cy.get(".govuk-heading-l").should("contain", "Manage important information");
        });

        it("shows a cancel button which returns me to the dashboard", () => {
            cy.visit("/supervision/deputies/3/manage-important-information");
            cy.get(".govuk-link").should("contain", "Cancel").click();
            cy.url().should(
                "not.contain",
                "/supervision/deputies/3/manage-important-information"
            );
            cy.get(".govuk-heading-l").should("contain", "Deputy details");
        });
    });

    describe("Manage important information form", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/3/manage-important-information");
        });

        it("autofills in existing data", () => {
            cy.get("#complaints-Yes").should('be.checked');
            cy.get("#panel-deputy-yes").should('be.checked');
            cy.get("#annual-billing-Schedule").should('be.checked');
            cy.get("#other-info-note").should('have.text', "Some important information is here");
        });

        it("allows me to edit and submit the form", () => {
            cy.setCookie("success-route", "importantInformation");
            cy.get("#complaints-No").click();
            cy.get("#panel-deputy-no").click();
            cy.get("#annual-billing-Invoice").click();
            cy.get("#other-info-note").clear().type("new data entered into box");
            cy.get(".govuk-button").click()
            cy.url().should(
                "not.contain",
                "/supervision/deputies/3/manage-important-information"
            );
            cy.get(".govuk-heading-l").should("contain", "Deputy details");
            cy.get(".moj-banner").should("contain", "Important information updated")
        });

        it("will show validation errors", () => {
            cy.setCookie("fail-route", "importantInformation");
            cy.get("#other-info-note").clear().type("data that is too long for the box");
            cy.get(".govuk-button").click();
            cy.get(".govuk-error-summary__title").should(
                "contain",
                "There is a problem"
            );
            cy.get(".govuk-list > li > a").should(
                "contain",
                "The other important information must be 1000 characters or fewer"
            );
        });
    });

    describe("Default values", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/4/manage-important-information");
        });

        it("shows the default values when no important information exists", () => {
            cy.get("#complaints-Unknown").should('be.checked');
            cy.get("#panel-deputy-no").should('be.checked');
            cy.get("#annual-billing-Unknown").should('be.checked');
            cy.get("#other-info-note").should('have.text', "");
        });

    });
});
