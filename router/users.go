package router

import (
	"lms/models"
	"time"

	db "lms/database"
	"lms/util"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("key")

// CreateUser route registers a User into the database
func CreateUser(c *fiber.Ctx) error {
	u := new(models.User)

	if err := c.BodyParser(u); err != nil {

		return c.JSON(fiber.Map{
			"error": true,
			"input": "Please review your input",
		})
	}

	// validate if the email, username and password are in correct format
	errors := util.ValidateRegister(u)
	if errors.Err {
		return c.JSON(errors)
	}

	if count := db.DB.Where(&models.User{Email: u.Email}).First(new(models.User)).RowsAffected; count > 0 {
		errors.Err, errors.Email = true, "Email is already registered"
	}
	if count := db.DB.Where(&models.User{Username: u.Username}).First(new(models.User)).RowsAffected; count > 0 {
		errors.Err, errors.Username = true, "Username is already registered"
	}
	if errors.Err {
		return c.JSON(errors)
	}

	// Hashing the password with a random salt
	password := []byte(u.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(
		password,
		8,
	)

	if err != nil {
		panic(err)
	}
	u.Password = string(hashedPassword)

	if err := db.DB.Create(&u).Error; err != nil {
		return c.JSON(fiber.Map{
			"error":   true,
			"general": "Something went wrong, please try again later. 😕",
		})
	}

	//redirect to home
	return c.JSON(errors)
}

// LoginUser route logins a user in the app
func LoginUser(c *fiber.Ctx) error {
	type LoginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	input := new(LoginInput)

	if err := c.BodyParser(input); err != nil {
		return c.JSON(fiber.Map{"redirected": false, "url": "", "msg": "Please review your input"})
	}

	// if input.Password == googlepassword {
	// 	return c.JSON(fiber.Map{"redirected": false, "url": "", "msg": "Invalid Credentials."})
	// }

	u := new(models.User)
	if res := db.DB.Where(
		&models.User{Email: input.Username}).Or(
		&models.User{Username: input.Username},
	).First(&u); res.RowsAffected <= 0 {
		return c.JSON(fiber.Map{"redirected": false, "url": "", "msg": "Invalid Credentials."})
	}

	//Comparing the Role
	if role := db.DB.Where(
		&models.User{Role: input.Role}).First(&u); role.RowsAffected <= 0 {
		return c.JSON(fiber.Map{"redirected": false, "url": "", "msg": "Invalid Credentials."})
	}

	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
		return c.JSON(fiber.Map{"redirected": false, "url": "", "msg": "Incorrect Password"})
	}

	// setting up the authorization cookies
	accessToken := util.GenerateTokens(u.Username)
	accessCookie := util.GetAuthCookies(accessToken)
	c.Cookie(accessCookie)
	c.Cookie(&fiber.Cookie{
		Name:     "username",
		Value:    u.Username,
		HTTPOnly: true,
		Secure:   true,
	})

	//redirect to private route
	if input.Role == "student" {
		return c.Redirect("/api/user/private/student", 301)
	} else {
		return c.Redirect("/api/user/private/teacher", 301)
	}
}

// this function logs out a user and reset the accesscookie to nil
func Logout(c *fiber.Ctx) error {
	c.ClearCookie()

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "0000",
		Expires:  time.Now().Add(1 * time.Second),
		HTTPOnly: true,
		Secure:   true,
	})

	return nil
}

// GetAccessToken generates and sends a new access token iff there is a valid refresh token
func GetAccessToken(c *fiber.Ctx) error {
	type RefreshToken struct {
		RefreshToken string `json:"refresh_token"`
	}

	reToken := new(RefreshToken)
	if err := c.BodyParser(reToken); err != nil {
		return c.JSON(fiber.Map{"error": true, "input": "Please review your input"})
	}

	refreshToken := reToken.RefreshToken

	refreshClaims := new(models.Claims)
	token, _ := jwt.ParseWithClaims(refreshToken, refreshClaims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if res := db.DB.Where(
		"expires_at = ? AND issued_at = ? AND issuer = ?",
		refreshClaims.ExpiresAt, refreshClaims.IssuedAt, refreshClaims.Issuer,
	).First(&models.Claims{}); res.RowsAffected <= 0 {
		// no such refresh token exist in the database
		c.ClearCookie("access_token", "refresh_token")
		return c.SendStatus(fiber.StatusForbidden)
	}

	if token.Valid {
		if refreshClaims.ExpiresAt < time.Now().Unix() {
			// refresh token is expired
			c.ClearCookie("access_token", "refresh_token")
			return c.SendStatus(fiber.StatusForbidden)
		}
	} else {
		// malformed refresh token
		c.ClearCookie("access_token", "refresh_token")
		return c.SendStatus(fiber.StatusForbidden)
	}

	_, accessToken := util.GenerateAccessClaims(refreshClaims.Issuer)

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})

	return c.JSON(fiber.Map{"access_token": accessToken})
}

// GetUsername returns the username of the user signed in
func GetUsername(c *fiber.Ctx) error {
	u := c.Cookies("username")
	return c.JSON(u)
}
