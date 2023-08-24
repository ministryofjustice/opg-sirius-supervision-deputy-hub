const navTabs = [
    ["Deputy details", "/supervision/deputies/1"],
    ["Contacts", "/supervision/deputies/1/contacts"],
    ["Clients", "/supervision/deputies/1/clients"],
    ["Timeline", "/supervision/deputies/1/timeline"],
    ["Notes", "/supervision/deputies/1/notes"],
    ["Tasks", "/supervision/deputies/1/tasks"],
    ["Assurance visits", "/supervision/deputies/1/assurance-visits"],
];

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
                cy.wrap($el).should("contain", navTabs[index][0]);
                cy.wrap($el)
                    .find("a")
                    .should("have.attr", "href")
                    .and("contain", navTabs[index][1]);
            });
    });
});

describe("Accessibility", () => {
   navTabs.forEach(([page, url]) => {
       it(`should render ${page} page accessibly`, () => {
           cy.visit(url);
           cy.injectAxe();
           cy.checkA11y();
       });
   })
});
