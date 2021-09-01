describe("Deputy Hub", () => {
  beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/supervision/deputies/public-authority/");
  });

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

  it("has working nav links within banner", () => {
    cy.get(".moj-header__navigation-list")
    .children()
    .each(($el, index) => {
        cy.wrap($el).should("contain", expected[index]);
    });
  })

  it("the nav link should contain supervision", () => {
    cy.get(".moj-header__navigation-list > :nth-child(1) > a").should("have.attr", "href", "http://localhost:8080/supervision")
  })  
  
  it("the nav link should contain lpa", () => {
    cy.get(".moj-header__navigation-list > :nth-child(2) > a").should("have.attr", "href", "http://localhost:8080/lpa")
  })
  
  it("the nav link should contain admin", () => {
    cy.get(".moj-header__navigation-list > :nth-child(3) > a").should("have.attr", "href", "http://localhost:8080/admin")
  })
  
  it("the nav link should contain logout", () => {
    cy.get(".moj-header__navigation-list > :nth-child(4) > a").should("have.attr", "href", "http://localhost:8080/auth/logout")
  }) 


  it("the footer should contain a link to the GOV.UK Prototype Kit", () => {
    cy.get(".govuk-footer__inline-list > :nth-child(1) > a").should("have.attr", "href", "https://govuk-prototype-kit.herokuapp.com/")
  })  
  
  it("the footer should contain a link to clear data", () => {
    cy.get(".govuk-footer__inline-list > :nth-child(2) > a").should("have.attr", "href", "/prototype-admin/clear-data")
  })
  
  it("the footer should contain a link to the open government licence", () => {
    cy.get(".govuk-footer__licence-description > .govuk-footer__link").should("have.attr", "href", "https://www.nationalarchives.gov.uk/doc/open-government-licence/version/3/")
  })
  
  it("the nav link should contain the crown copyright logo", () => {
    cy.get(".govuk-footer__copyright-logo").should("have.attr", "href", "https://www.nationalarchives.gov.uk/information-management/re-using-public-sector-information/uk-government-licensing-framework/crown-copyright/")
  }) 
});