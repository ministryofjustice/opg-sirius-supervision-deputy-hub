describe("Navigation bar", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/1");
    });

    const expected = [
        ["Deputy details", "/supervision/deputies/1"],
        ["Clients", "/supervision/deputies/1/clients"],
        ["Timeline", "/supervision/deputies/1/timeline"],
        ["Notes", "/supervision/deputies/1/notes"],
        ["Assurance visits", "/supervision/deputies/1/assurance-visits"],
    ];

    it("has titles and working nav links for all tabs in the correct order", () => {
        cy.get(".moj-sub-navigation__list")
            .children()
            .each(($el, index) => {
                cy.wrap($el).should("contain", expected[index][0]);
                cy.wrap($el)
                    .find("a")
                    .should("have.attr", "href")
                    .and("contain", expected[index][1]);
            });
    });
});
