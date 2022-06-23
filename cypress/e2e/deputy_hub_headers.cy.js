describe("Pro Deputy Hub", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/3");
    });

    describe("Header", () => {
        it("shows opg sirius within banner", () => {
            cy.contains(".moj-header__link", "OPG");
            cy.contains(".moj-header__link", "Sirius");
        });

        const expected = ["Workflow", "Supervision", "LPA", "Admin", "Logout"];

        it("has working nav links within header banner", () => {
            cy.get(".moj-header__navigation-list")
                .children()
                .each(($el, index) => {
                    cy.wrap($el).should("contain", expected[index]);
                    let $linkName = expected[index].toLowerCase();
                    cy.wrap($el)
                        .find("a")
                        .should("have.attr", "href")
                        .and("contain", `/${$linkName}`);
                });
        });
    });

    describe("Footer", () => {
        it("the footer should contain a link to the open government licence", () => {
            cy.get(
                ".govuk-footer__licence-description > .govuk-footer__link"
            ).should(
                "have.attr",
                "href",
                "https://www.nationalarchives.gov.uk/doc/open-government-licence/version/3/"
            );
        });

        it("the nav link should contain the crown copyright logo", () => {
            cy.get(".govuk-footer__copyright-logo").should(
                "have.attr",
                "href",
                "https://www.nationalarchives.gov.uk/information-management/re-using-public-sector-information/uk-government-licensing-framework/crown-copyright/"
            );
        });
    });

    describe("Pro deputy details", () => {
        it("the page should contain the deputy name and status", () => {
            cy.get(".hook_header_deputy_name").should(
                "contain",
                "firstname surname"
            );
            cy.get(".moj-badge").should("exist");
        });

        it("the page should contain the firm", () => {
            cy.get(".hook_header_firm_name > .govuk-link").should(
                "contain",
                "This is the Firm Name"
            );
        });

        it("the page should contain the deputy number", () => {
            cy.get(".hook_header_deputy_number").should(
                "contain",
                "Deputy Number: 1000"
            );
        });

        it("the page should contain the executive case manager", () => {
            cy.get(".hook_header_ecm").should(
                "contain",
                "Executive Case Manager: displayName"
            );
        });
    });
});
