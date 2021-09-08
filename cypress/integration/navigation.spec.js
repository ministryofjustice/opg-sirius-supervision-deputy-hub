describe("Navigation bar", () => {
  beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/supervision/deputies/public-authority/deputy/1/");
  });

  it("has working nav links for different tabs", () => {
      cy.get(".moj-sub-navigation__list > :nth-child(1) > a").should("contain", "Dashboard");
      cy.get(".moj-sub-navigation__list > :nth-child(1) > a").should("have.attr", "href", "/supervision/deputies/public-authority/deputy/1/");
  })

});