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

func Createcourse(c *fiber.Ctx) error {
	type courseinput struct {
		Title         string `json:"title" gorm:"unique" `
		Description   string `json:"description"`
		Enrollmentkey string `json:"enrollmentkey"`
	}

	input := new(courseinput)
	if err := c.BodyParser(input); err != nil {
		return c.JSON(fiber.Map{
			"error": true,
			"msg":   "incorrect input",
		})
	}
	VerifiedTeacher := c.Cookies("username")

	// role := db.DB.Model("user").Column(&ct)
	// var Id uint
	d := new(models.Courses)
	u := new(models.User)
	if res := db.DB.Where(
		&models.User{Username: VerifiedTeacher},
	).First(&u); res.RowsAffected <= 0 {
		return c.JSON(fiber.Map{"error": true, "url": "", "msg": "Invalid Credentials."})
	}

	if title := db.DB.Where(
		&models.Courses{Teacher: VerifiedTeacher}).Where(&models.Courses{Title: input.Title},).First(&d); title.RowsAffected > 0 {
		return c.JSON(fiber.Map{"error": true, "msg": "Course Title already exist."})
	}

	// Userid := row.Scan(&id)

	// fmt.Println(Userid)

	item := models.Courses{

		Title:         input.Title,
		Description:   input.Description,
		Enrollmentkey: input.Enrollmentkey,
		Teacher:       VerifiedTeacher,
		UserID:        u.ID,
	}

	fmt.Println(item)

	if err := db.DB.Create(&item).Error; err != nil {
		return c.JSON(fiber.Map{
			"error": true,
			"msg":   "insertion error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":  false,
		"status": "course created successfully",
	})
}

func DeleteCourse(c *fiber.Ctx) error {

	type iteminput struct {
		Id uint `json:"id"`
	}

	input := new(iteminput)
	if err := c.BodyParser(input); err != nil {
		return c.JSON(fiber.Map{
			"error":  true,
			"status": "incorrect input",
		})
	}
	u := new(models.Courses)
	VerifiedTeacher := c.Cookies("username")
	if res := db.DB.Where(
		&models.Courses{ID: input.Id},
	).Where(&models.Courses{Teacher: VerifiedTeacher}).First(&u); res.RowsAffected <= 0 {
		return c.JSON(fiber.Map{"error": true, "msg": "Invalid request."})
	}

	if err := db.DB.Delete(&models.Courses{}, input.Id).Error; err != nil {
		return c.JSON(fiber.Map{
			"error":  true,
			"status": "Deletion error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":  false,
		"status": "course deleted successfully",
	})
}

func GetYourCourses(c *fiber.Ctx) error {
	type Coursedata struct {
		ID            uint   `json:"id"`
		Title         string `json:"title"`
		Description   string `json:"description"`
		Enrollmentkey string `json:"enrollmentkey"`
	}

	var courses []Coursedata
	var course Coursedata
	var cart models.Courses

	VerifiedTeacher := c.Cookies("username")

	rows, _ := db.DB.Model(&models.Courses{}).Where("teacher = ?", VerifiedTeacher).Rows()
	defer rows.Close()
	for rows.Next() {
		db.DB.ScanRows(rows, &cart)

		course.ID = cart.ID
		course.Title = cart.Title
		course.Description = cart.Description
		course.Enrollmentkey = cart.Enrollmentkey
		courses = append(courses, course)
	}
	fmt.Println(courses)

	return c.JSON(courses)
}
