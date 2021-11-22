describe("Change ECM", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/public-authority/deputy/1/change-ecm");
    });

    it("has headers for different sections", () => {
        cy.get("h1").should("contain", "Change Executive Case Manager");
    })

    it("includes current ecm or leaves blank if none is set", () => {
        cy.get(".govuk-body").should("contain", "Current ECM:");
        cy.get(".govuk-label").should("contain", "Enter an Executive Case Manager name");
    });

    // it("has a drop down populated with members of the PA Deputy Team", () => {
    //     cy.get("#select-ecm").select('Cersei Lannister').should('have.value', '92');
    //     cy.get("#select-ecm").select('Jon Snow').should('have.value', '93');
    //     cy.get("#select-ecm").select('Eddard Stark').should('have.value', '94');
    // })

    it("directs me back to dashboard page if I press cancel", () => {
        cy.get(".govuk-link").should("contain", "Cancel").click();
        cy.url().should('not.include', '/change-ecm')
        cy.get("h1").should("contain", "Dashboard");
    })

    // it("allows me to fill in and submit the ecm form", () => {
    //     cy.get("#select-ecm").select('Jon Snow').should('have.value', '93');
    //     cy.get('form').submit()
    // })

    it("has a timeline event for when an ecm is changed", () => {
        cy.visit("/supervision/deputies/public-authority/deputy/1/timeline")
        cy.get(":nth-child(1) > .moj-timeline__header").should('contain', 'Executive Case Manager set to Public Authority deputy team');
        cy.get(":nth-child(1) > .moj-timeline__header > .moj-timeline__byline").should('contain', 'by Lay Team 1 - (Supervision')
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

});