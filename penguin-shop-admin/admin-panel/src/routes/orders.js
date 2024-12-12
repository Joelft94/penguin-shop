const express = require('express');
const router = express.Router();
const Order = require('../models/Order');
const auth = require('../middleware/auth');

// List all orders
router.get('/orders', auth, async (req, res) => {
    try {
        const orders = await Order.find()
            .populate('items.product')
            .sort({ createdAt: -1 });
        res.render('pages/orders/index', { orders });
    } catch (error) {
        console.error('Error fetching orders:', error);
        res.status(500).send('Error loading orders');
    }
});

// View single order
router.get('/orders/:id', auth, async (req, res) => {
    try {
        const order = await Order.findById(req.params.id)
            .populate('items.product');
        if (!order) {
            return res.redirect('/orders');
        }
        res.render('pages/orders/view', { order });
    } catch (error) {
        res.redirect('/orders');
    }
});

// Update order status
router.post('/orders/:id/status', auth, async (req, res) => {
    try {
        const { status } = req.body;
        await Order.findByIdAndUpdate(req.params.id, { status });
        res.redirect('/orders');
    } catch (error) {
        res.status(500).send('Error updating order status');
    }
});

module.exports = router;