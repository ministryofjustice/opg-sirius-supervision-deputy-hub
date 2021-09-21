describe("Navigation bar", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/public-authority/deputy/1/");
    });

    const expected = [
        ["Dashboard", "/supervision/deputies/public-authority/deputy/1/"],
        ["Clients", "/supervision/deputies/public-authority/deputy/1/clients"],
        ["Timeline", "/supervision/deputies/public-authority/deputy/1/timeline"],
    ];

    it("has working nav links for different tabs", () => {
        cy.get(".moj-sub-navigation__list")
            .children()
            .each(($el, index) => {
                cy.wrap($el).should("contain", expected[index][0]);
                cy.wrap($el).find('a').should("have.attr", "href").and("contain", expected[index][1]);
            });
    });
});