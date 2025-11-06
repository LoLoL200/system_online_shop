package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Structur
type Product struct {
	Name        string
	Description string
	Price       float64
	Category    string
	Stock       int
	ID          int
}

type Client struct {
	ID     int
	Name   string
	Status string
}

type Order struct {
	ID        int
	ClientID  int
	Items     []Product
	TotalCost float64
	Status    string
}

type Cart struct {
	ClientID int
	Items    []Product
}

type Store struct {
	Name      string
	Products  []Product
	Customers []Client
	Orders    []Order
	Carts     map[int]Cart
}

// initializeStore
func initializeStore() *Store {
	return &Store{
		Name:      "TechStore",
		Products:  make([]Product, 0),
		Customers: make([]Client, 0),
		Orders:    make([]Order, 0),
		Carts:     make(map[int]Cart),
	}
}

var reader = bufio.NewReader(os.Stdin)

// fF
func getFloatInput(prompt string) float64 {
	for {
		fmt.Print(prompt)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		value, err := strconv.ParseFloat(text, 64)
		if err != nil {
			fmt.Println("Введіть число, наприклад 333.33")
			continue
		}
		return value
	}
}

// fI
func getIntInput(prompt string) int {
	for {
		fmt.Print(prompt)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		value, err := strconv.Atoi(text)
		if err != nil {
			fmt.Println("Введіть ціле число, наприклад 666")
			continue
		}
		return value
	}
}

// fT
func getTextInput(prompt string) string {

	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)

}

// Place Order
func placeOrder(store *Store) {
	fmt.Println("============================Оформлення замовлення ============================")

	clientID := getIntInput("Введіть ID клієнта: ")

	// Search client
	var client *Client
	for i := range store.Customers {
		if store.Customers[i].ID == clientID {
			client = &store.Customers[i]
			break
		}
	}

	// Checking for absence of client
	if client == nil {
		fmt.Println("Клієнта з таким ID не знайдено.")
		return
	}

	// Get card
	cart, exists := store.Carts[clientID]
	if !exists || len(cart.Items) == 0 {
		fmt.Println("Кошик порожній або не існує.")
		return
	}

	// Calculate the order amount
	total := 0.0
	for _, item := range cart.Items {
		total += item.Price * float64(item.Stock)
	}

	// Status discount
	var discountRate float64
	switch client.Status {
	case "Regular":
		discountRate = 0.05
	case "VIP":
		discountRate = 0.10
	default:
		discountRate = 0.0
	}

	discountAmount := total * discountRate
	finalPrice := total - discountAmount

	// Update the remaining goods in stock
	for _, item := range cart.Items {
		for i := range store.Products {
			if store.Products[i].ID == item.ID {
				if store.Products[i].Stock < item.Stock {
					fmt.Printf("Недостатньо товару '%s' на складі.\n", item.Name)
					return
				}
				store.Products[i].Stock -= item.Stock
			}
		}
	}

	// Add Order
	order := Order{
		ID:        len(store.Orders) + 1,
		ClientID:  clientID,
		Items:     cart.Items,
		TotalCost: finalPrice,
		Status:    "Оформлено",
	}

	store.Orders = append(store.Orders, order)

	// Delete card
	delete(store.Carts, clientID)

	fmt.Println("===== Замовлення оформлено! =====")
	fmt.Printf("Номер замовлення: %d\n", order.ID)
	fmt.Printf("Клієнт: %s (%s)\n", client.Name, client.Status)
	fmt.Printf("Кількість товарів: %d\n", len(order.Items))
	fmt.Printf("Сума без знижки: %.2f грн\n", total)
	fmt.Printf("Знижка: %.0f%% (%.2f грн)\n", discountRate*100, discountAmount)
	fmt.Printf("До сплати: %.2f грн\n", finalPrice)
	fmt.Println("=================================")
}

