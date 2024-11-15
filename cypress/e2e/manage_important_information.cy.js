describe("Manage important Information", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe("Navigation", () => {
        it("shows a manage important information button on the dashboard which takes me to the page", () => {
            cy.setCookie("user", "finance-user");
            cy.visit("/supervision/deputies/3");

            cy.get('a[href*="/manage-important-information"]')
                .should('exist')
                .and('be.visible')
                .and('contain.text', 'Manage important information')
                .click();
            cy.url().should(
                "contain",
                "/supervision/deputies/3/manage-important-information"
            );
            cy.get(".govuk-heading-l").should(
                "contain",
                "Manage important information"
            );

            cy.get("#f-back-button").click();
            cy.get(".govuk-heading-l").contains("Deputy details");
        });

        it("shows a cancel button which returns me to the dashboard", () => {
            cy.setCookie("user", "finance-user");
            cy.visit("/supervision/deputies/3/manage-important-information");
            cy.get(".govuk-button-group > .govuk-link")
                .should("contain", "Cancel")
                .click();
            cy.url().should(
                "not.contain",
                "/supervision/deputies/3/manage-important-information"
            );
            cy.get(".govuk-heading-l").should("contain", "Deputy details");
        });
    });

    describe("Manage Pro Deputy important information form", () => {
        beforeEach(() => {
            cy.setCookie("user", "finance-user");
            cy.visit("/supervision/deputies/3/manage-important-information");
        });

        it("autofills in existing data", () => {
            cy.get("#complaints-Yes").should("be.checked");
            cy.get("#panel-deputy-yes").should("be.checked");
            cy.get("#annual-billing-Schedule").should("be.checked");
            cy.get("#f-otherImportantInformation").should(
                "have.text",
                "Some important information is here"
            );
        });

        it("allows me to edit and submit the form", () => {
            cy.setCookie("success-route", "/deputies/1");
            cy.get("#complaints-No").click();
            cy.get("#panel-deputy-no").click();
            cy.get("#annual-billing-Invoice").click();
            cy.get("#f-otherImportantInformation")
                .clear()
                .type("new data entered into box");
            cy.get(".govuk-button").click();
            cy.url().should(
                "not.contain",
                "/supervision/deputies/3/manage-important-information"
            );
            cy.get(".govuk-heading-l").should("contain", "Deputy details");
            cy.get(".moj-banner").should(
                "contain",
                "Important information updated"
            );
        });

        it("will show validation errors", () => {
            cy.setCookie("fail-route", "importantInformation");
            cy.get("#f-otherImportantInformation")
                .clear()
                .type("data that is too long for the box");
            cy.get(".govuk-button").click();
            cy.get(".govuk-error-summary__title").should(
                "contain",
                "There is a problem"
            );
            cy.get(".govuk-list > li > a").should(
                "contain",
                "The other important information must be 1000 characters or fewer"
            );
            cy.get('#f-otherImportantInformation.govuk-input--error').should("exist")
            cy.get('.govuk-character-count > .govuk-form-group--error')
            cy.get('#name-error-stringLengthTooLong').should("contain", "The other important information must be 1000 characters or fewer");
        });
    });

    describe("Manage Pro Deputy default values", () => {
        beforeEach(() => {
            cy.setCookie("user", "finance-user");
            cy.visit("/supervision/deputies/4/manage-important-information");
        });

        it("shows the default values when no important information exists", () => {
            cy.get("body").should("contain", 'Annual billing preference');

            cy.get("#complaints-Unknown").should("be.checked");
            cy.get("#panel-deputy-no").should("be.checked");
            cy.get("#annual-billing-Unknown").should("be.checked");
            cy.get("#f-otherImportantInformation").should("have.text", "");
        });
    });

    describe("Manage Pro Deputy default values as non-Finance Manager", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/4/manage-important-information");
        });

        it("hides the annual billing selection", () => {
            cy.get("body").should("not.contain", 'Annual billing preference');
        });
    });

    describe("Manage Pa Deputy important information form", () => {
        beforeEach(() => {
            cy.setCookie("user", "finance-user");
            cy.visit("/supervision/deputies/1/manage-important-information");
        });

        it("autofills in existing data", () => {
            cy.get("body").should("contain", 'Annual billing preference');

            cy.get("#monthly-spreadsheet-No").should("be.checked");
            cy.get("#independent-visitor-charges-Unknown").should("be.checked");
            cy.get("#bank-charges-Yes").should("be.checked");
            cy.get("#apad-Yes").should("be.checked");
            cy.get("#report-system-CASHFAC").should("be.checked");
            cy.get("#annual-billing-Schedule\\ and\\ Invoice").should(
                "be.checked"
            );
            cy.get("#f-otherImportantInformation").should(
                "have.text",
                "some info for the pa deputy"
            );
        });

        it("allows me to edit and submit the form", () => {
            cy.setCookie("success-route", "/deputies/1");
            cy.get("#independent-visitor-charges-No").click();
            cy.get("#bank-charges-No").click();
            cy.get("#apad-Unknown").click();
            cy.get("#annual-billing-Invoice").click();
            cy.get("#f-otherImportantInformation")
                .clear()
                .type("new data entered into box");
            cy.get(".govuk-button").click();
            cy.url().should(
                "not.contain",
                "/supervision/deputies/3/manage-important-information"
            );
            cy.get(".govuk-heading-l").should("contain", "Deputy details");
            cy.get(".moj-banner").should(
                "contain",
                "Important information updated"
            );
        });

        it("will show validation errors", () => {
            cy.setCookie("fail-route", "importantInformation");
            cy.get("#f-otherImportantInformation")
                .clear()
                .type("data that is too long for the box");
            cy.get(".govuk-button").click();
            cy.get(".govuk-error-summary__title").should(
                "contain",
                "There is a problem"
            );
            cy.get(".govuk-list > li > a").should(
                "contain",
                "The other important information must be 1000 characters or fewer"
            );
            cy.get('#f-otherImportantInformation.govuk-input--error').should("exist");
            cy.get('.govuk-character-count > .govuk-form-group--error').should("exist");
            cy.get('#name-error-stringLengthTooLong').should("contain", "The other important information must be 1000 characters or fewer");
        });
    });

    describe("Manage Pa Deputy important information form as non-Finance Manager", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/1/manage-important-information");
        });

        it("hides the annual billing selection", () => {
            cy.get("body").should("not.contain", 'Annual billing preference');
        });
    });
});
