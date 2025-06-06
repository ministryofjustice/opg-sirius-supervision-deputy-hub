describe("Contacts", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe("Navigation", () => {
        it("should navigate to Contacts tab", () => {
            cy.visit("/supervision/deputies/1");
            cy.get(".moj-sub-navigation__list").contains("Contacts").click();

            cy.url().should("include", "/supervision/deputies/1/contacts");
            cy.get(".govuk-heading-l").contains("Contacts");
        });
    });

    describe("Adding a Contact", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/1/contacts/add-contact");
        });

        it("shows content", () => {
            cy.get(".govuk-main-wrapper > header").contains("Add new contact");
            cy.get("#contact-form > :nth-child(2) > .govuk-label").contains("Name");
            cy.get("#contact-form > :nth-child(3) > .govuk-label").contains("Job title (optional)");
            cy.get("#contact-form > :nth-child(4) > .govuk-label").contains("Email");
            cy.get("#contact-form > :nth-child(5) > .govuk-label").contains("Phone");
            cy.get("#contact-form > :nth-child(6) > .govuk-label").contains("Other phone (optional)");
            cy.get("#contact-form > :nth-child(7) .govuk-label").contains("Notes (optional)");
            cy.get("#f-isNamedDeputy > .govuk-fieldset__legend").contains("Named deputy");
            cy.get("#f-isMainContact > .govuk-fieldset__legend").contains("Main contact");
            cy.get("#f-isMonthlySpreadsheetRecipient > .govuk-fieldset__legend").contains(
                "Monthly spreadsheet recipient",
            );
            cy.get(".govuk-button").contains("Save contact");
            cy.get(".govuk-button-group > .govuk-link").contains("Cancel");
        });

        it("should allow me to submit the form", () => {
            cy.setCookie("success-route", "/contacts/1");
            cy.get("#f-contactName").type("Test Contact");
            cy.get("#f-email").type("test@email.com");
            cy.get("#f-phoneNumber").type("0123456789");
            cy.get("#is-named-deputy-no").click();
            cy.get("#is-main-contact-no").click();
            cy.get("#contact-form").submit();
            cy.url().should("contain", "/supervision/deputies/1/contacts?success=newContact");
            cy.get(".moj-banner").should("contain", "Contact added");
        });

        it("shows error message when submitting empty data", () => {
            cy.setCookie("fail-route", "manageContactEmpty");
            cy.get("#contact-form").submit();
            cy.get(".govuk-error-summary__body").should("contain", "Enter a name");

            cy.get("#f-contactName.govuk-input--error").should("exist");
            cy.get("#contact-form > :nth-child(2).govuk-form-group--error").should("exist");
            cy.get(":nth-child(2) > #name-error-isEmpty").should("contain", "Enter a name");
        });

        it("shows error message when submitting invalid data", () => {
            cy.setCookie("fail-route", "manageContactInvalid");
            cy.get("#contact-form").submit();
            cy.get(".govuk-error-summary__body")
                .should("contain", "The name must be 255 characters or fewer")
                .should("contain", "The job title must be 255 characters or fewer")
                .should("contain", "Enter an email address in the correct format, like name@example.com")
                .should("contain", "The telephone number must be 255 characters or fewer")
                .should("contain", "The other telephone number must be 255 characters or fewer")
                .should("contain", "The note must be 255 characters or fewer");

            cy.get("#f-contactName.govuk-input--error").should("exist");
            cy.get("#contact-form > :nth-child(2).govuk-form-group--error").should("exist");
            cy.get(":nth-child(2) > #name-error-stringLengthTooLong").should(
                "contain",
                "The name must be 255 characters or fewer",
            );

            cy.get("#f-jobTitle.govuk-input--error").should("exist");
            cy.get("#contact-form > :nth-child(3).govuk-form-group--error").should("exist");
            cy.get(":nth-child(3) > #name-error-stringLengthTooLong").should(
                "contain",
                "The job title must be 255 characters or fewer",
            );

            cy.get("#f-email.govuk-input--error").should("exist");
            cy.get("#contact-form > :nth-child(4).govuk-form-group--error").should("exist");
            cy.get(":nth-child(4) > #name-error-emailAddressInvalidFormat").should(
                "contain",
                "Enter an email address in the correct format, like name@example.com",
            );

            cy.get("#f-phoneNumber.govuk-input--error").should("exist");
            cy.get("#contact-form > :nth-child(5).govuk-form-group--error").should("exist");
            cy.get(":nth-child(5) > #name-error-stringLengthTooLong").should(
                "contain",
                "The telephone number must be 255 characters or fewer",
            );

            cy.get("#f-otherPhoneNumber.govuk-input--error").should("exist");
            cy.get("#contact-form > :nth-child(6).govuk-form-group--error").should("exist");
            cy.get(":nth-child(6) > #name-error-stringLengthTooLong").should(
                "contain",
                "The other telephone number must be 255 characters or fewer",
            );

            cy.get("#f-contactNotes.govuk-input--error").should("exist");
            cy.get(".govuk-character-count > .govuk-form-group.govuk-form-group--error").should("exist");
            cy.get(".govuk-character-count > .govuk-form-group > #name-error-stringLengthTooLong").should(
                "contain",
                "The note must be 255 characters or fewer",
            );
        });
    });

    describe("Updating a Contact", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/1/contacts");
            cy.get(":nth-child(2) > :nth-child(4) > .govuk-button--secondary").click();
        });

        it("shows content", () => {
            cy.get(".govuk-main-wrapper > header").contains("Manage contact");
            cy.get("#contact-form > :nth-child(2) > .govuk-label").contains("Name");
            cy.get("#contact-form > :nth-child(3) > .govuk-label").contains("Job title (optional)");
            cy.get("#contact-form > :nth-child(4) > .govuk-label").contains("Email");
            cy.get("#contact-form > :nth-child(5) > .govuk-label").contains("Phone");
            cy.get("#contact-form > :nth-child(6) > .govuk-label").contains("Other phone (optional)");
            cy.get("#contact-form > :nth-child(7) .govuk-label").contains("Notes (optional)");
            cy.get("#f-isNamedDeputy > .govuk-fieldset__legend").contains("Named deputy");
            cy.get("#f-isMainContact > .govuk-fieldset__legend").contains("Main contact");
            cy.get("#f-isMonthlySpreadsheetRecipient > .govuk-fieldset__legend").contains(
                "Monthly spreadsheet recipient",
            );
            cy.get(".govuk-button").contains("Save contact");
            cy.get(".govuk-button-group > .govuk-link").contains("Cancel");
        });

        it("should allow me to submit the form", () => {
            cy.setCookie("success-route", "/contacts/1");
            cy.get("#f-contactName").type("{selectAll}{backspace}John Smith");
            cy.get("#contact-form").submit();
            cy.url().should("contain", "/supervision/deputies/1/contacts?success=updatedContact");
            cy.get(".moj-banner").should("contain", "John Smith's details updated");
        });

        it("shows error message when submitting empty data", () => {
            cy.setCookie("fail-route", "manageContactEmpty");
            cy.get("#contact-form").submit();
            cy.get(".govuk-error-summary__body")
                .should("contain", "Enter a name")
                .should("contain", "Select whether this contact is a main contact")
                .should("contain", "Select whether this contact is the named deputy");

            cy.get("#f-contactName.govuk-input--error").should("exist");
            cy.get(".govuk-form-group--error").should("exist");
            cy.get("#name-error-isEmpty").should("contain", "Enter a name");
        });

        it("shows error message when submitting invalid data", () => {
            cy.setCookie("fail-route", "manageContactInvalid");
            cy.get("#contact-form").submit();
            cy.get(".govuk-error-summary__body")
                .should("contain", "The name must be 255 characters or fewer")
                .should("contain", "The job title must be 255 characters or fewer")
                .should("contain", "Enter an email address in the correct format, like name@example.com")
                .should("contain", "The telephone number must be 255 characters or fewer")
                .should("contain", "The other telephone number must be 255 characters or fewer")
                .should("contain", "The note must be 255 characters or fewer");

            cy.get(".govuk-form-group--error").should("exist");
            cy.get("#f-contactName.govuk-input--error").should("exist");
            cy.get(":nth-child(2) > #name-error-stringLengthTooLong").should(
                "contain",
                "The name must be 255 characters or fewer",
            );
            cy.get("#f-jobTitle.govuk-input--error").should("exist");
            cy.get(":nth-child(3) > #name-error-stringLengthTooLong").should(
                "contain",
                "The job title must be 255 characters or fewer",
            );
            cy.get("#f-email.govuk-input--error").should("exist");
            cy.get("#name-error-emailAddressInvalidFormat").should(
                "contain",
                "Enter an email address in the correct format, like name@example.com",
            );
            cy.get("#f-phoneNumber.govuk-input--error").should("exist");
            cy.get(":nth-child(5) > #name-error-stringLengthTooLong").should(
                "contain",
                "The telephone number must be 255 characters or fewer",
            );
            cy.get("#f-otherPhoneNumber.govuk-input--error").should("exist");
            cy.get(":nth-child(6) > #name-error-stringLengthTooLong").should(
                "contain",
                "The other telephone number must be 255 characters or fewer",
            );
            cy.get("#f-contactNotes.govuk-input--error").should("exist");
            cy.get(".govuk-character-count > .govuk-form-group--error").should("exist");
            cy.get(".govuk-character-count > .govuk-form-group > #name-error-stringLengthTooLong").should(
                "contain",
                "The note must be 255 characters or fewer",
            );
        });
    });

    describe("List Contacts", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/1/contacts");
        });

        it("shows header title and button", () => {
            cy.get(".govuk-main-wrapper > header").contains("Contacts");
            cy.get(".govuk-button").contains("Add new contact");
        });

        it("displays 4 column headings", () => {
            cy.get(".govuk-table__row").find("th").should("have.length", 4);

            const expected = ["Contact", "Contact details", "Notes", "Action"];

            cy.get(".govuk-table__head > .govuk-table__row")
                .children()
                .each(($el, index) => {
                    cy.wrap($el).should("contain", expected[index]);
                });
        });

        it("should display contact data", () => {
            cy.get(":nth-child(1) > :nth-child(1) > .name").contains("Minimal Contact");
            cy.get(":nth-child(1) > :nth-child(2) > .email > a").contains("email@test.com");
            cy.get(":nth-child(1) > :nth-child(2) > .phone-number").contains("0123456789");
            cy.get(":nth-child(1) > :nth-child(4) > .govuk-button--secondary").contains("Manage contact");
            cy.get(":nth-child(1) > :nth-child(4) > .govuk-button--warning").contains("Delete contact");

            cy.get(":nth-child(2) > :nth-child(1) > .name").contains("Test Contact");
            cy.get(":nth-child(2) > :nth-child(1) > :nth-child(2)").contains("Main contact");
            cy.get(":nth-child(2) > :nth-child(1) > :nth-child(3)").contains("Named deputy");
            cy.get(":nth-child(2) > :nth-child(1)").contains("Monthly spreadsheet recipient");
            cy.get(":nth-child(2) > :nth-child(1) > .job-title").contains("Software Tester");
            cy.get(":nth-child(2) > :nth-child(2) > .email > a").contains("test@email.com");
            cy.get(":nth-child(2) > :nth-child(2) > .phone-number").contains("0123456789");
            cy.get(":nth-child(2) > :nth-child(2) > .other-phone-number").contains("9876543210");
            cy.get(":nth-child(2) > :nth-child(3) > .notes").contains("This is a test");
            cy.get(":nth-child(2) > :nth-child(4) > .govuk-button--secondary").contains("Manage contact");
            cy.get(":nth-child(2) > :nth-child(4) > .govuk-button--warning").contains("Delete contact");
        });
    });

    describe("Named deputy contact", () => {
        it("should not display 'Named Deputy' or 'Monthly Spreadsheet Recipient' radio buttons for Pro deputies", () => {
            cy.visit("/supervision/deputies/3/contacts/add-contact");
            cy.get("#f-isNamedDeputy > .govuk-fieldset__legend").should("not.be.visible");
            cy.get("#f-isMonthlySpreadsheetRecipient > .govuk-fieldset__legend").should("not.be.visible");
        });
    });

    describe("Deleting a Contact", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/3/contacts");
        });

        it("shows content from list contacts tab and allow me to cancel", () => {
            cy.get(":nth-child(1) > :nth-child(4) > .govuk-button--warning").click();
            cy.get(".govuk-heading-l").contains("Delete contact");
            cy.get("#contact-form > .govuk-heading-m").contains("Do you want to remove");
            cy.get("#delete-contact-button-group > .govuk-link").contains("Cancel").click();
            cy.url().should("contain", "/supervision/deputies/3/contacts");
            cy.get(".moj-banner").should("not.exist");
        });

        it("shows content from manage contact page and allow me to delete", () => {
            cy.setCookie("success-route", "/contacts/1");
            cy.get(":nth-child(1) > :nth-child(4) > .govuk-button--secondary").click();
            cy.get(".govuk-button--warning").click();
            cy.get(".govuk-heading-l").contains("Delete contact");
            cy.get("#contact-form > .govuk-heading-m").contains("Do you want to remove");
            cy.get("#delete-contact-button-group > button").contains("Delete contact").click();
            cy.url().should(
                "contain",
                "/supervision/deputies/3/contacts?success=deletedContact&contactName=Minimal%20Contact",
            );
            cy.get(".moj-banner").should("contain", "'s details removed");
        });
    });
});
