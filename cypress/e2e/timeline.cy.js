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
    });

    it("contains appropriate test data for a timeline event", () => {
        cy.visit("/supervision/deputies/1/timeline");

        cy.get('[data-cy="new-client-added-event"]').within(() => {
            cy.contains(".moj-timeline__title", "New client added to deputyship");
            cy.contains(".moj-timeline__byline", "by system admin (12345678)");

            cy.get("time").should("contain", "09/09/2021");

            cy.get(".moj-timeline__description")
                .get("li")
                .should("contain", "Order number: 03305972")
                .next()
                .should("contain", "Sirius ID: 7000-0000-1995")
                .next()
                .should("contain", "Order type: pfa")
                .next()
                .should("contain", "Client: Duke John Fearless");
        });
    });

    it("displays deputy contact timeline events", () => {
        cy.visit("/supervision/deputies/1/timeline");

        let events = [
            {
                name: "contact-added-event",
                title: "Mr Deputy Contact added as a contact",
                description: "Name: Mr Deputy Contact",
            },
            {
                name: "contact-edited-event",
                title: "Mr Deputy Contact's details updated",
                description: "Name: Mr Deputy Contact",
            },
            {
                name: "contact-set-as-main-event",
                title: "Main contact set to Mr Deputy Contact",
            },
            {
                name: "contact-removed-as-main-event",
                title: "Mr Deputy Contact removed as a Main contact",
            },
            {
                name: "contact-set-as-named-event",
                title: "Named deputy set to Mr Deputy Contact",
            },
            {
                name: "contact-removed-as-named-event",
                title: "Mr Deputy Contact removed as the Named deputy",
            },
            {
                name: "task-reassigned-event",
                title: "PDR report due task reassigned",
                description: "Assigned to Pro Team 2 - (Supervision)",
            },
        ];

        events.forEach((event) => {
            cy.get('[data-cy="' + event.name + '"]').within(() => {
                cy.get(".moj-timeline__title").should("contain.text", event.title);
                if ("description" in event) {
                    cy.get(".moj-timeline__description").contains(event.description);
                }
            });
        });
    });
});
