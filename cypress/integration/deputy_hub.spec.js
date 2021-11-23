Cypress.on('uncaught:exception', (err, runnable) => {
  // returning false here prevents Cypress from
  // failing the test
  return false
})

describe("Deputy Hub", () => {
  beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/supervision/deputies/public-authority/deputy/1/");
  });

  describe("Header", () => {
    it("shows opg sirius within banner", () => {
      cy.contains(".moj-header__link", "OPG");
      cy.contains(".moj-header__link", "Sirius");
    });

    const expected = [
      "Supervision",
      "LPA",
      "Admin",
      "Logout",
    ];

    it("has working nav links within header banner", () => {
      cy.get(".moj-header__navigation-list")
          .children()
          .each(($el, index) => {
            cy.wrap($el).should("contain", expected[index]);
            let $linkName = (expected[index].toLowerCase());
            cy.wrap($el).find('a').should("have.attr", "href").and("contain", `/${$linkName}`);

          });
    });
  });

  describe("Footer", () => {
    it("the footer should contain a link to the open government licence", () => {
      cy.get(".govuk-footer__licence-description > .govuk-footer__link").should("have.attr", "href", "https://www.nationalarchives.gov.uk/doc/open-government-licence/version/3/")
    })

    it("the nav link should contain the crown copyright logo", () => {
      cy.get(".govuk-footer__copyright-logo").should("have.attr", "href", "https://www.nationalarchives.gov.uk/information-management/re-using-public-sector-information/uk-government-licensing-framework/crown-copyright/")
    })
  });

  it("the page should contain the deputy organisation name", () => {
    cy.get(".govuk-heading-m").should("contain", "Test Organisation")
  })

  it("the page should contain the deputy number", () => {
    cy.get(".govuk-caption-m").eq(0).should("contain",  "Deputy Number: 11")
  })

  it("the page should contain the executive case manager label", () => {
    cy.get(".govuk-caption-m").eq(1).should("contain",  "Executive Case Manager:")
  })

  it("the page should contain the warning error", () => {
    cy.get('.govuk-list > li > a').should("contain",  "Assign an executive case manager")
  })
});
