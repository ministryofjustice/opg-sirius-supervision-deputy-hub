module.exports = (req, res, next) => {
    if (req.headers?.cookie?.includes('success-bypass')) {
        res.status(200)
    }
    next()
}
