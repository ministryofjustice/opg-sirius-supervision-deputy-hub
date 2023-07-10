describe("Add Contact", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/3/contacts/add-contact");
    });

    describe("Header", () => {
        it("shows content", () => {
            cy.get(".govuk-main-wrapper > header").contains("Create new contact");
            cy.get("#add-contact-form > :nth-child(2) > .govuk-label").contains("Name (required)");
            cy.get("#add-contact-form > :nth-child(3) > .govuk-label").contains("Job title");
            cy.get("#add-contact-form > :nth-child(4) > .govuk-label").contains("Email (required)");
            cy.get("#add-contact-form > :nth-child(5) > .govuk-label").contains("Phone (required)");
            cy.get("#add-contact-form > :nth-child(6) > .govuk-label").contains("Other phone");
            cy.get("#f-notes > .govuk-label").contains("Notes");
            cy.get("#f-isNamedDeputy > .govuk-fieldset__legend").contains("Named deputy (required)");
            cy.get("#f-isMainContact > .govuk-fieldset__legend").contains("Main contact (required)");
            cy.get(".govuk-button").contains("Save contact");
            cy.get(".govuk-button-group > .govuk-link").contains("Cancel");
        });
    });

    describe("Success submitting add contact form", () => {
        it("should allow me to submit the form", () => {
            cy.setCookie("success-route", "addContact");
            cy.get("#f-contactName").type("Test Contact");
            cy.get("#f-email").type("test@email.com");
            cy.get("#f-phoneNumber").type("0123456789");
            cy.get("#add-contact-form").submit();
            cy.url().should("contain", "/supervision/deputies/3/contacts?success=newContact");
            cy.get(".moj-banner").should("contain", "Contact added");
        });
    });

    describe("Error submitting empty add contact form", () => {
        it("shows error message when submitting empty data", () => {
        cy.setCookie("fail-route", "addContactEmpty");
        cy.get("#add-contact-form").submit();
        cy.get(".govuk-error-summary__body").should("contain", "Enter a name");
        cy.get(".govuk-error-summary__body").should("contain", "Enter a telephone number");
        cy.get(".govuk-error-summary__body").should("contain", "Enter an email address");
        cy.get(".govuk-error-summary__body").should("contain", "Select whether this contact is a main contact");
        cy.get(".govuk-error-summary__body").should("contain", "Select whether this contact is the named deputy");
        });
    });

    describe("Error submitting invalid add contact form", () => {
        it("shows error message when submitting invalid data", () => {
        cy.setCookie("fail-route", "addContactInvalid");
        cy.get("#add-contact-form").submit();
        cy.get(".govuk-error-summary__body").should("contain", "The name must be 255 characters or fewer");
        cy.get(".govuk-error-summary__body").should("contain", "The job title must be 255 characters or fewer");
        cy.get(".govuk-error-summary__body").should("contain", "Enter an email address in the correct format, like name@example.com");
        cy.get(".govuk-error-summary__body").should("contain", "The telephone number must be 255 characters or fewer");
        cy.get(".govuk-error-summary__body").should("contain", "The other telephone number must be 255 characters or fewer");
        cy.get(".govuk-error-summary__body").should("contain", "The note must be 255 characters or fewer");
        });
    });
});
