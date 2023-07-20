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

    describe('Showing Tasks', () => {
        it("should have required information on task page", () => {
            cy.visit("/supervision/deputies/1");
            cy.get(".moj-sub-navigation__list").contains("Tasks").click();

            cy.url().should("include", "/supervision/deputies/1/tasks");
            cy.contains(".govuk-heading-l", "Deputy tasks");

            cy.get(':nth-child(1) > .task_type').should('contain.text', 'Assurance visit follow up');
            cy.get(':nth-child(1) > .assigned_to').should('contain.text', 'Spongebob Squarepants');
            cy.get(':nth-child(1) > .due_date').should('contain.text', '29/01/2021');
            cy.get(':nth-child(1) > .task_type').should('not.contain.text', 'Notes');

            cy.get(':nth-child(2) > .task_type').should('contain.text', 'Notes');
            cy.get(':nth-child(2)  > .task_type > .govuk-details').click();
            cy.get(':nth-child(2) > .task_type > .govuk-details > .govuk-details__text').should('be.visible');
            cy.get(':nth-child(2) > .task_type > .govuk-details > .govuk-details__text').should('contain.text', 'Notes about the enquiry');
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

    describe("Task timeline", () => {
        it("displays a task timeline event", () => {
            cy.visit("/supervision/deputies/1/timeline");

            cy.get('[data-cy="task-created-event"]').within(() => {
                cy.contains(".moj-timeline__title", "Assurance visit follow up task created");
                cy.contains(".moj-timeline__byline", "by Lay Team 1 - (Supervision) (0123456789)");

                cy.get(".moj-timeline__description").get("li")
                    .should("contain", "Assigned to PA Team Workflow")
                    .next()
                    .should("contain", "Due date 13/07/2023")
                    .next()
                    .should("contain", "This is a note");
            });
        });
    })

    describe("Task note", () => {
        it("displays the task note title in the Notes tab correctly", () => {
            cy.visit("/supervision/deputies/3/notes");
            cy.contains(".moj-timeline__title", "General enquiry task created");
        });
    })

});
