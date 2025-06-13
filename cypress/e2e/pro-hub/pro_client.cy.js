describe("Clients tab", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/3/clients?sort=surname:asc");
    });

    it("has a tab header", () => {
        cy.get("h1").should("contain", "Clients");
    });

    it("displays 7 column headings", () => {
        cy.get(".govuk-table__row").find("th").should("have.length", 8);

        const expected = [
            "Client",
            "Accommodation type",
            "Order made date",
            "Status",
            "Supervision level",
            "Visits",
            "Report due",
            "Risk",
        ];

        cy.get(".govuk-table__head > .govuk-table__row")
            .children()
            .each(($el, index) => {
                cy.wrap($el).should("contain", expected[index]);
            });
    });

    it("lists clients with active/closed/duplicate orders", () => {
        cy.get(".govuk-table__body > .govuk-table__row").should("have.length", 3);
    });

    it("Clients surname have been sorted in order of descending", () => {
        cy.get('[aria-sort="ascending"] > a > button').click();
        cy.url().should("contain", "order-by=surname&sort=desc");
        cy.get('[aria-sort="descending"] > a > button').click();
        cy.url().should("contain", "order-by=surname&sort=asc");
    });

    it("Clients report due have been sorted in order of descending", () => {
        cy.get('a:contains("Report")').first().click();
        cy.url().should("contain", "order-by=reportDue&sort=asc");
        cy.get('a:contains("Report")').first().click();
        cy.url().should("contain", "order-by=reportDue&sort=desc");
    });

    it("Clients risk have been sorted in order of descending", () => {
        cy.get('a:contains("Risk")').first().click();
        cy.url().should("contain", "order-by=crec&sort=asc");
        cy.get('a:contains("Risk")').first().click();
        cy.url().should("contain", "order-by=crec&sort=desc");
    });

    it("Clients have a latest visit recorded", () => {
        cy.get(":nth-child(1) > .visit_type").should("contain", "01/01/2000");
        cy.get(":nth-child(1) > .visit_type").should("contain", "Standard visit");
        cy.get(":nth-child(1) > .visit_type").should("contain", "Low risk");
        cy.get(":nth-child(2) > .visit_type").should("contain", "03/03/2020");
        cy.get(":nth-child(2) > .visit_type").should("contain", "Urgent visit");
        cy.get(":nth-child(2) > .visit_type").should("contain", "High risk");
        cy.get(":nth-child(3) > .visit_type").should("contain", "-");
    });
});
