package api

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/gofiber/fiber/v2"
)

type NewUserData struct {
	Username string `json:"username"`
	Mail     string `json:"email"`
	Password string `json:"password"`
}

type Login struct {
	Loginid  string `json:"loginid"`
	Password string `json:"password"`
}

type Message struct {
	Message string
	Key     int
	Value   any
}

func LoginHandler(c *fiber.Ctx) error {
	var login Login
	if err := c.BodyParser(&login); err != nil {
		fmt.Println(err)
		return err
	}
	token, userid, err := LoginUser(login.Loginid, login.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error logging in!!",
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"code":   200,
		"token":  token,
		"userid": userid,
	})
}

func RegisterHandler(c *fiber.Ctx) error {
	var userData NewUserData
	if err := c.BodyParser(&userData); err != nil {
		fmt.Println(err)
		return err
	}
	doesExists := DoesExists("email", userData.Mail)
	if !doesExists {
		return c.JSON(fiber.Map{
			"message": "Email already exists!!",
		})
	}
	doesExists = DoesExists("username", userData.Username)
	if !doesExists {
		return c.JSON(fiber.Map{
			"message": "Username already exists!!",
		})
	}
	err := RegisterUser(userData.Username, userData.Password, userData.Mail)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"code":    500,
			"message": err,
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"code":    200,
		"message": "User registered successfully!!",
	})
}

func FileHandler(c *fiber.Ctx) error {
	// verify token and userid
	token := c.Get("token")
	userid := c.Get("userid")
	isValid, err := Verify(token, userid)
	if err != nil {
		return err
	}
	if !isValid {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to open file")
	}
	defer src.Close()

	// Read the file content into a byte slice
	fileContent, err := io.ReadAll(src)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read file content")
	}
	SaveFile(fileContent, file.Filename, userid)

	now := time.Now()
	file_name := now.Format("20060102150405") + "_" + file.Filename

	if err := c.SaveFile(file, "./src/api/uploads/"+file_name); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func SaveHandler(c *fiber.Ctx) error {
	// verify token and userid
	token := c.Get("token")
	userid := c.Get("userid")
	isValid, err := Verify(token, userid)
	if err != nil {
		return err
	}
	if !isValid {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	var data []interface{}
	if err := c.BodyParser(&data); err != nil {
		fmt.Println(err)
		return err
	}

	var buffer bytes.Buffer
	for _, v := range data {
		buffer.WriteString(fmt.Sprintf("%v", v))
	}

	byteData := buffer.Bytes()

	// Save the file to the server
	file_name := time.Now().Format("20060102150405")
	err = SaveFile(byteData, file_name, userid)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(fiber.Map{
		"message": "File saved successfully",
	})
}
