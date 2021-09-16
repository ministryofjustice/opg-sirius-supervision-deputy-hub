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

    it("lists all clients", () => {
        cy.get(".govuk-table__body > .govuk-table__row").should("have.length", 3);

        const expected = ["Luke Crete", "Council Rented", "Active", "Minimal", "", "", 0,
            "Ashley Jackson", "Supervised Sheltered Accommodation", "Active", "General", "", "", 0,
            "John Fearless", "Family Member/Friend's Home (including spouse/civil partner)",
            "Active", "General", "", "", 5];

        cy.get(".govuk-table__body > .govuk-table__row")
            .children()
            .each(($el, index) => {
                cy.wrap($el).should("contain", expected[index]);
            });
    });
});