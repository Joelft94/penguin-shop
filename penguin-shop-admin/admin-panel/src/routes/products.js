const express = require('express');
const router = express.Router();
const Product = require('../models/Product');
const auth = require('../middleware/auth');


router.get('/products', auth, async (req, res) => {
    try {
        const products = await Product.find().sort({ createdAt: -1 });
        res.render('pages/products/index', { products });
    } catch (error) {
        res.status(500).send('Error loading products');
    }
});


router.get('/products/new', auth, (req, res) => {
    res.render('pages/products/new');
});


router.post('/products', auth, async (req, res) => {
    try {
        const product = new Product(req.body);
        await product.save();
        res.redirect('/products');
    } catch (error) {
        res.render('pages/products/new', { 
            error: 'Error creating product',
            product: req.body 
        });
    }
});


router.get('/products/:id/edit', auth, async (req, res) => {
    try {
        const product = await Product.findById(req.params.id);
        if (!product) {
            return res.redirect('/products');
        }
        res.render('pages/products/edit', { product });
    } catch (error) {
        console.error('Edit form error:', error);
        res.redirect('/products');
    }
});


router.post('/products/:id', auth, async (req, res) => {
    try {
        const product = await Product.findById(req.params.id);
        if (!product) {
            return res.redirect('/products');
        }

        // Update product fields
        product.name = req.body.name;
        product.description = req.body.description;
        product.price = req.body.price;
        product.stock = req.body.stock;
        product.type = req.body.type;
        product.updatedAt = Date.now();

        await product.save();
        res.redirect('/products');
    } catch (error) {
        console.error('Update error:', error);
        res.redirect(`/products/${req.params.id}/edit`);
    }
});



router.post('/products/:id/delete', auth, async (req, res) => {
    try {
        await Product.findByIdAndDelete(req.params.id);
        res.redirect('/products');
    } catch (error) {
        res.status(500).send('Error deleting product');
    }
});

module.exports = router;