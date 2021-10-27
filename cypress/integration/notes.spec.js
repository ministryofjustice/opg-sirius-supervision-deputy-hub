describe("Notes", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/public-authority/deputy/1/notes");
    });

    describe("Notes timeline", () => {
        it("has a header called notes", () => {
            cy.get(".main > header").should("contain", "Notes");
        })

        it("has a button to add a note which directs me to the add note url", () => {
            cy.get(".govuk-button").should("contain", "Add a note").click();
            cy.url().should("contain", "/supervision/deputies/public-authority/deputy/1/notes/add-note");
        })
    })

    describe("Add a note", () => {
        it("successfully creates a new note", () => {
            cy.visit("/supervision/deputies/public-authority/deputy/1/notes/add-note");

            const title = "A new note";
            const note = "Something to write about";
            cy.get("#title").type(title);
            cy.get("#note").type(note);
            cy.get("#add-note-form").submit();

            cy.url().should("contain", "/supervision/deputies/public-authority/deputy/1/notes");

            cy.get(".moj-timeline").within(() => {
                cy.get(".moj-timeline__title").contains(title);
                cy.get(".moj-timeline__description").contains(note);
            })
        });

        it("shows error message when submitting invalid data", () => {
            cy.setCookie("fail-route", "notes")
            cy.visit("/supervision/deputies/public-authority/deputy/1/notes/add-note");
            cy.get("#add-note-form").submit();
            cy.get(".govuk-error-summary__title").should("contain", "There is a problem");
            cy.get(".govuk-error-summary__list").within(() => {
                cy.get("li:first").should("contain", "The title must be 255 characters or fewer");
                cy.get("li:last").should("contain", "Enter a note");
            })
        });
    });
});