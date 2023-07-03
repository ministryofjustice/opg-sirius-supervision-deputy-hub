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

        it("should add a task successfully", () => {
            cy.get('label:contains("Assurance visit review")').click();
            cy.get('label:contains("Due date (required")').type(new Date().toISOString().split("T")[0]);

            cy.get('.govuk-radios:has(label:contains("PA Team 1 - (Supervision) (Executive Case Manager)"))')
                .find('input')
                .should("be.checked");
            cy.get('#select-ecm').should("be.hidden");
            cy.get('label:contains("Someone else")').click();
            cy.get('#select-ecm').should("be.visible");
            cy.get("#select-ecm").type("S");
            cy.get("#select-ecm__listbox").find("li").should("have.length", 2);
            cy.get("#select-ecm").type("t");
            cy.get("#select-ecm__listbox").find("li").should("have.length", 1);
            cy.contains("#select-ecm__listbox", "Steven Toast").click();

            cy.get('label:contains("Notes")').type("Test note for task");
            cy.contains("button", "Save task").click();

            cy.url().should("contain", "/supervision/deputies/1/tasks");
            cy.contains(".moj-banner", "Assurance visit review task added");
        });

        it("displays validation errors", () => {
            cy.setCookie("fail-route", "addTask");
            cy.contains("button", "Save task").click();

            cy.get(".govuk-error-summary__title").should(
                "contain",
                "There is a problem"
            );
            cy.get(".govuk-error-summary__list").within(() => {
                cy.contains("li", "Select the task type");
            });
        });
    });
});
