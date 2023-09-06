describe("Error handling", () => {
    it("renders a 404-specific error page", () => {
        let urls = [
            "/supervision/deputies/nothing-here/123", // invalid route
            "/supervision/deputies/123" // valid route, non-existent deputy
        ]

        urls.forEach((url) => {
            cy.intercept(url).as("request")
            cy.visit(url, {failOnStatusCode: false});
            cy.wait("@request").then((response) => {
                expect(response.response.statusCode).to.eq(404)
            })
            cy.get("body").should("have.class", "sirius-pahub")
            cy.contains(".govuk-heading-l", "Page not found");
            cy.get(".govuk-body").should("contain", "If you typed the web address, check it is correct.")
            cy.get(".govuk-body").should("contain", "If you pasted the web address, check you copied the entire address.")
            cy.get(".govuk-body").should("contain", "Please use your browser to go back to the previous page, or return to the homepage.")
        })
    });

    describe("Non-404 error pages", () => {
        beforeEach(() => {
            cy.setCookie("Other", "other");
            cy.setCookie("XSRF-TOKEN", "abcde");
            cy.setCookie("fail-route", "500-example");
            cy.visit("/supervision/deputies/2/notes/add-note");
            cy.intercept("POST", "/supervision/deputies/2/notes/add-note").as("submit")
        });

        it("renders a 403-specific error page", () => {
            cy.setCookie("fail-code", "403");
            cy.get("#add-note-form").submit();
            cy.wait("@submit").then((submit) => {
                expect(submit.response.statusCode).to.eq(403)
            })
            cy.get("body").should("have.class", "sirius-pahub")
            cy.contains(".govuk-heading-l", "Forbidden");
            cy.get(".govuk-body").should("contain", "You do not have access to view this page.")
            cy.get(".govuk-body").should("contain", "Please use your browser to go back to the previous page, or return to the homepage.")
        })

        it("renders a generic error page", () => {
            cy.setCookie("fail-code", "503");
            cy.get("#add-note-form").submit();
            cy.wait("@submit").then((submit) => {
                expect(submit.response.statusCode).to.eq(503)
            })
            cy.get("body").should("have.class", "sirius-pahub")
            cy.contains(".govuk-heading-l", "Sorry, there is a problem with the service");
            cy.get(".govuk-body").should("contain", "Try again later.")
            cy.get(".govuk-body").invoke("text").should("match", /Further information: POST http(.*) returned 503/)
        })
    })
});
