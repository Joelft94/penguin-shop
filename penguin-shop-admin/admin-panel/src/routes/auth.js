const express = require('express');
const bcrypt = require('bcryptjs');
const jwt = require('jsonwebtoken');
const Admin = require('../models/Admin');
const router = express.Router();


router.get('/login', (req, res) => {
  res.render('pages/login', { error: null });
});


router.post('/login', async (req, res) => {
  try {
    const { username, password } = req.body;

    
    const admin = await Admin.findOne({ username });
    if (!admin) {
      return res.render('pages/login', { 
        error: 'Invalid credentials' 
      });
    }

    
    const isMatch = await bcrypt.compare(password, admin.password);
    if (!isMatch) {
      return res.render('pages/login', { 
        error: 'Invalid credentials' 
      });
    }

    
    const token = jwt.sign(
      { id: admin._id },
      process.env.JWT_SECRET,
      { expiresIn: '24h' }
    );

    
    req.session.token = token;
    res.redirect('/dashboard');
  } catch (error) {
    res.render('pages/login', { 
      error: 'An error occurred during login' 
    });
  }
});

router.post('/logout', (req, res) => {
  req.session.destroy((err) => {
    if (err) {
      return res.redirect('/dashboard');
    }
    res.redirect('/login');
  });
});

module.exports = router;