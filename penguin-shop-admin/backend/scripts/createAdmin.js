require('dotenv').config();
const mongoose = require('mongoose');
const Admin = require('../src/models/Admin');

const createInitialAdmin = async () => {
    try {
        await mongoose.connect(process.env.MONGODB_URI || 'mongodb://localhost:27017/penguin-shop');
        
    
        await Admin.deleteOne({ username: 'paula' });

        
        const admin = new Admin({
            username: 'joel',
            password: 'joel123' // This will be automatically hashed by the pre-save middleware
        });

        await admin.save();
        console.log('Admin user created successfully!');
        console.log('Username: joel');
        console.log('Password: joel123');
        
        await mongoose.connection.close();
        process.exit(0);
    } catch (error) {
        console.error('Error:', error);
        process.exit(1);
    }
};

createInitialAdmin();