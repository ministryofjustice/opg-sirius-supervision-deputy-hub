describe("Deputy details tab", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/1");
    });

    it("has headers for different sections", () => {
        cy.get("h1").should("contain", "Deputy details");
        cy.get("h2").should("contain", "Team details");
    });

    it("has a manage important information button", () => {
        cy.get(":nth-child(2) > .govuk-button")
            .should("contain", "Manage important information")
            .click();
        cy.url().should("include", "manage-important-information");
    });

    it("has a button which can take you to change ecm", () => {
        cy.get("#change-ecm").should("contain", "Change ECM").click();
        cy.url().should("include", "change-ecm");
    });

    it("lists active cases", () => {
        cy.get("#overview").should("contain", "3");
        cy.get("#overview").should("contain", "Active cases");
    });

    describe("Deputy contact details", () => {
        it("has rows in tables with accurate keys and values", () => {
            cy.get(
                "#team-details > :nth-child(1) > .govuk-summary-list__key"
            ).should("contain", "Deputy name");
            cy.get(
                "#team-details > :nth-child(1) > .govuk-summary-list__value"
            ).should("contain", "Test Organisation");
            cy.get(
                "#team-details > :nth-child(2) > .govuk-summary-list__key"
            ).should("contain", "Telephone");
            cy.get(
                "#team-details > :nth-child(2) > .govuk-summary-list__value"
            ).should("contain", "0115 876 5574");
            cy.get(
                "#team-details > :nth-child(3) > .govuk-summary-list__key"
            ).should("contain", "Email");
            cy.get(
                "#team-details > :nth-child(3) > .govuk-summary-list__value"
            ).should("contain", "deputyship@essexcounty.gov.uk");
            cy.get(
                "#team-details > :nth-child(4) > .govuk-summary-list__key"
            ).should("contain", "Postal address");
            cy.get(
                "#team-details > :nth-child(4) > .govuk-summary-list__value"
            ).should("contain", "Deputyship Team");
        });

        it("has a href link for email addresses", () => {
            cy.get(".govuk-summary-list__value > a").should(
                "have.attr",
                "href"
            );
        });

        it("displays warning when no ecm set", () => {
            cy.get(".govuk-list > li").should(
                "contain",
                "An executive case manager has not been assigned. Assign an executive case manager"
            );
        });

        it("has a button which can take you to manage team details", () => {
            cy.get(".govuk-main-wrapper > :nth-child(7) > :nth-child(1) > .govuk-button")
                .should("contain", "Manage team details")
                .click();
            cy.url().should("include", "manage-team-details");
        });
    });

    describe("Important information", () => {
        it("has rows in tables with accurate keys and values", () => {
            cy.get(
                ":nth-child(2) > .govuk-summary-list > :nth-child(1) > .govuk-summary-list__key"
            ).should("contain", "Monthly spreadsheet");
            cy.get(
                ":nth-child(2) > .govuk-summary-list > :nth-child(1) > .govuk-summary-list__value"
            ).should("contain", "No");
            cy.get(
                ":nth-child(2) > .govuk-summary-list > :nth-child(2) > .govuk-summary-list__key"
            ).should("contain", "Independent visitor charges");
            cy.get(
                ":nth-child(2) > .govuk-summary-list > :nth-child(2) > .govuk-summary-list__value"
            ).should("contain", "Unknown");
            cy.get(
                ":nth-child(2) > .govuk-summary-list > :nth-child(3) > .govuk-summary-list__key"
            ).should("contain", "Bank charges");
            cy.get(
                ":nth-child(2) > .govuk-summary-list > :nth-child(3) > .govuk-summary-list__value"
            ).should("contain", "Yes");
            cy.get(
                ":nth-child(2) > .govuk-summary-list > :nth-child(4) > .govuk-summary-list__key"
            ).should("contain", "APAD");
            cy.get(
                ":nth-child(2) > .govuk-summary-list > :nth-child(4) > .govuk-summary-list__value"
            ).should("contain", "Yes");
            cy.get(
                ":nth-child(2) > .govuk-summary-list > :nth-child(5) > .govuk-summary-list__key"
            ).should("contain", "Report system");
            cy.get(
                ":nth-child(2) > .govuk-summary-list > :nth-child(5) > .govuk-summary-list__value"
            ).should("contain", "CASHFAC");
            cy.get(
                ":nth-child(2) > .govuk-summary-list > :nth-child(6) > .govuk-summary-list__key"
            ).should("contain", "Annual billing preference");
            cy.get(
                ":nth-child(2) > .govuk-summary-list > :nth-child(6) > .govuk-summary-list__value"
            ).should("contain", "Schedule and Invoice");
            cy.get(
                ":nth-child(2) > .govuk-summary-list > :nth-child(7) > .govuk-summary-list__key"
            ).should("contain", "Other important information");
            cy.get(
                ":nth-child(2) > .govuk-summary-list > :nth-child(7) > .govuk-summary-list__value"
            ).should("contain", "some info for the pa deputy");
        });
    });
});
