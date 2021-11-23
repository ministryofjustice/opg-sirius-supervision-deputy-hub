
describe("Change ECM", () => {
    beforeEach(() => {
        Cypress.on('uncaught:exception', (err, runnable) => {
            if (err.message.includes('selectElement is not defined')){return false}
        })
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

    it("has a drop down populated with members of the PA Deputy Team", () => {
        cy.get("#select-ecm").should('contain', 'Jon Snow')
        cy.get("#select-ecm").should('contain', 'Cersei Lannister')
        cy.get("#select-ecm").should('contain', 'Eddard Stark')
        cy.get("#select-ecm").should('not.contain', 'Fake PA Deputy Name')
    });

    it("directs me back to dashboard page if I press cancel", () => {
        cy.get(".govuk-link").should("contain", "Cancel").click();
        cy.url().should('not.include', '/change-ecm')
        cy.get("h1").should("contain", "Dashboard");
    });

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