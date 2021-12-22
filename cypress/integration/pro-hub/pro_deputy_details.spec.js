describe("Pro Deputy Hub", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/public-authority/deputy/3");
    });

    describe("Deputy details", () => {
        it("shows all the deputy details", () => {
            cy.get(".hook_deputy_name").contains("firstname surname");
            cy.get(".hook_deputy_phone_number").contains("1111111");
            cy.get(".hook_deputy_email").contains("email@something.com");
        });
    });
});
