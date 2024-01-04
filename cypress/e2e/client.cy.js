describe("Clients tab", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    it("should navigate to the Clients tab", () => {
        cy.visit("/supervision/deputies/1");
        cy.get(".moj-sub-navigation__list").contains("Clients").click();
        cy.url().should("include", "/supervision/deputies/1/clients");
        cy.get("h1").should("contain", "Clients");
    });

    describe("clients table", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/1/clients?sort=surname:asc");
        });

        it("displays 7 column headings", () => {
            cy.get(".govuk-table__row").find("th").should("have.length", 7);

            const expected = [
                "Client",
                "Accommodation type",
                "Status",
                "Supervision level",
                "Visits",
                "Report due",
                "Risk",
            ];

            cy.get(".govuk-table__head > .govuk-table__row")
                .children()
                .each(($el, index) => {
                    cy.wrap($el).should("contain", expected[index]);
                });
        });

        it("lists clients with active/closed/duplicate orders", () => {
            cy.get(".govuk-table__body > .govuk-table__row").should(
                "have.length",
                3
            );
        });

        it("Clients surname have been sorted in order of ascending by default", () => {
            cy.get(":nth-child(1) > .client_name_ref > .govuk-link").should(
                "contain",
                "Burgundy"
            );
            cy.get(":nth-child(2) > .client_name_ref > .govuk-link").should(
                "contain",
                "Dauphin"
            );
            cy.get(":nth-child(3) > .client_name_ref > .govuk-link").should(
                "contain",
                "Here"
            );
        });

        it("Clients surname have been sorted in order of descending", () => {
            cy.get('[aria-sort="ascending"] > a > button').click();
            cy.url().should("contain", "order-by=surname&sort=desc");
            cy.get('[aria-sort="descending"] > a > button').click();
            cy.url().should("contain", "order-by=surname&sort=asc");
        });

        it("displays REM warning label", () => {
            cy.get(":nth-child(1) > .rem-warning").should("contain", "REM warning");
            cy.get(":nth-child(2) > .rem-warning").should("not.exist");
        });
    });
});
