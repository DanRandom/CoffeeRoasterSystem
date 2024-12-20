package orders

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Order struct {
	ID          int64   `json:"id"`
	UserID      int64   `json:"user_id"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type OrderItems struct {
	ID        int64   `json:"id"`
	OrderID   int64   `json:"order_id"`
	CoffeeID  int64   `json:"coffee_id"`
	Quantity  int64   `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

// InspectOrders fetches all pending orders and returns them
func InspectOrders(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Fetch pending orders
		rows, err := db.Query(`
            SELECT id, user_id, total_amount, status, created_at, updated_at
            FROM orders
            WHERE status = 'pending'
        `)
		if err != nil {
			log.Printf("Error fetching orders: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to load orders. Please try again later.",
			})
			return
		}
		defer rows.Close()

		var orders []Order
		for rows.Next() {
			var order Order
			if err := rows.Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.Status, &order.CreatedAt, &order.UpdatedAt); err != nil {
				log.Printf("Error scanning order: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Error loading order data.",
				})
				return
			}
			orders = append(orders, order)
		}

		c.JSON(http.StatusOK, gin.H{
			"title":  "Pending Orders",
			"orders": orders,
		})
	}
}

// InspectOrders fetches all pending orders and returns them
func InspectAllOrders(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Fetch pending orders
		rows, err := db.Query(`
            SELECT id, user_id, total_amount, status, created_at, updated_at
            FROM orders`)
		if err != nil {
			log.Printf("Error fetching orders: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to load orders. Please try again later.",
			})
			return
		}
		defer rows.Close()

		var orders []Order
		for rows.Next() {
			var order Order
			if err := rows.Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.Status, &order.CreatedAt, &order.UpdatedAt); err != nil {
				log.Printf("Error scanning order: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Error loading order data.",
				})
				return
			}
			orders = append(orders, order)
		}

		c.JSON(http.StatusOK, gin.H{
			"title":  "Pending Orders",
			"orders": orders,
		})
	}
}

// UpdateOrderStatus updates the status of a specific order
func UpdateOrderStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("id")
		var request struct {
			Status string `json:"status" binding:"required"`
		}

		// Parse the new status from the request body
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input: Status is required",
			})
			return
		}

		// Update order status in the database
		_, err := db.Exec(`
            UPDATE orders
            SET status = ?, updated_at = datetime('now')
            WHERE id = ?
        `, request.Status, orderID)

		if err != nil {
			log.Printf("Error updating order status: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update order status. Please try again later.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Order status updated successfully",
		})
	}
}

func DeleteOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coffee ID"})
			return
		}

		query := "DELETE FROM orders WHERE id = ?"
		res, err := db.Exec(query, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"title":          "RoastingRooster Coffee Inventory",
			"dbActionStatus": "SUCCESS",
		})
	}
}

func GetUserOrders(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		rows, err := db.Query(`
            SELECT id, user_id, total_amount, status, created_at, updated_at
            FROM orders
            WHERE user_id = ?
        `, id)
		if err != nil {
			log.Printf("Error fetching orders: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to load orders. Please try again later.",
			})
			return
		}
		defer rows.Close()

		var orders []Order
		for rows.Next() {
			var order Order
			if err := rows.Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.Status, &order.CreatedAt, &order.UpdatedAt); err != nil {
				log.Printf("Error scanning order: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Error loading order data.",
				})
				return
			}
			orders = append(orders, order)
		}

		c.JSON(http.StatusOK, gin.H{
			"title":  "Pending Orders",
			"orders": orders,
		})
	}
}

func GetOrderItems(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		rows, err := db.Query(`
            SELECT id, order_id, coffee_id, quantity, unit_price
            FROM order_items
            WHERE order_id = ?
        `, id)
		if err != nil {
			log.Printf("Error fetching order items: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to load orders. Please try again later.",
			})
			return
		}
		defer rows.Close()
		var order OrderItems
		for rows.Next() {
			if err := rows.Scan(&order.ID, &order.OrderID, &order.CoffeeID, &order.Quantity, &order.UnitPrice); err != nil {
				log.Printf("Error scanning order: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Error loading order data.",
				})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"title":       "Pending Orders",
			"order_items": order,
		})
	}
}
