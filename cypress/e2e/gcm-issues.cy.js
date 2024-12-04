describe("Documents", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/1/gcm-issues/open-issues");
    });

    describe("Content", () => {
        it("shows correct headers", () => {
            cy.get('.govuk-heading-l').contains("General Case Manager issues");
            cy.get('.govuk-tabs__list-item--selected').contains('Open issues');

            cy.get(".govuk-table__row").find("th").should("have.length", 5);

            const expected = [
            "",
                "Client",
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
            cy.get('.govuk-table__body > .govuk-table__row > :nth-child(2)').contains('Hamster Person');
            cy.get('.govuk-table__body > .govuk-table__row > :nth-child(3)').contains('PROTeam1 User1');
            cy.get('.govuk-table__body > .govuk-table__row > :nth-child(4)').contains('13/08/2024');
            cy.get('.govuk-table__body > .govuk-table__row > :nth-child(5)').contains('Missing information');
        });

        it("only shows note drop down if the document has notes", () => {
           cy.get(':nth-child(1) > [span="2"] > .govuk-details').contains("Notes").click();
           cy.get(':nth-child(1) > [span="2"] > .govuk-details > .govuk-details__text').should('contain.text', 'Not happy');
        });

        it("allows sorting", () => {
           cy.get('#issue-added-sort').click()
           cy.url().should("include","order-by=createdDate&sort=desc");
           cy.get('#issue-added-sort').click()
           cy.url().should("include","order-by=createdDate&sort=asc");
           cy.get('#issue-sort').click()
           cy.url().should("include","order-by=issueType&sort=asc");
           cy.get('#issue-sort').click()
           cy.url().should("include","order-by=issueType&sort=desc");
        });
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
            cy.get('.govuk-heading-l').contains("Add a GCM issue");
            cy.get('#f-client-case-number').type('12345');
            cy.get('#find-client').click();
            cy.get('#client_name').contains('Hamster Person');
            cy.get('#DEPUTY_FEES_INCORRECT').click();
            cy.get('#f-gcm-note').type('Some thoughts about this issue');
            cy.get('.govuk-button-group > .govuk-button').contains("Save GCM issue").click();
            cy.url().should("not.include","/gcm-issues/add");
            cy.url().should("include","/gcm-issues/open-issues?success=addGcmIssue");
            cy.get('.moj-banner--success').should('be.visible');
            cy.get('.moj-banner--success').should('contain.text', "GCM Issue added");
        });
    });

    describe("Close GCM Issue", () => {
        it("allows me to close a GCM Issue and see correct succes message", () => {
            cy.setCookie("success-route", "/close-gcm-issue/1");
            cy.get('#gcm-issue-1').click();
            cy.get('#close-gcm-issue').click();
            cy.get('.moj-banner--success').contains('You have closed 1 number(s) of GCM issues.')
        });
    });
});
