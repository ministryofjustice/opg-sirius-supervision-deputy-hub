describe("Timeline", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    it("should navigate to and from the Timeline tab", () => {
        cy.visit("/supervision/deputies/1");
        cy.get(".moj-sub-navigation__list").contains("Timeline").click();

        cy.url().should("include", "/supervision/deputies/1/timeline");
        cy.get(".govuk-heading-l").contains("Timeline");
    })

    it("contains appropriate test data for a timeline event", () => {
        cy.visit("/supervision/deputies/1/timeline");

        cy.get('[data-cy="new-client-added-event"]').within(() => {
            cy.contains(".moj-timeline__title", "New client added to deputyship");
            cy.contains(".moj-timeline__byline", "by system admin (12345678)");

            cy.get("time").should("contain", "09/09/2021");

            cy.get(".moj-timeline__description").get("li")
                .should("contain", "Order number: 03305972")
                .next()
                .should("contain", "Sirius ID: 7000-0000-1995")
                .next()
                .should("contain", "Order type: pfa")
                .next()
                .should("contain", "Client: Duke John Fearless");
        });
    });
});
