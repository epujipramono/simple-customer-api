package handler

import (
	"customer-api/internal/model"
	"customer-api/internal/usecase"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CustomerAPIHandler struct {
	usecase *usecase.CustomerAPIUsecase
}

func NewCustomerAPIHandler(uc *usecase.CustomerAPIUsecase) *CustomerAPIHandler {
	return &CustomerAPIHandler{usecase: uc}
}

func (h *CustomerAPIHandler) GetAllCustomers(c *gin.Context) {
	data, err := h.usecase.GetAllCustomers(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *CustomerAPIHandler) GetCustomerByID(c *gin.Context) {
	var req model.Customer
	if err := c.ShouldBindQuery(&req); err != nil {
		respondError(c, err)
		return
	}
	req.ID, _ = strconv.ParseInt(c.Param("id"), 10, 64)
	data, err := h.usecase.GetCustomerByID(c.Request.Context(), req.ID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *CustomerAPIHandler) AddNewCustomer(c *gin.Context) {
	var req model.Customer
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidRequest(c)
		return
	}
	req.ID, _ = strconv.ParseInt(c.Param("id"), 10, 64)

	//validate input
	if err := validateInput(req); err != nil {
		respondError(c, err)
		return
	}
	customer := &model.Customer{Name: req.Name, Email: req.Email, Phone: req.Phone}
	err := h.usecase.AddNewCustomer(c.Request.Context(), customer)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Customer added successfully."})
}

func (h *CustomerAPIHandler) UpdateCustomerByID(c *gin.Context) {
	var req model.Customer
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidRequest(c)
		return
	}
	req.ID, _ = strconv.ParseInt(c.Param("id"), 10, 64)

	if err := validateInput(req); err != nil {
		respondError(c, err)
		return
	}

	customer := &model.Customer{ID: req.ID, Name: req.Name, Email: req.Email, Phone: req.Phone}
	err := h.usecase.UpdateCustomerByID(c.Request.Context(), customer)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Customer updated successfully"})
}

func (h *CustomerAPIHandler) DeleteCustomerByID(c *gin.Context) {
	var req model.Customer
	if err := c.ShouldBindQuery(&req); err != nil {
		respondError(c, err)
		return
	}
	req.ID, _ = strconv.ParseInt(c.Param("id"), 10, 64)
	fmt.Println(req.ID)
	err := h.usecase.DeleteCustomerByID(c.Request.Context(), req.ID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

func validateInput(req model.Customer) error {
	if req.Name == "" || req.Email == "" || req.Phone == "" {
		return errors.New("Name or email or phone is empty")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
	if emailRegex.MatchString(req.Email) == false {
		return errors.New("Incorrect email format")
	}
	phoneRegex := regexp.MustCompile(`^[0-9]+$`)
	if phoneRegex.MatchString(req.Phone) == false {
		return errors.New("Incorrect phone format")
	}
	return nil
}

func invalidRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
}

func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errors.New("Name or email or phone is empty")):
		c.JSON(http.StatusBadRequest, gin.H{"message": "Name, email, and phone are required."})
	case errors.Is(err, errors.New("Incorrect email format")):
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	case errors.Is(err, errors.New("Incorrect phone format")):
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	case errors.Is(err, errors.New("Customer not found")):
		c.JSON(http.StatusNotFound, gin.H{"message": "Customer data not found."})
	case errors.Is(err, errors.New("Duplicate email")):
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email already exists."})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
	}
}
