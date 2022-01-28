describe("Navigation bar", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/deputies/1");
    });

    const expected = [
        ["Deputy details", "/deputies/1"],
        ["Clients", "/deputies/1/clients"],
        [
            "Timeline",
            "/deputies/1/timeline",
        ],
        ["Notes", "/deputies/1/notes"],
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
