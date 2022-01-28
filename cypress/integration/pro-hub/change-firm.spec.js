describe("Change Firm", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe("Changing a firm", () => {
        beforeEach(() => {
            cy.visit(
                "/supervision/deputies/public-authority/3/change-firm"
            );
        });

        it("shows title for page", () => {
            cy.get(".govuk-grid-column-full > header").should(
                "contain",
                "Change firm"
            );
        });

        it("shows current firm name", () => {
            cy.get(".govuk-body").should("contain", "Current firm");
        });

        it("has a save button that can redirect to add-note page", () => {
            cy.get("#new-firm").click();
            cy.get(".govuk-button")
                .should("contain", "Save and continue")
                .click();
            cy.url().should(
                "contain",
                "/supervision/deputies/public-authority/3/add-firm"
            );
        });

        it("has a cancel button that can redirect to deputy details page", () => {
            cy.get(".govuk-link").should("contain", "Cancel").click();
            cy.url().should(
                "contain",
                "/supervision/deputies/public-authority/3"
            );
        });
    });
});
