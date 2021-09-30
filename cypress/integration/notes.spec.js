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

        it("shows the error banner", () => {
            cy.visit("/supervision/deputies/public-authority/deputy/1/notes/add-note");
            cy.get("#title").type("empty string")
            cy.get("#note").type("empty string")
            // cy.intercept('POST', '/api/v1/deputy/1/create-note', {
            //     statusCode: 500,
            // }).as('addNote')
            cy.intercept('POST', '/api/v1/deputy/1/create-note', (req) => {
                req.reply({
                statusCode: 500,
                    body: {
                    name: 'Peter Pan'
                },
                })
            }).as('addNote')
            cy.get(".govuk-button").should("contain", "Save note").click()
            cy.wait('@addNote').should('have.property', 'response.statusCode', 500)
        })

        it("shows error banner can return specific statusCode", () => {
            cy.visit("/supervision/deputies/public-authority/deputy/1/notes/add-note");

            cy.intercept('POST', '/supervision/deputies/public-authority/deputy/1/notes/add-note', {statusCode: 500}).as('addNote')

            cy.get(".govuk-button").should("contain", "Save note").click()
            cy.wait('@addNote')
            cy.get('@addNote').then( xhr => {
                console.log(xhr)
                expect(xhr.response.statusCode).to.equal(500)
            })
        })

    })
});