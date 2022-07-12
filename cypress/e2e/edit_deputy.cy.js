describe("Edit deputy tab", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    it("should navigate to and from the Deputy tab", () => {
        cy.visit("/supervision/deputies/1");
        cy.contains(".govuk-button", "Manage team details").click();
        cy.url().should("include", "/supervision/deputies/1/manage-team-details");
        cy.get("#f-back-button").click();
        cy.get(".govuk-heading-l").contains("Deputy details");
    });

    it("the success banner shows on success", () => {
        cy.visit("/supervision/deputies/1/manage-team-details");

        cy.get("#f-team").focus().clear();
        cy.get("#f-team").type("New Team Name");
        cy.get("form").submit();
        cy.get(
            "body > div > main > div.moj-banner.moj-banner--success > div"
        ).should("contain", "Team details updated");
    });
});
