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

        it("lists active cases", () => {
            cy.get("#overview").should("contain", "3");
            cy.get("#overview").should("contain", "Active cases");

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
            cy.get('[label="sort-name-column-ascending"] > button').click();
            cy.get(":nth-child(1) > .client_name_ref > .govuk-link").should(
                "contain",
                "Here"
            );
            cy.get(":nth-child(2) > .client_name_ref > .govuk-link").should(
                "contain",
                "Dauphin"
            );
            cy.get(":nth-child(3) > .client_name_ref > .govuk-link").should(
                "contain",
                "Burgundy"
            );
        });

        it("Clients report due dates have been sorted in order of ascending", () => {
            cy.get('[label="sort-due-date-column-none"] > button').click();
            cy.get(":nth-child(1) > .reports").should("contain", "21/12/2015");
            cy.get(":nth-child(2) > .reports").should("contain", "01/10/2018");
            cy.get(":nth-child(3) > .reports").should("contain", "-");
        });

        it("Clients report due dates have been sorted in order of descending", () => {
            cy.get('[label="sort-due-date-column-none"] > button').click();
            cy.get('[label="sort-due-date-column-ascending"] > button').click();
            cy.get(":nth-child(1) > .reports").should("contain", "-");
            cy.get(":nth-child(2) > .reports").should("contain", "01/10/2018");
            cy.get(":nth-child(3) > .reports").should("contain", "21/12/2015");
        });

        it("Clients risk have been sorted in order of ascending", () => {
            cy.get('[label="sort-aria-column-none"] > button').click();
            cy.get(":nth-child(1) > .data-crec").should("contain", "2");
            cy.get(":nth-child(2) > .data-crec").should("contain", "3");
            cy.get(":nth-child(3) > .data-crec").should("contain", "4");
        });

        it("Clients risk have been sorted in order of descending", () => {
            cy.get('[label="sort-aria-column-none"] > button').click();
            cy.get('[label="sort-aria-column-ascending"] > button').click();
            cy.get(":nth-child(1) > .data-crec").should("contain", "4");
            cy.get(":nth-child(2) > .data-crec").should("contain", "3");
            cy.get(":nth-child(3) > .data-crec").should("contain", "2");
        });

        it("Clients have a latest visit recorded", () => {
            cy.get(":nth-child(1) > .visit_type").should("contain", "01/01/2000");
            cy.get(":nth-child(1) > .visit_type").should(
                "contain",
                "Standard visit"
            );
            cy.get(":nth-child(1) > .visit_type").should("contain", "Green");
            cy.get(":nth-child(2) > .visit_type").should("contain", "03/03/2020");
            cy.get(":nth-child(2) > .visit_type").should("contain", "Urgent visit");
            cy.get(":nth-child(2) > .visit_type").should("contain", "Red");
            cy.get(":nth-child(3) > .visit_type").should("contain", "-");
        });
    });
});
