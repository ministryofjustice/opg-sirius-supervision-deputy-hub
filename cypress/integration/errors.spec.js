describe("Error handling", () => {
    it("renders the error page when the URL does not match a valid route", () => {
        cy.visit("/supervision/deputies/client/1");
        cy.contains(".govuk-heading-l", "Page not found");
    });
});
