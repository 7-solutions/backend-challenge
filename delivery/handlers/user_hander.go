package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/testGolang/backend-challenge/application/usecases"
	"github.com/testGolang/backend-challenge/domain"
)

type UserHandler struct {
	uc *usecases.UserUseCase
}

func NewUserHandler(uc *usecases.UserUseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

// Register godoc
// @Summary     Register a new user
// @Description Create a user and return the created object
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       registerRequest  body     object  true  "Register payload"
// @Success     201              {object} domain.User
// @Failure     400              {object} map[string]string
// @Router      /users/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {

		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	user, err := h.uc.Register(body.Name, body.Email, body.Password)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(user)
}

// Login godoc
// @Summary     Authenticate user
// @Description Verify credentials and return a JWT token
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       loginRequest  body      domain.LoginRequest  true  "Login payload"
// @Success     200           {object}  map[string]string
// @Failure     400           {object}  map[string]string
// @Failure     401           {object}  map[string]string
// @Router      /users/login [post]
// Login handles user authentication and returns a JWT token
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var body domain.LoginRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	token, err := h.uc.Authenticate(body.Email, body.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"token": token})
}

// List godoc
// @Summary     List all users
// @Description Get a list of all registered users
// @Tags        users
// @Security    BearerAuth
// @Param       Authorization  header  string  true  "Bearer {token}"
// @Produce     json
// @Success     200 {array} domain.User
// @Failure     500 {object} map[string]string
// @Router      /users [get]
func (h *UserHandler) List(c *fiber.Ctx) error {
	users, err := h.uc.List()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

// Get godoc
// @Summary     Get user by ID
// @Description Retrieve a single user by their ID
// @Tags        users
// @Security    BearerAuth
// @Produce     json
// @Param       id   path      string  true  "User ID"
// @Success     200  {object}  domain.User
// @Failure     404  {object}  map[string]string
// @Router      /users/{id} [get]
func (h *UserHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.uc.Get(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

// Update godoc
// @Summary     Update user details
// @Description Update a user’s name or email
// @Tags        users
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       id            path      string  true  "User ID"
// @Param       updateRequest body      domain.UpdateUserRequest  true  "Update payload"
// @Success     200           {object} domain.User
// @Failure     400           {object} map[string]string
// @Router      /users/{id} [put]
func (h *UserHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var body domain.UpdateUserRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	user, err := h.uc.Update(id, body.Name, body.Email, body.Password)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

// Delete godoc
// @Summary     Delete a user
// @Description Remove a user by their ID
// @Tags        users
// @Security    BearerAuth
// @Param       id   path      string  true  "User ID"
// @Success     204  "No Content"
// @Failure     400  {object} map[string]string
// @Router      /users/{id} [delete]
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.uc.Delete(id); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
