describe("Timeline", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    it("should navigate to and from the Timeline tab", () => {
        cy.visit("/supervision/deputies/1");
        cy.get(".moj-sub-navigation__list").contains("Timeline").click();

        cy.url().should("include", "/supervision/deputies/1/timeline");
        cy.get(".main > header").contains("Timeline");
    })

    it("contains appropriate test data for a timeline event", () => {
        cy.visit("/supervision/deputies/1/timeline");

        cy.get(".moj-timeline__title").should(
            "contain",
            "New client added to deputyship"
        );
        cy.get(".moj-timeline__byline").should(
            "contain",
            "by system admin (12345678)"
        );
        cy.get("time").should("contain", "09/09/2021 14:01:59");
        cy.get(".govuk-list > :nth-child(1)").should(
            "contain",
            "Order number: 03305972"
        );
        cy.get(".govuk-list > :nth-child(2)").should(
            "contain",
            "Sirius ID: 7000-0000-1995"
        );
        cy.get(".govuk-list > :nth-child(3)").should(
            "contain",
            "Order type: pfa"
        );
        cy.get(".govuk-list > :nth-child(4)").should(
            "contain",
            "Client: Duke John Fearless"
        );
    });
});
