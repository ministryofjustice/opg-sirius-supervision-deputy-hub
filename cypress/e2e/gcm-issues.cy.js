describe("Documents", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/1/gcm-issues");
    });

    describe("Content", () => {
        it("shows correct headers", () => {
            cy.get('.govuk-heading-l').contains("General Case Manager issues");
            cy.get('#tab_open').contains('Open issues').should("have.attr", "aria-selected");

            cy.get(".govuk-table__row").find("th").should("have.length", 6);

            const expected = [
                "",
                "Case number",
                "Client name",
                "General Case Manager",
                "Issue added",
                "Issue"
            ];

            cy.get(".govuk-table__head > .govuk-table__row")
                .children()
                .each(($el, index) => {
                    cy.wrap($el).should("contain", expected[index]);
                });
        });

        it("shows correct body", () => {
            cy.get('.govuk-table__body > .govuk-table__row > :nth-child(2)').contains('48217682')
            cy.get('.govuk-table__body > .govuk-table__row > :nth-child(3)').contains('Hamster Person');
            cy.get('.govuk-table__body > .govuk-table__row > :nth-child(4)').contains('PROTeam1 User1');
            cy.get('.govuk-table__body > .govuk-table__row > :nth-child(5)').contains('13/08/2024');
            cy.get('.govuk-table__body > .govuk-table__row > :nth-child(6)').contains('Missing information');


            // cy.contains("a","Order_documents.pdf").invoke('attr','href')
            //    .should('include','api/v1/documents/5/download')

        });

         // it("only shows note drop down if the document has notes", () => {
         //    cy.get(':nth-child(2) > #document-name').contains("important_file.png");
         //    cy.get(':nth-child(2) > #notes-summary').should("not.exist");
         // });
    });

    describe("Add GCM Issue", () => {
        it("allows me to navigate to page from GCM issues and returns to GCM issues on cancel", () => {
            cy.get("#add-a-gcm-issue").contains("Add a GCM issue").click();
            cy.url().should("include","/gcm-issues/add");
            cy.get("a.govuk-link").contains("Cancel").click();
            cy.url().should("not.include","/add");
        });

        it("allows me to add a GCM issue", () => {
            cy.setCookie("success-route", "/add-gcm-issue/1");
            cy.get("#add-a-gcm-issue").contains("Add a GCM issue").click();
            cy.get('.govuk-heading-l').contains("Add a document");
            cy.get('input[type=file]').selectFile('cypress/fixtures/example.json')
            cy.get('#f-documentType > .govuk-fieldset__legend').contains("Type of document")
            cy.get('#type-ASSURANCE_VISIT').click();
            cy.get('#f-documentDirection > .govuk-fieldset__legend').contains("Direction")
            cy.get('#direction-INCOMING').click();
            cy.get(':nth-child(5) > .govuk-label').contains("Date")
            cy.get('#f-documentDate').type("2021-02-01");
            cy.get('.govuk-character-count > .govuk-form-group > .govuk-label').contains("Notes")
            cy.get('#f-notes').type("Some notes");
            cy.get('#add-document-submit-form').click();
            cy.get('.moj-banner--success').should('be.visible');
            cy.get('.moj-banner--success').contains("Document example.json added");
            cy.url().should("include","?success=addDocument&filename=example.json");
        });
    });

    describe("Replace document", () => {

        it("allows me to navigate to page from documents and returns to documents on cancel", () => {
            cy.get('.govuk-table__body > :nth-child(1)').contains("Replace").click();
            cy.url().should("include","/documents/5/replace");
            cy.get("#replace-document-cancel-button").contains("Cancel").click();
            cy.url().should("not.include","/replace");
        });

        it("contains information on document to be changed", () => {
            cy.get('.govuk-table__body > :nth-child(1)').contains("Replace").click();
            cy.get('#main-content > header').contains("Replace a document");
            cy.get('.govuk-table').contains("Current document details");
            cy.get('.govuk-table').contains("Order_documents.pdf");
            cy.get('.govuk-table').contains("Catch-up call");
            cy.get('.govuk-table').contains("Internal");
            cy.get('.govuk-table').contains("30/05/2024");
            cy.get('.govuk-table').contains("test");
        });

        it("allows me to replace document", () => {
            cy.get('.govuk-table__body > :nth-child(1)').contains("Replace").click();
            cy.setCookie("success-route", "/replace-document/1");
            cy.get('input[type=file]').selectFile('cypress/fixtures/example.json')
            cy.get('#f-documentType > .govuk-fieldset__legend').contains("Type of document")
            cy.get('#type-ASSURANCE_VISIT').click();
            cy.get('#f-documentDirection > .govuk-fieldset__legend').contains("Direction")
            cy.get('#direction-INCOMING').click();
            cy.get(':nth-child(4) > .govuk-label').contains("Date")
            cy.get('#f-documentDate').type("2021-02-01");
            cy.get('.govuk-character-count > .govuk-form-group > .govuk-label').contains("Notes")
            cy.get('#f-notes').type("Some notes");
            cy.get('#replace-document-submit-form').click();
            cy.get('.moj-banner--success').should('be.visible');
            cy.get('.moj-banner--success').contains("Document Order_documents.pdf has been replaced by example.json");
            cy.url().should("include","?success=replaceDocument&previousFilename=Order_documents.pdf&filename=example.json");
        });
    });

    describe("Validation messages", () => {
        it("shows error messages when submitting invalid data when adding document", () => {
            cy.get("#add-a-document-button").contains("Add a document").click();
            cy.get('#add-document-submit-form').click();
            cy.get('.govuk-error-summary').contains("Select a file to attach")
            cy.get('#add-document-form > :nth-child(2)').contains("Select a file to attach")
        });
        it("shows error messages when submitting invalid data when replacing document", () => {
            cy.get('.govuk-table__body > :nth-child(1)').contains("Replace").click();
            cy.get('#replace-document-submit-form').click();
            cy.get('.govuk-error-summary').contains("Select a file to attach")
            cy.get('#replace-document-form > :nth-child(2)').contains("Select a file to attach")
        });
    });
});
