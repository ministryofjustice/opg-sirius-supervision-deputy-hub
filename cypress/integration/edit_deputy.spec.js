describe("Clients tab", () => {
    beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/supervision/deputies/public-authority/deputy/1/manage-team-details");
    });

    it("the success banner shows on success", () => {
        cy.get('#f-team').focus().clear();
        cy.get('#f-team').type("New Team Name")
        cy.get('form').submit()
        cy.get('.moj-banner__message').contains("Success You have successfully edited your details.")
    })
});
