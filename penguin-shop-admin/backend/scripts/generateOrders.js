// scripts/generateOrders.js
require('dotenv').config();
const mongoose = require('mongoose');
const Product = require('../src/models/Product');
const Order = require('../src/models/Order');

const iglooAddresses = [
    'Igloo #123, South Pole Street',
    'Igloo #456, Penguin Avenue',
    'Igloo #789, Arctic Circle',
    'Igloo #321, Ice Block B',
    'Igloo #654, Frozen Lake Road'
];

const generateRandomOrder = async (products) => {
    // With one product, we'll just randomize the quantity
    const product = products[0];
    const quantity = Math.floor(Math.random() * 3) + 1;
    const totalAmount = product.price * quantity;

    console.log(`Creating order with quantity: ${quantity}`);

    const order = new Order({
        iglooAddress: iglooAddresses[Math.floor(Math.random() * iglooAddresses.length)],
        items: [{
            product: product._id,
            quantity: quantity
        }],
        totalAmount: totalAmount,
        status: ['pending', 'processing', 'delivered'][Math.floor(Math.random() * 3)]
    });

    return order;
};

const generateOrders = async () => {
    try {
        await mongoose.connect(process.env.MONGODB_URI || 'mongodb://localhost:27017/penguin-shop');
        console.log('Connected to MongoDB');

        const products = await Product.find();
        console.log(`Found ${products.length} products:`);
        products.forEach(p => console.log(`- ${p.name} (${p._id})`));

        if (products.length === 0) {
            console.log('No products found. Please add some products first.');
            process.exit(1);
        }

        // Clear existing orders
        await Order.deleteMany({});
        console.log('Cleared existing orders');

        // Generate orders
        const numberOfOrders = 10;
        console.log(`\nGenerating ${numberOfOrders} test orders...`);

        for (let i = 0; i < numberOfOrders; i++) {
            console.log(`\nGenerating order ${i + 1} of ${numberOfOrders}`);
            const order = await generateRandomOrder(products);
            await order.save();
            console.log(`Created order: ${order.orderNumber} with total: $${order.totalAmount}`);
        }

        const finalCount = await Order.countDocuments();
        console.log(`\nTotal orders created: ${finalCount}`);
        
        process.exit(0);
    } catch (error) {
        console.error('Error generating orders:', error);
        process.exit(1);
    }
};

generateOrders();