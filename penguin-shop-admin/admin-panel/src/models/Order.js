const mongoose = require('mongoose');

const orderSchema = new mongoose.Schema({
    orderNumber: {
        type: String,
        unique: true
    },
    items: [{
        product: {
            type: mongoose.Schema.Types.ObjectId,
            ref: 'Product',
            required: true
        },
        quantity: {
            type: Number,
            required: true,
            min: 1
        }
    }],
    iglooAddress: {
        type: String,
        required: true
    },
    totalAmount: {
        type: Number,
        required: true
    },
    status: {
        type: String,
        enum: ['pending', 'processing', 'delivered'],
        default: 'pending'
    },
    createdAt: {
        type: Date,
        default: Date.now
    }
});

// Generate order number
orderSchema.pre('validate', async function(next) {
    if (this.orderNumber) return next();
    
    const date = new Date();
    const dateStr = date.getFullYear().toString() +
        (date.getMonth() + 1).toString().padStart(2, '0') +
        date.getDate().toString().padStart(2, '0');
    
    const lastOrder = await this.constructor.findOne({}, {}, { sort: { 'createdAt': -1 } });
    let sequence = '001';
    
    if (lastOrder && lastOrder.orderNumber) {
        const lastSequence = parseInt(lastOrder.orderNumber.split('-')[2]);
        sequence = (lastSequence + 1).toString().padStart(3, '0');
    }
    
    this.orderNumber = `ORD-${dateStr}-${sequence}`;
    next();
});

module.exports = mongoose.model('Order', orderSchema);