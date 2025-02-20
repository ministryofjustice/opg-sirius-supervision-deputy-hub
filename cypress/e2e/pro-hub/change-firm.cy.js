describe("Change Firm", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe("Changing a firm", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/3");
            cy.get("#change-firm").click();
        });

        it("shows title for page", () => {
            cy.get(".govuk-grid-column-full > header").should(
                "contain",
                "Change firm"
            );
        });

        it("shows current firm name", () => {
            cy.get(".govuk-body").should("contain", "Current firm");
        });

        it("has a save button that can redirect to add-note page", () => {
            cy.get("#new-firm").click();
            cy.get(".govuk-button")
                .should("contain", "Save and continue")
                .click();
            cy.url().should("contain", "/supervision/deputies/3/add-firm");
        });

        it("has a cancel button that can redirect to deputy details page", () => {
            cy.get(".govuk-button-group > .govuk-link")
                .should("contain", "Cancel")
                .click();
            cy.url().should("contain", "/supervision/deputies/3");
        });
    });

    describe("Changing a firm to use existing firm", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/3/change-firm");
        });

        it("has a dropdown with the existing firm options", () => {
            cy.get("#f-existing-firm").click();
            cy.get("#select-existing-firm-dropdown > .govuk-label").should(
                "contain",
                "Enter a firm name or number"
            );
            cy.get("#select-existing-firm").click().type("Firm");
            cy.get("#select-existing-firm__listbox")
                .find("li")
                .should("have.length", 2);
        });

        it("will redirect and show success banner when deputy allocated to firm", () => {
            cy.setCookie("success-route", "/firms/1");
            cy.get("#f-existing-firm").click();
            cy.get("#select-existing-firm-dropdown > .govuk-label").should(
                "contain",
                "Enter a firm name or number"
            );
            cy.get("#select-existing-firm").click().type("Great");
            cy.contains(
                "#select-existing-firm__option--0",
                "Great Firm Corp - 1000002"
            ).click();
            cy.get("#existing-firm-or-new-firm-form").submit();
            cy.get(".moj-banner").should("contain", "Firm changed to");
            cy.get("h1").should("contain", "Deputy details");
        });

        it("will allow searching based on firm id", () => {
            cy.setCookie("success-route", "/firms/1");
            cy.get("#f-existing-firm").click();
            cy.get("#select-existing-firm-dropdown > .govuk-label").should(
                "contain",
                "Enter a firm name or number"
            );
            cy.get("#select-existing-firm").click().type("1000002");
            cy.contains(
                "#select-existing-firm__option--0",
                "Great Firm Corp - 1000002"
            ).click();
            cy.get("#existing-firm-or-new-firm-form").submit();
            cy.get(".moj-banner").should("contain", "Firm changed to");
            cy.get("h1").should("contain", "Deputy details");
        });

        it("will show a validation error if no options available", () => {
            cy.setCookie("fail-route", "allocateToFirm");
            cy.get("#f-existing-firm").click();
            cy.get("#select-existing-firm")
                .click()
                .type("Unknown option for firm name");
            cy.get("#existing-firm-or-new-firm-form").submit();
            cy.get(".govuk-error-summary__title").should(
                "contain",
                "There is a problem"
            );
            cy.get(".govuk-error-summary__list").within(() => {
                cy.get("li:first").should(
                    "contain",
                    "Enter a firm name or number"
                );
            });
        });

        it("will show a validation error if form submitted when autocomplete empty", () => {
            cy.setCookie("fail-route", "allocateToFirm");
            cy.get("#f-existing-firm").click();
            cy.get("#existing-firm-or-new-firm-form").submit();
            cy.get(".govuk-error-summary__title").should(
                "contain",
                "There is a problem"
            );
            cy.get(".govuk-error-summary__list").within(() => {
                cy.get("li:first").should(
                    "contain",
                    "Enter a firm name or number"
                );
            });
        });
    });

    describe("Change Firm timeline event", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/3/timeline");
        });

        it("has a timeline event for when the firm is changed", () => {
            cy.get("[data-cy=deputy-allocated-firm-event]").first().within(() => {
                cy.contains(".moj-timeline__title", "Deputy firm updated");
                cy.contains(".moj-timeline__byline", "case manager (12345678)");
                cy.get(".moj-timeline__description > .govuk-list").children()
                    .first().should("contain", "New firm:")
                    .next().should("contain", "Old firm:");
            });
        });

        it("does not show the firm on timeline event if its the first firm set", () => {
            cy.contains("[data-cy=deputy-allocated-firm-event] > .moj-timeline__description > .govuk-list > li", "My First Firm")
                .parent().should("have.length", 1);
        });
    });
});
