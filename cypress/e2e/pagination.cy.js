describe("Pagination", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/1/clients?sort=surname:asc");
    });

    it("shows correct number of total clients", () => {
        cy.get("#top-pagination").contains(".moj-pagination__results", 'Showing 1 to 3 of 3 clients')
    })

    it("can select 25 from task view value dropdown", () => {
        cy.get("#top-pagination .display-rows").select('25')
            .invoke('val').should('contain', 'limit=25')
    })

    it("can select 50 from task view value dropdown", () => {
        cy.get("#top-pagination .display-rows").select('50')
            .invoke('val').should('contain', 'limit=50')
    })

    it("can select 100 from task view value dropdown", () => {
        cy.get("#top-pagination .display-rows").select('100')
            .invoke('val').should('contain', 'limit=100')
    })
});