func shopping(store *Store) {
	clientID := getIntInput("Введіть ID клієнта для перегляду кошика: ")
	//  Does the client exist?
	var clientExists bool
	for _, c := range store.Customers {
		if c.ID == clientID {
			clientExists = true
			break
		}
	}
	if !clientExists {
		fmt.Println("Клієнта з таким ID не знайдено.")
		return
	}
	// Check for card to client
	cart, exists := store.Carts[clientID]
	if !exists {
		fmt.Println("У цього клієнта немає кошика.")
		return
	}

	// Check for empty card client
	fmt.Printf("-- У кошику клієнта ID %d --\n", clientID)
	if len(cart.Items) == 0 {
		fmt.Println("Кошик порожній.")
		return
	}

	// Displaying the customer's shopping cart
	for i, item := range cart.Items {
		fmt.Printf("%d. %s (%d шт.) - %.2f грн/шт\n", i+1, item.Name, item.Stock, item.Price)
	}
}

func discount(store *Store) {
	fmt.Println("============================Застосування знижки ============================")

	clientID := getIntInput("Введіть ID клієнта: ")

	// Search client
	var client *Client
	for i := range store.Customers {
		if store.Customers[i].ID == clientID {
			client = &store.Customers[i]
			break
		}
	}

	// Check for client void
	if client == nil {
		fmt.Println("Клієнта з таким ID не знайдено.")
		return
	}

	// Get shopping card
	cart, exists := store.Carts[clientID]
	if !exists || len(cart.Items) == 0 {
		fmt.Println("Кошик порожній або не існує.")
		return
	}

	// Calculate the order amount
	total := 0.0
	for _, item := range cart.Items {
		total += item.Price * float64(item.Stock)
	}

	// Determine the discount amount
	var discountRate float64
	switch client.Status {
	case "Regular":
		discountRate = 0.05
	case "VIP":
		discountRate = 0.10
	default:
		discountRate = 0.0
	}

	// Calculating the discounted amount
	discountAmount := total * discountRate
	finalPrice := total - discountAmount

	// Output the result
	fmt.Println("===== Розрахунок знижки =====")
	fmt.Printf("Клієнт: %s (Статус: %s)\n", client.Name, client.Status)
	fmt.Printf("Сума без знижки: %.2f грн\n", total)
	fmt.Printf("Знижка: %.0f%% (%.2f грн)\n", discountRate*100, discountAmount)
	fmt.Printf("До сплати: %.2f грн\n", finalPrice)
	fmt.Println("==============================")
}

func removeShopping(store *Store) {
	fmt.Println("--- Видалення з кошика ---")
	clientID := getIntInput("Введіть ID клієнта: ")

	// Checking if the client exists
	var clientExists bool
	for _, c := range store.Customers {
		if c.ID == clientID {
			clientExists = true
			break
		}
	}

	// Check for card to client ID
	if !clientExists {
		fmt.Println("Клієнта з таким ID не знайдено.")
		return
	}

	// Get shopping card client
	cart, exists := store.Carts[clientID]
	if !exists || len(cart.Items) == 0 {
		fmt.Println("Кошик порожній або не існує.")
		return
	}

	// Showing the contents of the cart
	fmt.Println("Товари у кошику:")
	for i, item := range cart.Items {
		fmt.Printf("%d. %s (%d шт.) - %.2f грн/шт [ID:%d]\n", i+1, item.Name, item.Stock, item.Price, item.ID)
	}

	productID := getIntInput("Введіть ID товару для видалення: ")

	// Search item to shopping cart
	found := false
	newItems := make([]Product, 0)
	for _, item := range cart.Items {
		if item.ID == productID {
			found = true
			continue
		}
		newItems = append(newItems, item)
	}

	// Product inspection
	if !found {
		fmt.Println("Товар з таким ID не знайдено у кошику.")
		return
	}

	// Update shopping cart
	cart.Items = newItems
	store.Carts[clientID] = cart

	fmt.Printf("Товар з ID %d видалено з кошика клієнта ID %d!\n", productID, clientID)
}

