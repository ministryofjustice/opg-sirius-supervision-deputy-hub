describe("Change ECM", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/public-authority/deputy/1/change-ecm");
    });

    it("has headers for different sections", () => {
        cy.get("h1").should("contain", "Change Executive Case Manager");
    })

    it("leaves current ecm blank if none is set", () => {
        cy.get(".govuk-body").should("contain", "Current ECM:");
        cy.get(".govuk-label").should("contain", "Enter an Executive Case Manager name");
    });

    it("shows ecm if is set", () => {
        cy.visit("/supervision/deputies/public-authority/deputy/2/change-ecm");
        cy.get(".govuk-body").should("contain", "Current ECM:");
        cy.get(".govuk-body").should("contain", "LayTeam1 User2");
    });

    it("has a drop down populated with members of the PA Deputy Team", () => {
        cy.get("#select-ecm").type('S');
        cy.get("#select-ecm__listbox").find('li').should('have.length', 3);
        cy.get("#select-ecm").type('now');
        cy.get("#select-ecm__listbox").find('li').should('have.length', 1);
    })

    it("directs me back to dashboard page if I press cancel", () => {
        cy.get(".govuk-link").should("contain", "Cancel").click();
        cy.url().should('not.include', '/change-ecm')
        cy.get("h1").should("contain", "Dashboard");
    })

    it("allows me to fill in and submit the ecm form", () => {
        cy.visit("/supervision/deputies/public-authority/deputy/1/change-ecm");
        cy.setCookie("success-route", "ecm");
        cy.get("#select-ecm").type('S');
        cy.contains("#select-ecm__listbox", 'Jon Snow').click();
        cy.get('form').submit();
        cy.url().should("contain", "?success=ecm");
        cy.get("h1").should("contain", "Dashboard");
        cy.get(".moj-banner--success").should("contain", "Ecm changed to Public Authority Deputy Team" );
    })

    it("has a timeline event for when an ecm is automatically allocated on deputy creation", () => {
        cy.visit("/supervision/deputies/public-authority/deputy/1/timeline")
        cy.get(":nth-child(2) > .moj-timeline__header").should('contain', 'Executive Case Manager set to Public Authority deputy team');
        cy.get(":nth-child(2) > .moj-timeline__header > .moj-timeline__byline").should('contain', 'by Lay Team 1 - (Supervision')
    })

    it("has a timeline event for when an ecm is allocated", () => {
        cy.visit("/supervision/deputies/public-authority/deputy/1/timeline")
        cy.get(":nth-child(1) > .moj-timeline__header").should('contain', 'Executive Case Manager changed to PATeam1 User1');
        cy.get(":nth-child(1) > .moj-timeline__header > .moj-timeline__byline").should('contain', 'by case manager (12345678)')
    })

    it("displays warning when no ecm chosen and form submitted", () => {
        cy.setCookie("fail-route", "ecm")
        cy.get("#select-ecm").type('S');
        cy.get('form').submit();
        cy.get('.govuk-error-summary').should('contain', 'There is a problem');
        cy.get('.govuk-list > li > a').should('contain', 'Select an executive case manager');
    })

});

describe("Change ECM links to Dashboard", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/public-authority/deputy/1/");
    });

    it("has a link from the dashboard page", () => {
        cy.get(".moj-button-menu__wrapper > .govuk-button").should("contain", "Change ECM").click();
        cy.url().should('include', '/change-ecm');
        cy.get("h1").should("contain", "Change Executive Case Manager");
    })

    it("displays a warning if ECM is not set which links to the Change ECM page", () => {
        cy.get(".govuk-list > li > a").should("contain", "Assign an executive case manager").click();
        cy.url().should('include', '/change-ecm');
        cy.get("h1").should("contain", "Change Executive Case Manager");
    })

    it("does not display warning when ecm set", () => {
        cy.visit("/supervision/deputies/public-authority/deputy/2/");
        cy.get(".govuk-\\!-margin-bottom-2").should("contain", "LayTeam1 User2")
        cy.get(".govuk-list > li > a").should("not.exist");
    })

});