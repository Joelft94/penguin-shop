package handlers

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"penguin-store/database"
	"penguin-store/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var templates *template.Template

func init() {
	// Debug loading of templates
	log.Printf("Loading templates...")
	
	templates = template.Must(template.ParseGlob("templates/*.html"))
	log.Printf("Templates loaded successfully")
}

type PageData struct {
	Products    []models.Product
	Product     *models.Product
	OrderNumber string
}

// HandleProducts shows the product listing
func HandleProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var products []models.Product
	cursor, err := database.DB.Collection("products").Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error fetching products: %v", err)
		http.Error(w, "Error fetching products", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &products); err != nil {
		log.Printf("Error parsing products: %v", err)
		http.Error(w, "Error parsing products", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Products: products,
	}

	log.Printf("Found %d products", len(products))
	
	// Debug template execution
	err = templates.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully rendered products page")
}


// HandleOrder handles both GET and POST requests for orders
func HandleOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handleGetOrder(w, r)
		return
	}
	if r.Method == "POST" {
		handlePostOrder(w, r)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func handleGetOrder(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("product_id")
	if productID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		log.Printf("Invalid product ID: %v", err)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	var product models.Product
	err = database.DB.Collection("products").FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		log.Printf("Error fetching product: %v", err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	data := PageData{
		Product: &product,
	}

	log.Printf("Rendering order template with product: %+v", product)
	log.Printf("Debug - Template data: %+v", data)
	err = templates.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully rendered order page")
}

func handlePostOrder(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	productID := r.FormValue("product_id")
	quantityStr := r.FormValue("quantity")
	iglooAddress := r.FormValue("iglooAddress")

	if iglooAddress == "" {
		http.Error(w, "Igloo address is required", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity < 1 {
		quantity = 1
	}

	ctx := context.Background()
	
	// Get product for price calculation and stock check
	var product models.Product
	err = database.DB.Collection("products").FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		log.Printf("Error fetching product: %v", err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Check stock
	if quantity > product.Stock {
		http.Error(w, "Not enough stock available", http.StatusBadRequest)
		return
	}

	// Calculate total amount
	totalAmount := product.Price * float64(quantity)

	// Create order
	order := models.Order{
		ID:           primitive.NewObjectID(),
		OrderNumber:  generateOrderNumber(),
		IglooAddress: iglooAddress,
		Items: []models.OrderItem{
			{
				ProductID: objID,
				Quantity:  quantity,
			},
		},
		TotalAmount: totalAmount,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	// Insert order
	_, err = database.DB.Collection("orders").InsertOne(ctx, order)
	if err != nil {
		log.Printf("Error creating order: %v", err)
		http.Error(w, "Error creating order", http.StatusInternalServerError)
		return
	}

	// Update product stock
	_, err = database.DB.Collection("products").UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$inc": bson.M{"stock": -quantity}},
	)
	if err != nil {
		log.Printf("Error updating stock: %v", err)
	}

	// Redirect to success page
	http.Redirect(w, r, "/order-success?order="+order.OrderNumber, http.StatusSeeOther)
}

func HandleOrderSuccess(w http.ResponseWriter, r *http.Request) {
	orderNumber := r.URL.Query().Get("order")
	if orderNumber == "" {
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}

	data := PageData{
		OrderNumber: orderNumber,
	}

	err := templates.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

func generateOrderNumber() string {
	timestamp := time.Now().Format("20060102")
	randomStr := primitive.NewObjectID().Hex()[:6]
	return "ORD-" + timestamp + "-" + randomStr
}