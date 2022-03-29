describe("Pro Deputy Hub", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/3");
    });

    it("has a button which can take you to change firm", () => {
        cy.get(".moj-button-menu__wrapper > .govuk-button").should(
            "contain",
            "Change firm"
        ).click();
        cy.url().should("include", "change-firm");
    });

    describe("Deputy contact details", () => {
        it("shows all the deputy details", () => {
            cy.get(".hook_deputy_name").contains("firstname surname");
            cy.contains(".hook_deputy_firm_name", "This is the Firm Name")
                .find('a')
                .should("have.attr", "href")
                .and('contain', "/supervision/deputies/firm/0");
            cy.get(".hook_deputy_phone_number").contains("1111111");
            cy.get(".hook_deputy_email").contains("email@something.com");
        });

        it("has a button which can take you to manage deputy contact details", () => {
            cy.get("[data-cy=manage-deputy-contact-details-btn]").should(
                "contain",
                "Manage deputy contact details"
            ).click();
            cy.url().should("include", "manage-deputy-contact-details");
        });
    });

    describe("Important information", () => {

        it("has rows in tables with accurate keys and values", () => {
            cy.get(":nth-child(2) > .govuk-summary-list > :nth-child(1) > .govuk-summary-list__key").should("contain", "Complaints");
            cy.get(":nth-child(2) > .govuk-summary-list > :nth-child(1) > .govuk-summary-list__value").should("contain", "Yes");
            cy.get(":nth-child(2) > .govuk-summary-list > :nth-child(2) > .govuk-summary-list__key").should("contain", "Panel deputy");
            cy.get(":nth-child(2) > .govuk-summary-list > :nth-child(2) > .govuk-summary-list__value").should("contain", "Yes");
            cy.get(":nth-child(2) > .govuk-summary-list > :nth-child(3) > .govuk-summary-list__key").should("contain", "Annual billing preference");
            cy.get(":nth-child(2) > .govuk-summary-list > :nth-child(3) > .govuk-summary-list__value").should("contain", "Schedule");
            cy.get(":nth-child(2) > .govuk-summary-list > :nth-child(4) > .govuk-summary-list__key").should("contain", "Other important information");
            cy.get(":nth-child(2) > .govuk-summary-list > :nth-child(4) > .govuk-summary-list__value").should("contain", "Some important information is here");
        });
    });
});

