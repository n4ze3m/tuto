package handlers

import (
	"strings"
	"time"

	"github.com/NazeemNato/tuto/src/database"
	"github.com/NazeemNato/tuto/src/middlewares"
	"github.com/NazeemNato/tuto/src/models"
	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	user := models.User{
		Firstname:    data["first_name"],
		Lastname:     data["last_name"],
		Email:        data["email"],
		IsAmbassador: strings.Contains(c.Path(), "/api/ambassador"),
	}
	// set password
	user.SetPassword(data["password"])

	database.DB.Create(&user)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email=?", data["email"]).First(&user)

	if user.Id == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}

	if !user.ComparePassword(data["password"]) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid password"})
	}

	IsAmbassador := strings.Contains(c.Path(), "/api/ambassador")
	var scope string

	if IsAmbassador {
		scope = "ambassador"
	} else {
		scope = "admin"
	}

	if !IsAmbassador && user.IsAmbassador {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	token, err := middlewares.GenerateJWT(user.Id, scope)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "something went wrong"})
	}

	cookie := fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{"message": "Success"})
}

func User(c *fiber.Ctx) error {
	id, _ := middlewares.GetUserId(c)

	var user models.User
	database.DB.Where("id = ?", id).First(&user)

	if strings.Contains(c.Path(), "/api/ambassador") {
		ambassador := models.Ambassador(user)
		ambassador.CalculateRevenue(database.DB)
		return c.JSON(ambassador)
	}

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{"message": "Success"})
}

func UpdateInfo(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	id, _ := middlewares.GetUserId(c)

	user := models.User{
		Firstname: data["first_name"],
		Lastname:  data["last_name"],
		Email:     data["email"],
	}

	user.Id = id

	database.DB.Model(&user).Updates(&user)

	return c.JSON(fiber.Map{"message": "Profile updated"})
}

func UpdatePassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	id, _ := middlewares.GetUserId(c)

	user := models.User{}

	user.Id = id

	user.SetPassword(data["password"])

	database.DB.Model(&user).Updates(&user)

	return c.JSON(fiber.Map{"message": "Password updated"})
}