func addShopping(store *Store) {
	fmt.Println("============================Додавання до кошика ============================")

	clientID := getIntInput("Введіть ID клієнта: ")

	//  Does the client exist?
	var clientExists bool
	for _, c := range store.Customers {
		if c.ID == clientID {
			clientExists = true
			break
		}
	}
	if !clientExists {
		fmt.Println("Клієнта з таким ID не знайдено.")
		return
	}

	productID := getIntInput("Введіть ID товару: ")
	quantity := getIntInput("Введіть кількість: ")

	// Search for product
	var selectedProduct *Product
	for i := range store.Products {
		if store.Products[i].ID == productID {
			selectedProduct = &store.Products[i]
			break
		}
	}
	if selectedProduct == nil {
		fmt.Println("Товар з таким ID не знайдено.")
		return
	}

	// Checking for availability in stock
	if quantity > selectedProduct.Stock {
		fmt.Printf(" Недостатньо товару '%s' на складі. Є лише %d шт.\n",
			selectedProduct.Name, selectedProduct.Stock)
		return
	}

	// Create a new or retrieve a customer cart
	cart, exists := store.Carts[clientID]
	if !exists {
		cart = Cart{ClientID: clientID, Items: []Product{}}
	}

	// Add the product (copy) with the required quantity
	productCopy := *selectedProduct
	productCopy.Stock = quantity
	cart.Items = append(cart.Items, productCopy)

	// Update cart in store
	store.Carts[clientID] = cart

	fmt.Printf("Товар \"%s\" (%d шт.) додано до кошика клієнта ID %d!\n",
		selectedProduct.Name, quantity, clientID)
}

func shoppingCard(store *Store) {
	fmt.Println("============================ Меню кошика ============================")
	clientID := getIntInput("Ведіть ID клієнта:")

	clientExists := false
	for _, c := range store.Customers {
		if c.ID == clientID {
			clientExists = true
			break
		}
	}

	if !clientExists {
		fmt.Println(" Клієнта з таким ID не знайдено. Повернення до головного меню.")
		return
	}

	fmt.Printf("1. Додати товар до кошика\n2. Видалити товар з кошика\n3. Переглянути кошик\n4. Застосувати знижку\n5. Оформити замовлення\n6. Повернутися до головного меню\n")
	inputCard := getIntInput("> ")
	switch inputCard {
	case 1:
		addShopping(store)
	case 2:
		removeShopping(store)
	case 3:
		shopping(store)
	case 4:
		discount(store)
	case 5:
		placeOrder(store)
	case 6:
		menuShop(store)
	}
}

// 2

// redact Status
func redactStatus(store *Store) {
	fmt.Println("============================Зміна статусу============================")
	radact := getTextInput("Хочите змінити статус?")
	if radact != "ТАК" && radact != "так" && radact != "Так" {
		fmt.Println("Ви повернитись з початку")
		return
	}
	inputIdClient := getIntInput("Введіть ID клієнта:")
	found := false
	for i := range store.Customers {
		if store.Customers[i].ID == inputIdClient {
			found = true
			renameStatus := getTextInput("Введіть новий статус клієнта (Base, Regular, VIP):")

			if renameStatus != "Base" && renameStatus != "Regular" && renameStatus != "VIP" {
				fmt.Println("Невірний статус, спробуйте ще раз")
				return
			}

			store.Customers[i].Status = renameStatus
			fmt.Printf("Статус клієнта з ID %d змінено на %s.\n", inputIdClient, renameStatus)
			break
		}
	}

	if !found {
		fmt.Printf("Клієнта з ID %d не знайдено.\n", inputIdClient)
	}

}

// Search client by ID
func searchClientID(store *Store) {
	searchID := getIntInput("Введіть ID клієнта:")
	for _, client := range store.Customers {
		if searchID == client.ID {
			fmt.Printf("%d Name:%s\n Status: %s \n",
				client.ID, client.Name, client.Status)
		} else if searchID != client.ID {
			fmt.Println("Немає такого ID, ведіть будь ласка правильний ID")
		} else {
			fmt.Println("Ведіть будьласка число")
		}

	}
}

