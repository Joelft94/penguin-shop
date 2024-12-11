const express = require('express');
const router = express.Router();
const auth = require('../middleware/auth');

router.get('/dashboard', auth, async (req, res) => {
    try {
        // For now, we'll use dummy data
        // Later we'll fetch this from the database
        const dashboardData = {
            username: 'joel', // We'll make this dynamic later
            totalProducts: 15,
            totalOrders: 48,
            todayRevenue: 2850,
            pendingOrders: 5,
            recentActivities: [
                {
                    description: 'New order #1234 received',
                    time: '5 minutes ago'
                },
                {
                    description: 'Product "Fresh Salmon" updated',
                    time: '1 hour ago'
                },
                {
                    description: 'New product "Ice Tray XL" added',
                    time: '2 hours ago'
                },
                {
                    description: 'Order #1233 marked as delivered',
                    time: '3 hours ago'
                }
            ]
        };

        res.render('pages/dashboard', dashboardData);
    } catch (error) {
        console.error('Dashboard error:', error);
        res.status(500).send('Error loading dashboard');
    }
});

module.exports = router;