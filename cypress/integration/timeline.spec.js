describe("Timeline", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/1/timeline");
    });

    it("has a header called timeline", () => {
        cy.get(".main > header").should("contain", "Timeline");
    });

    it("lists timeline events in date ascending order", () => {
        const timelineItems = cy.get(".moj-timeline__item");

        timelineItems.first().within((item) => {
            cy.wrap(item).contains(
                ".moj-timeline__title",
                "Deputy firm updated"
            );
            cy.wrap(item).contains(
                ".moj-timeline__byline",
                "by case manager (12345678)"
            );
            cy.wrap(item).contains("time", "22/03/2022 15:57:11");
            cy.wrap(item).contains(
                ".govuk-list > :nth-child(1)",
                "New firm: another firm - 1000001"
            );
            cy.wrap(item).contains(
                ".govuk-list > :nth-child(2)",
                "Old firm: new firm - 1000000"
            );
        });

        timelineItems.next().within((item) => {
            cy.wrap(item).contains(
                ".moj-timeline__title",
                "Deputy firm updated"
            );
            cy.wrap(item).contains(
                ".moj-timeline__byline",
                "by case manager (12345678)"
            );
            cy.wrap(item).contains("time", "22/03/2022 15:56:53");
        });

        timelineItems.next().within((item) => {
            cy.wrap(item).contains(
                ".moj-timeline__title",
                "Executive Case Manager changed to PATeam1 User1"
            );
            cy.wrap(item).contains(
                ".moj-timeline__byline",
                "by case manager (12345678)"
            );
            cy.wrap(item).contains("time", "24/11/2021 14:01:59");
        });

        timelineItems.next().within((item) => {
            cy.wrap(item).contains(
                ".moj-timeline__title",
                "Executive Case Manager set to Public Authority deputy team"
            );
            cy.wrap(item).contains(
                ".moj-timeline__byline",
                "by Lay Team 1 - (Supervision) (0123456789)"
            );
            cy.wrap(item).contains("time", "10/10/2021 15:01:59");
        });
    });

    it("contains appropriate test data for a timeline event", () => {
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
