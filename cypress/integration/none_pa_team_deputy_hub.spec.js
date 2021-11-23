describe("Deputy Hub", () => {
  beforeEach(() => {
      Cypress.on('uncaught:exception', (err, runnable) => {
          if (err.message.includes('selectElement is not defined')){return false}
      })
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/supervision/deputies/public-authority/deputy/2/");
  });

    it("the page should contain the deputy number", () => {
        cy.get(".govuk-caption-m").eq(0).should("contain",  "Deputy Number: 23")
    })

    it("the page should not contain the warning error", () => {
    cy.get('.moj-banner__message > a').should('not.exist')
    })
});
