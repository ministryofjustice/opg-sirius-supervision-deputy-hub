describe("Pagination", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/1/clients?sort=surname:asc");
    });

    it("shows correct number of total clients", () => {
        cy.get("#top-pagination").contains(".moj-pagination__results", 'Showing 1 to 3 of 3 clients');
    });

    it("has the correct options in the per-page dropdown", () => {
        cy.get("#top-pagination .display-rows option").then(options => {
            const values = [...options].map(o => o.value);
            expect(values).to.deep.eq(["?limit=25&page=1", "?limit=50&page=1", "?limit=100&page=1"]);
        });
    });

    it("redirects page to new per-page limit", () => {
        cy.get("#top-pagination .display-rows").select('50');
        cy.location('search').should('contain', 'limit=50');
    });
});
