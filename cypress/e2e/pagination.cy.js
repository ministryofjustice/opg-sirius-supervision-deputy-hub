describe("Pagination", () => {
    before(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");

    });

    it("is visible on the Deputy client list page", () => {
        cy.visit("/supervision/deputies/1/clients?order-by=surname&sort=asc");
        cy.get("#top-pagination").should("exist");
        cy.get("#bottom-pagination").should("exist");
        cy.get(".moj-pagination__results").should("contain.text", "Showing 1 to 3 of 3 clients")
    })
});
