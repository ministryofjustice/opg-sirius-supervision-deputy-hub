describe("Manage Deputy Contact Details", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe("Navigation", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/3");
        });

        it("should navigate to the 'Manage deputy contact details' page", () => {
            cy.get("[data-cy=manage-deputy-contact-details-btn]").click();
            cy.contains(".govuk-heading-l", "Manage deputy contact details");
        });
    });

    describe("Form functionality", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/3/manage-deputy-contact-details");
        });

        it("should navigate to dashboard on cancel", () => {
            cy.get("[data-cy=cancel-btn]").click();
            cy.contains(".govuk-heading-l", "Deputy details");
        });

        it("should populate fields with current contact details", () => {
            cy.get("input[name=deputy-first-name]").should(
                "have.value",
                "firstname"
            );
            cy.get("input[name=deputy-last-name]").should(
                "have.value",
                "surname"
            );
            cy.get("input[name=address-line-1]").should(
                "have.value",
                "addressLine1"
            );
            cy.get("input[name=address-line-2]").should(
                "have.value",
                "addressLine2"
            );
            cy.get("input[name=address-line-3]").should(
                "have.value",
                "addressLine3"
            );
            cy.get("input[name=town]").should("have.value", "town");
            cy.get("input[name=county]").should("have.value", "county");
            cy.get("input[name=postcode]").should("have.value", "postcode");
            cy.get("input[name=telephone]").should("have.value", "1111111");
            cy.get("input[name=email]").should(
                "have.value",
                "email@something.com"
            );
        });

        it("should show success banner on submit", () => {
            cy.get("form").submit();

            cy.contains(".govuk-heading-l", "Deputy details");
            cy.get(
                "body > div > main > div.moj-banner.moj-banner--success > div"
            ).should("contain", "Deputy details updated");
        });

        it("should show validation errors", () => {
            cy.setCookie("fail-route", "contact-details");
            cy.get("input:not([type=hidden]):not([disabled])").each((input) => {
                cy.wrap(input).clear();
            });
            cy.get("form").submit();
            cy.get(".govuk-error-summary__title").should(
                "contain",
                "There is a problem"
            );
            cy.get(".govuk-error-summary__list").within(() => {
                cy.get("li")
                    .eq(0)
                    .should(
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
                        "The deputy first name is required and can't be empty"
                    );
                cy.get("li")
                    .eq(6)
                    .should(
                        "contain",
                        "The postcode must be 255 characters or fewer"
                    );
                cy.get("li")
                    .eq(7)
                    .should(
                        "contain",
                        "The deputy surname is required and can't be empty"
                    );
                cy.get("li")
                    .eq(8)
                    .should(
                        "contain",
                        "The town or city must be 255 characters or fewer"
                    );
                cy.get("li")
                    .eq(9)
                    .should(
                        "contain",
                        "The telephone number must be 255 characters or fewer"
                    );
            });

            cy.get('#contact-details-form > :nth-child(2).govuk-form-group--error').should("exist");
            cy.get('#f-firstname.govuk-input--error').should("exist");
            cy.get('#contact-details-form > :nth-child(2) > #name-error-isEmpty')
                .should("contain", "The deputy first name is required and can't be empty");

            cy.get('#contact-details-form > :nth-child(3).govuk-form-group--error').should("exist");
            cy.get('#f-surname.govuk-input--error').should("exist");
            cy.get('#contact-details-form > :nth-child(3) > #name-error-isEmpty')
                .should("contain", "The deputy surname is required and can't be empty");

            cy.get('.govuk-fieldset > :nth-child(3).govuk-form-group--error').should("exist");
            cy.get('#f-addressLine1.govuk-input--error').should("exist");
            cy.get('.govuk-fieldset > :nth-child(3) > #name-error-stringLengthTooLong')
                .should("contain", "The building or street must be 255 characters or fewer");

            cy.get('.govuk-fieldset > :nth-child(4).govuk-form-group--error').should("exist");
            cy.get('#f-addressLine2.govuk-input--error').should("exist");
            cy.get('.govuk-fieldset > :nth-child(4) > #name-error-stringLengthTooLong')
                .should("contain", "Address line 2 must be 255 characters or fewer");

            cy.get('.govuk-fieldset > :nth-child(5).govuk-form-group--error').should("exist");
            cy.get('#f-addressLine3.govuk-input--error').should("exist");
            cy.get('.govuk-fieldset > :nth-child(5) > #name-error-stringLengthTooLong')
                .should("contain", "Address line 3 must be 255 characters or fewer");

            cy.get('.govuk-fieldset > :nth-child(6).govuk-form-group--error').should("exist");
            cy.get('#f-town.govuk-input--error').should("exist");
            cy.get('.govuk-fieldset > :nth-child(6) > #name-error-stringLengthTooLong')
                .should("contain", "The town or city must be 255 characters or fewer");

            cy.get('.govuk-fieldset > :nth-child(7).govuk-form-group--error').should("exist");
            cy.get('#f-county.govuk-input--error').should("exist");
            cy.get('.govuk-fieldset > :nth-child(7) > #name-error-stringLengthTooLong')
                .should("contain", "The county must be 255 characters or fewer");

            cy.get('.govuk-fieldset > :nth-child(8).govuk-form-group--error').should("exist");
            cy.get('#f-postcode.govuk-input--error').should("exist");
            cy.get('.govuk-fieldset > :nth-child(8) > #name-error-stringLengthTooLong')
                .should("contain", "The postcode must be 255 characters or fewer");

            cy.get('#contact-details-form > :nth-child(5).govuk-form-group--error').should("exist");
            cy.get('#f-workPhoneNumber.govuk-input--error').should("exist");
            cy.get('#contact-details-form > :nth-child(5) > #name-error-stringLengthTooLong')
                .should("contain", "The telephone number must be 255 characters or fewer");

            cy.get('#contact-details-form > :nth-child(6).govuk-form-group--error').should("exist");
            cy.get('#f-email.govuk-input--error').should("exist");
            cy.get('#contact-details-form > :nth-child(6) > #name-error-stringLengthTooLong')
                .should("contain", "The email must be 255 characters or fewer");

        });
    });

    it("should show 'Deputy name' field when deputy is organisation", () => {
        cy.visit("/supervision/deputies/4/manage-deputy-contact-details");

        cy.get("input[name=organisation-name]").should(
            "have.value",
            "Organisation Ltd"
        );
    });

    describe("Deputy contact details changed timeline event", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/1/timeline");
        });

        it("has a timeline event for when the deputy contact details have changed", () => {
            cy.get("[data-cy=deputy-contact-details-event]").within(() => {
                cy.contains(".moj-timeline__title", "Deputy contact details changed");
                cy.contains(".moj-timeline__byline", "case manager (12345678)");
                cy.get(".moj-timeline__description > .govuk-list").children()
                    .first().should("contain", "Address line 1: Town Hall")
                    .next().should("contain", "Address line 2: City Centre");
            });
        });
    });
});