// All clients
func allClient(store *Store) {
	if len(store.Customers) == 0 {
		fmt.Println("Клієнтів немає(")
	}
	fmt.Println("============================Список клієнтів============================")
	for i, client := range store.Customers {
		fmt.Printf("%d Name: %s\n Status:%s\n",
			i+1, client.Name, client.Status)
	}
}

// Add new Clients
func addClient(store *Store) {
	nameClient := getTextInput("Введіть назву клієнта:")
	status := getTextInput("Введіть статус  клієнта (тільки ці три варіанти Base,Regular,VIP):")

	if status != "Base" && status != "Regular" && status != "VIP" {
		fmt.Println(" Невірний статус клієнта! Спробуйте ще раз.")
		return
	}

	client := Client{
		ID:     len(store.Customers),
		Name:   nameClient,
		Status: status,
	}
	store.Customers = append(store.Customers, client)
	fmt.Printf("%d Додан до бази клієнт: Ім'я клієнта: %s\n статус:%s\n", client.ID, client.Name, client.Status)

}

// Menu for clients
func menuClient(store *Store) {
	fmt.Println("============================Меню клієнта ============================")
	fmt.Printf("1. Додати клієнта\n2. Переглянути всіх клієнтів\n3. Знайти клієнта за ID\n4. Змінити статус клієнта\n5. Повернутися до головного меню\n")
	inputClient := getIntInput("> ")
	switch inputClient {
	case 1:
		addClient(store)
	case 2:
		allClient(store)
	case 3:
		searchClientID(store)
	case 4:
		redactStatus(store)
	case 5:
		menuShop(store)
	}
}

// 1

// Search product by Category
func searchByCategory(store *Store) {
	searchID := getTextInput("Введіть категорію продука:")

	found := false
	for _, product := range store.Products {
		if searchID == product.Category {
			fmt.Printf("ID: %d\nName: %s\nDescription: %s\nPrice: %.2f\nCategory: %s\nStock: %d\n",
				product.ID, product.Name, product.Description, product.Price, product.Category, product.Stock)
			found = true
		}
	}
	if !found {
		fmt.Println("Немає товарів у цій категорії")
	}

}

// Search product by ID
func searchByID(store *Store) {
	searchID := getIntInput("Введіть ID продука:")
	fmt.Println("", searchID)
	for _, product := range store.Products {
		if searchID == product.ID {
			fmt.Printf("ID: %d\nName: %s\nDescription: %s\nPrice: %.2f\nCategory: %s\nStock: %d\n",
				product.ID, product.Name, product.Description, product.Price, product.Category, product.Stock)
			return
		}
	}
	fmt.Println("Товар з таким ID не знайдено")

}

// Add Product
func addProduct(store *Store) {
	product := Product{
		ID:          len(store.Products) + 1,
		Name:        getTextInput("Введіть назву товару:"),
		Description: getTextInput("Введіть опис:"),
		Price:       getFloatInput("Введіть ціну:"),
		Category:    getTextInput("Введіть категорію:"),
		Stock:       getIntInput("Введіть кількість на складі:"),
	}
	product.ID = len(store.Products) + 1
	store.Products = append(store.Products, product)
	fmt.Printf("%d Додан до бази продукт: %s\n", product.ID, product.Name)
}

// All list Product
func allProduct(store *Store) {
	if len(store.Products) == 0 {
		fmt.Println("Товарів не знайдено")
	}
	fmt.Println("============================Список товарів============================")
	for i, product := range store.Products {
		fmt.Printf("%d Name: %s\n Description:%s\n Prise%.2f\n  Category:%s\n QuantityInStock:%d \n",
			i+1, product.Name, product.Description, product.Price, product.Category, product.Stock)
	}
}

