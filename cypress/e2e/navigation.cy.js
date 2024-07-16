import navTabs from "../fixtures/navigation.json";

describe("Navigation bar", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/1");
    });

    it("has titles and working nav links for all tabs in the correct order", () => {
        cy.get(".moj-sub-navigation__list")
            .children()
            .each(($el, index) => {
                if(!$el.attr('class', 'hide')) {
                    cy.wrap($el).should("contain", navTabs[index][0]);
                    cy.wrap($el)
                        .find("a")
                        .should("have.attr", "href")
                        .and("contain", navTabs[index][1]);
                }
            });
    });
});
