describe("Deputy Hub", () => {
  beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/supervision/deputies/public-authority/deputy/2/");
  });

  it("the page should not contain the warning error", () => {
    cy.get('.moj-banner__message > a').should('not.exist')
  })
});
