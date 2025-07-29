describe("Footer", () => {
    beforeEach(() => {
        cy.visit("/supervision/deputies/1");
    });

    it("should show the accessibility link", () => {
        cy.get('[data-cy="accessibilityStatement"]').should("contain", "Accessibility statement");
    });
});
