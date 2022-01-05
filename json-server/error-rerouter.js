const getFailRoute = (req) => {
    return req.headers?.cookie?.match(/fail-route=(?<failRoute>\w+);/)?.groups
        .failRoute;
};

module.exports = (req, res, next) => {
    if (req.method === "POST") {
        const failRoute = getFailRoute(req);

        if (failRoute) {
            req.method = "GET";
            req.url = `/errors/${failRoute}`;
            res.status(400);
        }
    }
    next();
};
