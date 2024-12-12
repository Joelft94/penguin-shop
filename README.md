# 🐧 Penguin Shop
A full-stack e-commerce application built for penguins, featuring an admin panel and a customer-facing store. The project consists of two main components:

An Express.js-based admin dashboard
A Go-based customer storefront

## System Architecture
Admin Panel (Node.js/Express)

Authentication System: JWT-based authentication for admin users
MongoDB Integration: Direct database management for products and orders
Dashboard Features: Sales monitoring, order management, and product inventory control

Customer Store (Go)

Product Catalog: Public-facing store displaying available products
Order System: Shopping cart and order placement functionality
Template System: Go templates for dynamic HTML rendering

## Features
### Admin Dashboard

Secure login system with session management
Product management (CRUD operations)
Order tracking and status updates
Real-time inventory management
Sales statistics and reporting

### Customer Store

Browse products by category (fish, ice-trays, sleds)
Place orders with igloo delivery address
Real-time stock checking
Order confirmation system

### Database Schema

### Product Model
```
{
  name: String,
  description: String,
  price: Number,
  stock: Number,
  type: Enum['fish', 'ice-tray', 'sled'],
  createdAt: Date,
  updatedAt: Date
}
```

### Order Model

```
{
  orderNumber: String,
  items: [{
    product: ObjectId,
    quantity: Number
  }],
  iglooAddress: String,
  totalAmount: Number,
  status: Enum['pending', 'processing', 'delivered'],
  createdAt: Date
}
```
### Admin Model
```
{
  username: String,
  password: String (hashed),
  createdAt: Date,
  updatedAt: Date
}
```
## Setup Instructions

### Prerequisites

Node.js (v14 or higher)
Go (v1.16 or higher)
MongoDB (v4.4 or higher)

### Admin Panel Setup

Install dependencies:
```
npm install
```
Set up environment variables:

```
MONGODB_URI=mongodb://localhost:27017/penguin-shop
JWT_SECRET=your-secret-key
SESSION_SECRET=your-session-secret
```
Create initial admin user:
```
node scripts/createAdmin.js
```
Generate sample orders (optional):
```
node scripts/generateOrders.js
```
Customer Store Setup

Set environment variables:
```
MONGODB_URI=mongodb://localhost:27017
PORT=8080
```
Build and run:
```
go build
./penguin-store
```


