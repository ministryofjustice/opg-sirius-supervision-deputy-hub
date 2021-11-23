describe("Edit deputy tab", () => {
    beforeEach(() => {
        Cypress.on('uncaught:exception', (err, runnable) => {
            if (err.message.includes('selectElement is not defined')){return false}
        })
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/public-authority/deputy/1/manage-team-details");
    });

    it("the success banner shows on success", () => {
        cy.get('#f-team').focus().clear();
        cy.get('#f-team').type("New Team Name")
        cy.get('form').submit()
        cy.get("body > div > main > div.moj-banner.moj-banner--success > div").should("contain", "Team details updated");
    })
});
