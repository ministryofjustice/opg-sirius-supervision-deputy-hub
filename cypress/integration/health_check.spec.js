describe("Health check", () => {
    it("returns 200 status", () => {
        cy.request({
            url: "/supervision/deputies/health-check",
        }).then((resp) => {
            expect(resp.status).to.eq(200);
        });
    });
});
