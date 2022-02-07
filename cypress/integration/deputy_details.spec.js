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

    it("has rows in tables with accurate keys and values", () => {
        cy.get("#team-details > :nth-child(1) > .govuk-summary-list__key").should("contain", "Deputy name");
        cy.get("#team-details > :nth-child(1) > .govuk-summary-list__value").should("contain", "Test Organisation");
        cy.get("#team-details > :nth-child(2) > .govuk-summary-list__key").should("contain", "Telephone");
        cy.get("#team-details > :nth-child(2) > .govuk-summary-list__value").should("contain", "0115 876 5574");
        cy.get("#team-details > :nth-child(3) > .govuk-summary-list__key").should("contain", "Email");
        cy.get("#team-details > :nth-child(3) > .govuk-summary-list__value").should("contain", "deputyship@essexcounty.gov.uk");
        cy.get("#team-details > :nth-child(4) > .govuk-summary-list__key").should("contain", "Postal address");
        cy.get("#team-details > :nth-child(4) > .govuk-summary-list__value").should("contain", "Deputyship Team");
    });

    it("has a href link for email addresses", () => {
        cy.get(".govuk-summary-list__value > a").should("have.attr", "href");
    });

    it("displays warning when no ecm set", () => {
        cy.get(".govuk-list > li").should(
            "contain",
            "An executive case manager has not been assigned. Assign an executive case manager"
        );
    });
});
