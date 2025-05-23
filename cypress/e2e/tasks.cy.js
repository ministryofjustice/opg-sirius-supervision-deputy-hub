describe("Tasks", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe("Navigation", () => {
        it("should navigate to and from the Tasks pages", () => {
            cy.visit("/supervision/deputies/1");
            cy.get(".moj-sub-navigation__list").contains("Tasks").click();

            cy.url().should("include", "/supervision/deputies/1/tasks");
            cy.contains(".govuk-heading-l", "Deputy tasks");

            cy.contains(".govuk-button", "Add a new task").click();
            cy.url().should("include", "/supervision/deputies/1/tasks/add-task");

            cy.contains(".govuk-link", "Cancel").click();
            cy.url().should("include", "/supervision/deputies/1/tasks");

            cy.contains(".govuk-button", "Manage task").first().click();
            cy.url().should("include", "/supervision/deputies/1/tasks/190");

            cy.contains(".govuk-link", "Cancel").click();
            cy.url().should("include", "/supervision/deputies/1/tasks");

            cy.contains(".govuk-button", "Mark as complete").first().click();
            cy.url().should("include", "/supervision/deputies/1/tasks/complete/190");

            cy.contains(".govuk-link", "Cancel").click();
            cy.url().should("include", "/supervision/deputies/1/tasks");
        });
    });

    describe("Showing Tasks", () => {
        it("should have required information on task page", () => {
            cy.visit("/supervision/deputies/1");
            cy.get(".moj-sub-navigation__list").contains("Tasks").click();

            cy.url().should("include", "/supervision/deputies/1/tasks");
            cy.contains(".govuk-heading-l", "Deputy tasks");

            cy.get(":nth-child(1) > .task_type").should("contain.text", "Assurance visit follow up");
            cy.get(":nth-child(1) > .assigned_to").should("contain.text", "Spongebob Squarepants");
            cy.get(":nth-child(1) > .due_date").should("contain.text", "29/01/2026");
            cy.get(":nth-child(1) > .task_type").should("contain.text", "Notes");

            cy.get(":nth-child(4) > .task_type").should("not.contain.text", "Notes");

            cy.get(":nth-child(2) > .task_type > .govuk-details").click();
            cy.get(":nth-child(2) > .task_type > .govuk-details > .govuk-details__text").should("be.visible");
            cy.get(":nth-child(2) > .task_type > .govuk-details > .govuk-details__text").should(
                "contain.text",
                "Notes about the enquiry",
            );
        });
    });

    describe("Adding a Task", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/1/tasks/add-task");
        });

        it("should add a task successfully", () => {
            cy.get('label:contains("Assurance visit review")').click();
            cy.get('label:contains("Due date")').type(new Date().toISOString().split("T")[0]);

            cy.get('.govuk-radios:has(label:contains("PA Team 1 - (Supervision) (Executive Case Manager)"))')
                .find("input")
                .should("be.checked");
            cy.get("#f-select-ecm").should("be.hidden");
            cy.get('label:contains("Someone else")').click();
            cy.get("#f-select-ecm").should("be.visible");
            cy.get("#f-select-ecm").type("S");
            cy.get("#f-select-ecm__listbox").find("li").should("have.length", 3);
            cy.get("#f-select-ecm").type("ta");
            cy.get("#f-select-ecm__listbox").find("li").should("have.length", 1);
            cy.contains("#f-select-ecm__listbox", "Eddard Stark").click();

            cy.get('label:contains("Notes")').type("Test note for task");
            cy.contains("button", "Save task").click();

            cy.url().should("contain", "/supervision/deputies/1/tasks");
            cy.contains(".moj-banner", "Assurance visit review task added");
        });

        it("displays validation errors", () => {
            cy.setCookie("fail-route", "addTask");
            cy.contains("button", "Save task").click();

            cy.get(".govuk-error-summary__title").should("contain", "There is a problem");
            cy.get(".govuk-error-summary__list").within(() => {
                cy.contains("li", "Select the task type");
                cy.contains("li", "This must be a real date");
            });

            cy.get("#add-task-form > :nth-child(2) > :nth-child(1).govuk-form-group--error").should("exist");
            cy.get("#add-task-form > :nth-child(2) > :nth-child(1)  #name-error-isEmpty").should(
                "contain",
                "Select the task type",
            );

            cy.get("#add-task-form > :nth-child(2) > :nth-child(2).govuk-form-group--error").should("exist");
            cy.get("#add-task-form > :nth-child(2) > :nth-child(2) #name-error-dateFalseFormat").should(
                "contain",
                "This must be a real date",
            );
            cy.get("#f-dueDate.govuk-input--error").should("exist");
        });
    });

    describe("Task timeline", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/1/timeline");
        });
        it("displays a task timeline event", () => {
            cy.get('[data-cy="task-created-event"]').within(() => {
                cy.contains(".moj-timeline__title", "Assurance visit follow up task created");
                cy.contains(".moj-timeline__byline", "by Lay Team 1 - (Supervision) (0123456789)");

                cy.get(".moj-timeline__description")
                    .get("li")
                    .should("contain", "Assigned to PA Team Workflow")
                    .next()
                    .should("contain", "Due date 13/07/2023")
                    .next()
                    .should("contain", "This is a note");
            });
        });
        it("displays an edited task timeline event", () => {
            cy.get('[data-cy="task-updated-event"]').within(() => {
                cy.contains(".moj-timeline__title", "Assurance visit follow up task updated");
                cy.contains(".moj-timeline__byline", "by Lay Team 1 - (Supervision) (0123456789)");

                cy.get(".moj-timeline__description")
                    .get("li")
                    .should("contain", "Due date 12/08/2023")
                    .next()
                    .should("contain", "editing and updating task notes");
            });
        });
    });

    describe("Task note", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/3/notes");
        });
        it("displays the task note title in the Notes tab correctly", () => {
            cy.contains(".moj-timeline__title", "General enquiry task created");
        });
        it("displays edited task notes in the Notes tab correctly", () => {
            cy.contains(".moj-timeline__title", "General enquiry task updated");
        });
    });

    describe("Editing Tasks", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/1/tasks");
            cy.contains(".govuk-button", "Manage task").first().click();
        });
        it("prefills all editable fields and allows them to be changed", () => {
            cy.setCookie("success-route", "/tasks/190");

            cy.url().should("include", "/supervision/deputies/1/tasks/190");
            cy.get(".govuk-heading-l").should("include.text", "Manage");

            cy.get("#assignedto-current-assignee").should("be.checked");
            cy.get("#assignedto-ecm").should("not.be.checked");
            cy.get("#assignedto-other").should("not.be.checked");
            cy.get("#f-2-note").should("contain.text", "Notes about the task");

            let nextYear = new Date().getFullYear() + 1;
            cy.get("#duedate").type(nextYear + "-01-01");
            cy.get("#assignedto-other").check();
            cy.get("#f-select-ecm").type("S");
            cy.get("#f-select-ecm__listbox").find("li").should("have.length", 3);
            cy.get("#f-select-ecm").type("ta");
            cy.get("#f-select-ecm__listbox").find("li").should("have.length", 1);
            cy.contains("#f-select-ecm__listbox", "Eddard Stark").click();
            cy.get("#f-2-note").clear();
            cy.get("#f-2-note").type("updated notes");

            cy.contains("button", "Save task").click();
            cy.get(".moj-banner--success").should("contain.text", "Assurance visit follow up task updated");
        });

        it("shows validation errors", () => {
            cy.setCookie("fail-route", "manageTask");
            cy.get("#f-2-note").clear();
            cy.get("#f-2-note").type("updated notes");
            cy.contains("button", "Save task").click();

            cy.get(".govuk-error-summary__title").should("contain", "There is a problem");
            cy.get(".govuk-error-summary__list").within(() => {
                cy.contains("li", "Enter a name of someone who works on the Public Authority team");
            });
        });
    });

    describe("Complete a task", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/1/tasks");
            cy.contains(".govuk-button", "Mark as complete").first().click();
        });
        it("displays the task details for task about to be completed and allows notes to be added", () => {
            cy.setCookie("success-route", "/tasks/190");
            cy.get(".govuk-heading-l").should("contain.text", "Complete Task");
            cy.get(":nth-child(1) > .govuk-table__cell").should("contain.text", "Assurance visit follow up");
            cy.get(":nth-child(2) > .govuk-table__cell").should("contain.text", "Notes about the task");
            cy.get(":nth-child(3) > .govuk-table__cell").should("contain.text", "29/01/2026");
            cy.get(":nth-child(4) > .govuk-table__cell").should("contain.text", "Spongebob Squarepants");
            cy.get("#f-notes").type("Notes for the event about to be completed");
            cy.get(".govuk-button").contains("Complete task").click();
            cy.get(".moj-banner--success").should("be.visible");
            cy.get(".moj-banner--success").should("contain.text", "Assurance visit follow up task completed");
        });
        it("shows validation message if notes over 1000 characters", () => {
            cy.setCookie("fail-route", "completeTask");
            cy.get("#f-notes").type("Notes for the event about to be completed");
            cy.get(".govuk-button").contains("Complete task").click();
            cy.get(".govuk-error-summary").should("be.visible");
            cy.get(".govuk-error-summary__title").should("contain", "There is a problem");
            cy.get(".govuk-error-summary__body").should("contain", "The note must be 1000 characters or fewer");
            cy.get(".govuk-form-group--error").should("exist");
            cy.get("#f-notes.govuk-input--error").should("exist");
            cy.get("#name-error-stringLengthTooLong").should("contain", "The note must be 1000 characters or fewer");
        });
    });
});
