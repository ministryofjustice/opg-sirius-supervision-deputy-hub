import navTabs from "../fixtures/navigation.json";
import "cypress-axe";


describe("Accessibility", () => {
    navTabs.forEach(([page, url]) =>
        it(`should render ${page} page accessibly`, () => {
            cy.visit(url);
            cy.injectAxe();
            cy.checkA11y();
        })
    )
});
