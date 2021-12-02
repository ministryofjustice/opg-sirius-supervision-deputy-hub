describe("Dashboard tab", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/public-authority/deputy/1");
    });

    it("has headers for different sections", () => {
        cy.get("h1").should("contain", "Dashboard");
        cy.get("h2").should("contain", "Team details");
    });

    const expected = [
        ["Deputy name", "Test Organisation"],
        ["Telephone", "0115 876 5574"],
        ["Email", "deputyship@essexcounty.gov.uk"],
        ["Postal address", "Deputyship Team"],
    ];

    it("has rows in tables with accurate keys and values", () => {
        cy.get(".govuk-summary-list")
            .children()
            .each(($el, index) => {
                cy.wrap($el).should("contain", expected[index][0]);
                cy.wrap($el).should("contain", expected[index][1]);
            });
    });

    it("has a href link for email addresses", () => {
        cy.get(".govuk-summary-list__value > a").should("have.attr", "href");
    });

    it("displays warning when no ecm set", () => {
        cy.get(".govuk-list > li").should(
            "contain",
            "An executive case manager has not been assigned. Assign an executive case manager"
        );
    });
});
