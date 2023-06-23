describe("Tasks", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe('Navigation', () => {
       it("should navigate to and from the Tasks pages", () => {
           cy.visit("/supervision/deputies/1");
           cy.get(".moj-sub-navigation__list").contains("Tasks").click();

           cy.url().should("include", "/supervision/deputies/1/tasks");
           cy.contains(".govuk-heading-l", "Deputy tasks");

           cy.contains(".govuk-button", "Add a new task").click();
           cy.url().should("include","/supervision/deputies/1/tasks/add-task");

           cy.contains(".govuk-link", "Cancel").click();
           cy.url().should("include", "/supervision/deputies/1/tasks");
       });
    });

    describe("Adding a Task", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/1/tasks/add-task");
        });
    });
});
