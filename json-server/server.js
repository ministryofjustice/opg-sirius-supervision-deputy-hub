const jsonServer = require("json-server")

const server = jsonServer.create()
const router = jsonServer.router("db.json")
const middlewares = jsonServer.defaults()

server.use(...middlewares)

const getFailRoute = (req) => {
    if (req.headers && req.headers.cookie) {
        const match = req.headers.cookie.match(/fail-route=(?<failRoute>\w+);/)

        if (match && match.groups) {
            return req.headers.cookie.match(/fail-route=(?<failRoute>\w+);/).groups.failRoute
        }
    }
}

const updatePath = (req, failRoute) => {
    const re = /(?<=\/api\/v1\/)(.*)/

    req.url = req.originalUrl.replace(re, `errors/${failRoute}`)
}

server.use((req, res, next) => {
    if (req.method === "POST") {
        const failRoute = getFailRoute(req)

        if (failRoute) {
            console.log('redirecting')
            req.method = "GET"
            res.status(400)
            updatePath(req, failRoute)
        }
    }
    next()
})

server.use(jsonServer.rewriter({
    "/api/v1/*": "/$1",
    "/deputies/pa/:personId/clients": "/clients",
    "/timeline/:personId": "/timeline",
    "/deputy/:personId/notes": "/notes",
    "/users/current": "/users"
}))

server.use(router)
server.listen(3000, "0.0.0.0", () => {
    console.log("JSON Server is running on http://localhost:3000/api/v1/")
})