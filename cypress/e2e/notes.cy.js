describe("Notes", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe('Navigation', () => {
       it("should navigate to and from the Notes pages", () => {
           cy.visit("/supervision/deputies/1");
           cy.get(".moj-sub-navigation__list").contains("Notes").click();

           cy.url().should("include", "/supervision/deputies/1/notes");
           cy.get(".govuk-heading-l").contains("Notes");

           cy.get(".govuk-button").contains("Add a note").click();
           cy.url().should("include","/supervision/deputies/1/notes/add-note");

           cy.get("#f-back-button").click();
           cy.get(".govuk-heading-l").contains("Notes");
       });
    });

    describe("Viewing notes", () => {
        it("should display notes correctly", () => {
            cy.visit("/supervision/deputies/2/notes");
            cy.get(
                ":nth-last-child(1) > .moj-timeline__header > .moj-timeline__title"
            ).should("contain", "New note title");
            cy.get(
                ":nth-last-child(1) > .moj-timeline__description"
            ).should("contain", "Note text entered");
        });
    })

    describe("Adding a note", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/2/notes/add-note");
        });

        describe("Successfully adding a note", () => {
            it("has a add a note page with expected fields", () => {
                cy.get(":nth-child(2) > .govuk-label").should(
                    "contain",
                    "Title (required)"
                );
                cy.get(
                    ".govuk-character-count > .govuk-form-group > .govuk-label"
                ).should("contain", "Note (required)");
                cy.get("#f-2-note-info").should(
                    "contain",
                    "You have 1000 characters remaining"
                );
                cy.get(".govuk-button").should("contain", "Save note");
                cy.get(".govuk-link").should("contain", "Cancel");
            });

            it("allows me to enter note information which amends character count", () => {
                cy.get("#f-1-title").type("example note title");
                cy.get("#f-2-note").type("example note text");
                cy.get("#f-2-note-info + .govuk-character-count__status").should(
                    "contain",
                    "You have 983 characters remaining"
                );
            });

            it("shows success banner with Note added message", () => {
                cy.get("#f-1-title").type("title");
                cy.get("#f-2-note").type("note");
                cy.get("#add-note-form").submit();
                cy.url().should("contain", "/supervision/deputies/2/notes");
                cy.get(
                    "body > div > main > div.moj-banner.moj-banner--success > div"
                ).should("contain", "Note added");
            });
        });

        it("redirects me to main notes page if I cancel adding a note", () => {
            cy.get(".govuk-button-group > .govuk-link")
                .should("contain", "Cancel")
                .click();
            cy.get(".govuk-heading-l").should("contain", "Notes");
            cy.url().should("contain", "/supervision/deputies/2/notes");
        });

        it("shows error message when submitting invalid data", () => {
            cy.setCookie("fail-route", "notes");
            cy.get("#add-note-form").submit();
            cy.get(".govuk-error-summary__title").should(
                "contain",
                "There is a problem"
            );
            cy.get(".govuk-error-summary__list").within(() => {
                cy.get("li:first").should(
                    "contain",
                    "The title must be 255 characters or fewer"
                );
                cy.get("li:last").should("contain", "Enter a note");
            });
        });

        it("displays an error if the request fails with a 500 error", () => {
            cy.setCookie("fail-route", "500-example");
            cy.setCookie("fail-code", "500");
            cy.get("#add-note-form").submit();
            cy.get(".govuk-heading-l").should(
                "contain",
                "Sorry, there is a problem with the service"
            );
            cy.get(".govuk-body").should(
                "contain",
                " returned 500"
            );
        });
    });
});
