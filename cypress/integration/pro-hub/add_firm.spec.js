describe("Firm", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe("Adding a firm", () => {
        beforeEach(() => {
            cy.visit(
                "/supervision/deputies/public-authority/3/add-firm"
            );
        });

        it("shows error message when submitting invalid data", () => {
            cy.setCookie("fail-route", "firm");
            cy.get("#add-firm-form").submit();
            cy.get(".govuk-error-summary__title").should(
                "contain",
                "There is a problem"
            );
            cy.get(".govuk-error-summary__list").within(() => {
                cy.get("li:first").should(
                    "contain",
                    "The building or street must be 255 characters or fewer"
                );
                cy.get("li")
                    .eq(1)
                    .should(
                        "contain",
                        "Address line 2 must be 255 characters or fewer"
                    );
                cy.get("li")
                    .eq(2)
                    .should(
                        "contain",
                        "Address line 3 must be 255 characters or fewer"
                    );
                cy.get("li")
                    .eq(3)
                    .should(
                        "contain",
                        "The county must be 255 characters or fewer"
                    );
                cy.get("li")
                    .eq(4)
                    .should(
                        "contain",
                        "The email must be 255 characters or fewer"
                    );
                cy.get("li")
                    .eq(5)
                    .should(
                        "contain",
                        "The firm name is required and can't be empty"
                    );
                cy.get("li")
                    .eq(6)
                    .should(
                        "contain",
                        "The telephone number must be 255 characters or fewer"
                    );
                cy.get("li")
                    .eq(7)
                    .should(
                        "contain",
                        "The postcode must be 255 characters or fewer"
                    );
                cy.get("li")
                    .eq(8)
                    .should(
                        "contain",
                        "The town or city must be 255 characters or fewer"
                    );
            });
        });

        it("allows me to fill in and submit the firm form", () => {
            cy.setCookie("success-route", "firm");
            cy.get("#f-firmName").type("The Firm Name");
            cy.get("#add-firm-form").submit();
            cy.get(".govuk-heading-l").should("contain", "Deputy details");
        });
    });
});
