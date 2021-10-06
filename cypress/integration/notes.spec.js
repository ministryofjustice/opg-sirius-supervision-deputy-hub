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
        it("has a add a note page with expected fields", () => {
            cy.get(".govuk-button").should("contain", "Add a note").click();
            cy.get(":nth-child(2) > .govuk-label").should("contain", "Title (required)")
            cy.get(".govuk-character-count > .govuk-form-group > .govuk-label").should("contain", "Note (required)")
            cy.get("#note-info").should("contain", "You have 1000 characters remaining")
            cy.get(".govuk-button").should("contain", "Save note")
            cy.get(".govuk-link").should("contain", "Cancel")
        })

        it("allows me to enter note information which amends character count", () => {
            cy.visit("/supervision/deputies/public-authority/deputy/1/notes/add-note");
            cy.get("#title").type("example note title")
            cy.get("#note").type("example note text")
            cy.get("#note-info").should("contain", "You have 983 characters remaining")
        })

        it("redirects me to main notes page if I cancel adding a note", () => {
            cy.visit("/supervision/deputies/public-authority/deputy/1/notes/add-note");
            cy.get(".govuk-link").should("contain", "Cancel").click()
            cy.get(".main > header").should("contain", "Notes");
            cy.url().should("contain", "/supervision/deputies/public-authority/deputy/1/notes");
        })

        it("shows success banner with Note added message", () => {
            cy.visit("/supervision/deputies/public-authority/deputy/1/notes/add-note");
            cy.get("#title").type("title")
            cy.get("#note").type("note")
            cy.get('form').submit()
            cy.url().should("contain", "/supervision/deputies/public-authority/deputy/1/notes");
            cy.get("body > div > main > div.moj-banner.moj-banner--success > div").should("contain", "Note added");
        })

        it("shows new note on the timeline", () => {
            cy.visit("/supervision/deputies/public-authority/deputy/1/notes/add-note");
            cy.get("#title").type("New note title")
            cy.get("#note").type("Note text entered")
            cy.get('form').submit()
            cy.url().should("contain", "/supervision/deputies/public-authority/deputy/1/notes");
            cy.get(":nth-last-child(1) > .moj-timeline__header > .moj-timeline__title").should("contain", "New note title");
            cy.get(':nth-last-child(1) > pre > .moj-timeline__description').should("contain", "Note text entered");
        })
    })
});