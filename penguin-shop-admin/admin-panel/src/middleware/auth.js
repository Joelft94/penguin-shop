const jwt = require('jsonwebtoken');

const auth = async (req, res, next) => {
  try {

    if (!req.session.token) {
      return res.redirect('/login');
    }

    const decoded = jwt.verify(req.session.token, process.env.JWT_SECRET);
    req.adminId = decoded.id;
    next();
  } catch (error) {
    res.redirect('/login');
  }
};

module.exports = auth;