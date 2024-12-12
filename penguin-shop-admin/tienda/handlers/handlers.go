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
	log.Printf("Loading templates...")

	// Use more explicit parsing
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Template parsing error: %v", err)
	}

	templates = tmpl

	// Verify template names
	if templates == nil {
		log.Fatal("No templates found!")
	}

	log.Printf("Parsed templates: %v", templates.DefinedTemplates())
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
	log.Printf("=== FULL REQUEST DETAILS ===")
	log.Printf("Method: %s", r.Method)
	log.Printf("URL: %s", r.URL.String())
	log.Printf("Content-Type: %s", r.Header.Get("Content-Type"))

	for key, values := range r.Header {
		for _, value := range values {
			log.Printf("Header - %s: %s", key, value)
		}
	}

	switch r.Method {
	case http.MethodGet:
		handleGetOrder(w, r)
	case http.MethodPost:
		handlePostOrder(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("=== handleGetOrder START ===")

	productID := r.URL.Query().Get("product_id")
	log.Printf("Product ID from query: %s", productID)

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

	log.Printf("Product fetched: %+v", product)

	data := PageData{
		Product: &product,
	}

	log.Println("Attempting to render 'simple.html' template")
	err = templates.ExecuteTemplate(w, "simple.html", data)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}

	log.Println("=== handleGetOrder END ===")
}

func handlePostOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("=== POST ORDER ATTEMPT ===")

	// Parse form with maximum memory
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Printf("FORM PARSE ERROR: %v", err)
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}

	log.Println("FORM CONTENTS:")
	for key, values := range r.Form {
		for _, value := range values {
			log.Printf("  %s: %s", key, value)
		}
	}

	// Extract form values
	productID := r.FormValue("product_id")
	quantityStr := r.FormValue("quantity")
	iglooAddress := r.FormValue("iglooAddress")

	log.Printf("EXTRACTED VALUES:")
	log.Printf("  Product ID: %s", productID)
	log.Printf("  Quantity: %s", quantityStr)
	log.Printf("  Igloo Address: %s", iglooAddress)

	// Rest of the existing implementation...
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		log.Printf("INVALID PRODUCT ID: %v", err)
		http.Error(w, "Invalid product", http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity < 1 {
		log.Printf("QUANTITY PARSE ERROR: %v", err)
		quantity = 1
	}

	// Fetch product
	var product models.Product
	err = database.DB.Collection("products").FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		log.Printf("PRODUCT FETCH ERROR: %v", err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

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
		TotalAmount: product.Price * float64(quantity),
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	// Insert order
	result, err := database.DB.Collection("orders").InsertOne(ctx, order)
	if err != nil {
		log.Printf("ORDER INSERT ERROR: %v", err)
		http.Error(w, "Error creating order", http.StatusInternalServerError)
		return
	}

	log.Printf("ORDER INSERTED: ID = %v", result.InsertedID)
	log.Println("=== handlePostOrder END ===")

	count, countErr := database.DB.Collection("orders").CountDocuments(ctx, bson.M{"orderNumber": order.OrderNumber})
	if countErr != nil {
		log.Printf("Count error: %v", countErr)
	}
	log.Printf("Orders with this number: %d", count)

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
