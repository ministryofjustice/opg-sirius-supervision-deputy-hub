describe("Documents", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/1/documents");
    });

    describe("Content", () => {
        it("shows correct headers", () => {
            cy.get('.govuk-heading-l').contains("Documents");
            cy.get("#doc-name-header").contains("Name and details");
            cy.get("#doc-type-header").contains("Document type");
            cy.get("#doc-added-by-header").contains("Added by");
            cy.get("#doc-direction-header").contains("Direction");
            cy.get("#doc-date-header").contains("Date");
            cy.get("#doc-actions-header").contains("Actions");
        });

        it("shows correct body", () => {
            cy.get("#document-name").contains("Order_documents.pdf");
            cy.contains("a","Order_documents.pdf").invoke('attr','href')
               .should('include','api/v1/documents/5/download')
            cy.get("#document-type").contains("Catch-up call");
            cy.get("#document-added-by").contains("case manager");
            cy.get("#document-direction").contains("Internal");
            cy.get("#document-date").contains("30/05/2024");
        });
    });

    describe("Add document", () => {
        it("allows me to navigate to page from documents and returns to documents on cancel", () => {
            cy.get("#add-a-document-button").contains("Add a document").click();
            cy.url().should("include","/documents/add");
            cy.get("#add-document-cancel-button").contains("Cancel").click();
            cy.url().should("not.include","/add");
        });

        it("allows me to add a document", () => {
            cy.setCookie("success-route", "/add-document/1");
            cy.get("#add-a-document-button").contains("Add a document").click();
            cy.get('.govuk-heading-l').contains("Add a document");
            cy.get('input[type=file]').selectFile('cypress/fixtures/example.json')
            cy.get('#f-type > .govuk-fieldset__legend').contains("Type of document")
            cy.get('#type-ASSURANCE_VISIT').click();
            cy.get('#f-direction > .govuk-fieldset__legend').contains("Direction")
            cy.get('#direction-INCOMING').click();
            cy.get(':nth-child(5) > .govuk-label').contains("Date")
            cy.get('#f-date').type("2021-02-01");
            cy.get('.govuk-character-count > .govuk-form-group > .govuk-label').contains("Notes")
            cy.get('#f-notes').type("Some notes");
            cy.get('#add-document-submit-form').click();
            cy.get('.moj-banner--success').should('be.visible');
            cy.get('.moj-banner--success').contains("Document example.json added");
            cy.url().should("include","?success=addDocument&filename=example.json");
        });
    });

    describe("Validation messages", () => {
        it("shows error messages when submitting invalid data", () => {
            cy.get("#add-a-document-button").contains("Add a document").click();
            let notes = "a";
            cy.get('#f-notes').type(notes.repeat(1001), { delay: 0 });
            cy.get('#add-document-submit-form').click();
            cy.get('.govuk-error-summary').contains("Select a date")
            cy.get('.govuk-error-summary').contains("Select a direction")
            cy.get('.govuk-error-summary').contains("Error uploading the file")
            cy.get('.govuk-error-summary').contains("The note must be 1000 characters or fewer")
            cy.get('.govuk-error-summary').contains("Select a type")
            cy.get('#add-document-form > :nth-child(2)').contains("Error uploading the file")
            cy.get('#f-type').contains("Select a type")
            cy.get('#f-direction').contains("Select a direction")
            cy.get('#add-document-form > :nth-child(5)').contains("Select a date")
            cy.get('.govuk-character-count > .govuk-form-group').contains("The note must be 1000 characters or fewer")
            cy.get('.govuk-character-count > .govuk-form-group').contains("You have 1 character too many")
        });
    });
});
