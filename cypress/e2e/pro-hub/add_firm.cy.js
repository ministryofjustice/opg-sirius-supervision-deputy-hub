describe("Firm", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe("Adding a firm", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/3/add-firm");
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

            cy.get('#add-firm-form > :nth-child(2).govuk-form-group--error').should("exist");
            cy.get('#f-firmName.govuk-input--error').should("exist");
            cy.get('#name-error-isEmpty').should("contain", "The firm name is required and can't be empty");

            cy.get('.govuk-fieldset > :nth-child(2).govuk-form-group--error').should("exist");
            cy.get('#f-addressLine1.govuk-input--error').should("exist");
            cy.get(':nth-child(2) > #name-error-stringLengthTooLong').should("contain", "The building or street must be 255 characters or fewer");

            cy.get('.govuk-fieldset > :nth-child(3).govuk-form-group--error').should("exist");
            cy.get('#f-addressLine2.govuk-input--error').should("exist");
            cy.get(':nth-child(3) > #name-error-stringLengthTooLong').should("contain", "Address line 2 must be 255 characters or fewer");

            cy.get('.govuk-fieldset > :nth-child(4).govuk-form-group--error').should("exist");
            cy.get('#f-addressLine3.govuk-input--error').should("exist");
            cy.get(':nth-child(4) > #name-error-stringLengthTooLong').should("contain", "Address line 3 must be 255 characters or fewer");

            cy.get('.govuk-fieldset > :nth-child(5).govuk-form-group--error').should("exist");
            cy.get('#f-town.govuk-input--error').should("exist");
            cy.get(':nth-child(5) > #name-error-stringLengthTooLong').should("contain", "The town or city must be 255 characters or fewer");

            cy.get('.govuk-fieldset > :nth-child(6).govuk-form-group--error').should("exist");
            cy.get('#f-county.govuk-input--error').should("exist");
            cy.get(':nth-child(6) > #name-error-stringLengthTooLong').should("contain", "The county must be 255 characters or fewer");

            cy.get('.govuk-fieldset > :nth-child(7).govuk-form-group--error').should("exist");
            cy.get('#f-postcode.govuk-input--error').should("exist");
            cy.get(':nth-child(7) > #name-error-stringLengthTooLong').should("contain", "The postcode must be 255 characters or fewer");

            cy.get('#add-firm-form > :nth-child(4).govuk-form-group--error').should("exist");
            cy.get('#f-phoneNumber.govuk-input--error').should("exist");
            cy.get('#add-firm-form > :nth-child(4) > #name-error-stringLengthTooLong').should("contain", "The telephone number must be 255 characters or fewer");

            cy.get('#add-firm-form > :nth-child(5).govuk-form-group--error').should("exist");
            cy.get('#f-email.govuk-input--error').should("exist");
            cy.get('#add-firm-form > :nth-child(5) > #name-error-stringLengthTooLong').should("contain", "The email must be 255 characters or fewer");
        });

        it("allows me to fill in and submit the firm form", () => {
            cy.setCookie("success-route", "/firms/2");
            cy.get("#f-firmName").type("The Firm Name");
            cy.get("#add-firm-form").submit();
            cy.get(".moj-banner").should("contain", "Firm added");
            cy.get(".govuk-heading-l").should("contain", "Deputy details");
        });
    });
});
