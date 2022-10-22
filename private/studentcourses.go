package private

import (
	//"math/rand"

	"fmt"

	db "lms/database"

	"lms/models"

	//"main.go/util"

	//"golang.org/x/crypto/bcrypt"

	//"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	// "github.com/google/uuid"
)

func AvailableCourses(c *fiber.Ctx) error {
	type Coursedata struct {
		ID            uint   `json:"id"`
		Title         string `json:"title"`
		Teacher       string `json:"teacher"`
		Description   string `json:"description"`
		Enrollmentkey string `json:"enrollmentkey"`
	}

	var courses []Coursedata
	var course Coursedata
	var cart models.Courses

	//VerifiedTeacher := c.Cookies("username")

	rows, _ := db.DB.Model(&models.Courses{}).Rows()
	defer rows.Close()
	for rows.Next() {
		db.DB.ScanRows(rows, &cart)

		course.ID = cart.ID
		course.Title = cart.Title
		course.Teacher = cart.Teacher
		course.Description = cart.Description
		course.Enrollmentkey = cart.Enrollmentkey
		courses = append(courses, course)
	}
	fmt.Println(courses)

	return c.JSON(courses)
}