// Update only price and stock
func updateProduct(store *Store) {
	fmt.Println("============================Оновлення товару (ціна та кількість)============================")

	if len(store.Products) == 0 {
		fmt.Println("Немає жодного товару для оновлення.")
		return
	}

	productID := getIntInput("Введіть ID товару, який потрібно оновити: ")

	var product *Product
	for i := range store.Products {
		if store.Products[i].ID == productID {
			product = &store.Products[i]
			break
		}
	}

	if product == nil {
		fmt.Println("Товар з таким ID не знайдено.")
		return
	}

	fmt.Printf("Ви обрали товар: %s (%.2f грн, Кількість: %d)\n",
		product.Name, product.Price, product.Stock)

	// Update price and quantity
	newPrice := getFloatInput("Введіть нову ціну: ")
	newStock := getIntInput("Введіть нову кількість на складі: ")

	product.Price = newPrice
	product.Stock = newStock

	fmt.Println(" Товар успішно оновлено!")
	fmt.Printf("ID: %d\nНазва: %s\nНова ціна: %.2f грн\nНова кількість: %d шт\n",
		product.ID, product.Name, product.Price, product.Stock)
}

// Menu for products
func menuProducts(store *Store) {

	fmt.Println("============================ Меню товарів ============================")
	fmt.Printf("1. Додати товар\n2. Переглянути всі товари\n3. Знайти товар за ID\n")
	fmt.Printf("4. Пошук за категорією\n5.Оновити товари\n6. Повернутися до головного меню\n")
	inputProduct := getIntInput("> ")
	switch inputProduct {
	case 1:
		addProduct(store)
	case 2:
		allProduct(store)
	case 3:
		searchByID(store)
	case 4:
		searchByCategory(store)
	case 5:
		updateProduct(store)
	case 6:
		menuShop(store)
	}

}
func booking(store *Store) {
	if len(store.Orders) == 0 {
		fmt.Println("Замовлень ще немає.")
		return
	}
	for _, order := range store.Orders {
		// Search client
		var clientName string
		for _, c := range store.Customers {
			if c.ID == order.ClientID {
				clientName = c.Name
				break
			}
		}
		if clientName == "" {
			clientName = "Невідомий клієнт"
		}

		fmt.Printf("\n============================\n")
		fmt.Printf("Номер замовлення: %d\n", order.ID)
		fmt.Printf("Клієнт: %s (ID: %d)\n", clientName, order.ClientID)
		fmt.Printf("Статус замовлення: %s\n", order.Status)
		fmt.Println("Товари:")

		for i, item := range order.Items {
			fmt.Printf("  %d. %s — %d шт × %.2f грн = %.2f грн\n",
				i+1, item.Name, item.Stock, item.Price, item.Price*float64(item.Stock))
		}

		fmt.Printf("Загальна сума (зі знижкою): %.2f грн\n", order.TotalCost)
		fmt.Println("============================")
	}

}

// Store statistics: number of clients and products
func storeStatistics(store *Store) {
	fmt.Println("=== Статистика магазину ===")

	totalClients := len(store.Customers)
	totalProducts := len(store.Products)

	fmt.Printf("Кількість клієнтів: %d\n", totalClients)
	fmt.Printf("Кількість товарів: %d\n", totalProducts)

	fmt.Println("============================")
}

// Menu shop
func menuShop(store *Store) {
	fmt.Printf("====Головне меню:=====\n1. Управління товарами\n2. Управління клієнтами\n3. Кошик покупок\n")
	fmt.Printf("4. Замовлення\n5. Статистика магазину\n6. Вихід\n============================\n")
	inputShop := getIntInput("> ")
	switch inputShop {
	case 1:
		menuProducts(store)
	case 2:
		menuClient(store)
	case 3:
		shoppingCard(store)
	case 4:
		booking(store)
	case 5:
		storeStatistics(store)
	case 6:
		os.Exit(0)
	}
}

func main() {
	store := initializeStore()
	for {
		fmt.Println("=== Онлайн-магазин TechStore ===")
		menuShop(store)
	}

}
