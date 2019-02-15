package customer

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/Chatchaikan/finalexam/database"
	"github.com/gin-gonic/gin"
)

func getCustomersHandler(c *gin.Context) {

	status := c.Query("status")
	if status != "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not support yet"})
		return
	}

	stmt, err := database.Conn().Prepare("SELECT id, name, email, status FROM customers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "can't prepare query all customers statment" + err.Error()})
		return
	}

	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "can't query all customers" + err.Error()})
		return
	}

	var customers = []Customer{}

	for rows.Next() {
		c := Customer{}
		err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Status)
		if err != nil {
			//c.JSON(http.StatusInternalServerError, gin.H{"message": "can't Scan row into variable" + err.Error()})
			return
		}

		customers = append(customers, c)
	}

	c.JSON(http.StatusOK, customers)
}

func createCustomersHandler(c *gin.Context) {
	var customer Customer
	err := c.ShouldBindJSON(&customer)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	row := database.InsertCustomer(customer.Name, customer.Email, customer.Status)
	err = row.Scan(&customer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "can't Scan row into variable" + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

func getCustomerByIDHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	stmt, err := database.Conn().Prepare("SELECT id, name, email, status FROM customers WHERE id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	row := stmt.QueryRow(id)

	customer := Customer{}
	err = row.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "data not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func updateCustomerHandler(c *gin.Context) {
	customer := Customer{}
	err := c.ShouldBindJSON(&customer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := database.Conn().Prepare("UPDATE customers SET status=$4, email=$3, name=$2 WHERE id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	if _, err := stmt.Exec(id, customer.Name, customer.Email, customer.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	customer.ID = id

	c.JSON(http.StatusOK, customer)
}

func deleteCustomerHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	// Customer: handle error

	stmt, err := database.Conn().Prepare("DELETE FROM customers WHERE id = $1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if _, err := stmt.Exec(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}

func loginMiddleware(c *gin.Context) {
	log.Println("starting middleware")
	authKey := c.GetHeader("Authorization")
	if authKey != "token2019" {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		c.Abort()
		return
	}
	c.Next()
}

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/customers", loginMiddleware, getCustomersHandler)
	r.GET("/customers/:id", getCustomerByIDHandler)
	r.PUT("/customers/:id", updateCustomerHandler)
	r.DELETE("/customers/:id", deleteCustomerHandler)
	r.POST("/customers", createCustomersHandler)

	return r
}

var db *sql.DB

func CreateTable() {

	ctb := `
	CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);`

	//_, err := db.Exec(ctb)
	_, err := database.Conn().Exec(ctb)
	if err != nil {
		log.Fatal("can't create table", err)
	}
}
