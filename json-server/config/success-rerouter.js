const getSuccessRoute = (req) => {
    return req.headers?.cookie?.match(/success-route=(?<successRoute>[\w\/]+);/)?.groups.successRoute
}

module.exports = (req, res, next) => {
    if (["POST", "PUT", "PATCH"].includes(req.method)) {
        const successRoute = getSuccessRoute(req);

        if (req) {
            req.method === "POST" ? res.status(201) : res.status(200);
            req.method = "GET";
            req.url = successRoute;
        }
    }
    next();
};
