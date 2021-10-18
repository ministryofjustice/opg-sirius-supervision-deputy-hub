describe("Clients tab", () => {
    beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/supervision/deputies/public-authority/deputy/1/clients");
    });

    it("has a tab header", () => {
      cy.get("h1").should("contain", "Clients");
    });

    it("displays 7 column headings", () => {
        cy.get('.govuk-table__row').find('th').should('have.length', 7)

        const expected = ["Client", "Accomodation type", "Status", "Supervision level", "Visits", "Report due", "CREC"];

        cy.get(".govuk-table__head > .govuk-table__row")
            .children()
            .each(($el, index) => {
                cy.wrap($el).should("contain", expected[index]);
            });
    });

    it("lists clients with active/closed/duplicate orders", () => {
        cy.get(".govuk-table__body > .govuk-table__row").should("have.length", 3);
    });

    it("Clients have a report due dates", () => {
        cy.get(':nth-child(1) > .reports').should("contain", "21/12/2015");
        cy.get(':nth-child(2) > .reports').should("contain", "01/10/2018");
        cy.get(':nth-child(3) > .reports').should("contain", "-");
    });
});
