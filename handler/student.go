package handler

import (
	"net/http"
	"strconv"

	"crud-golang/db"

	"github.com/labstack/echo/v4"
)

type Student struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"not null" json:"name"`
	Email     string `gorm:"uniqueIndex;not null" json:"email"`
	CourseID  uint   `gorm:"not null" json:"course_id"`
	Course    Course `json:"course"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func CreateStudent(c echo.Context) error {
	student := new(Student)
	if err := c.Bind(student); err != nil {
		return err
	}

	if err := db.GetDB().Create(student).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, student)
}

func GetStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	student, err := getStudentByID(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, student)
}

func UpdateStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	student, err := getStudentByID(id)
	if err != nil {
		return err
	}

	if err := c.Bind(student); err != nil {
		return err
	}

	if err := db.GetDB().Save(student).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, student)
}

func DeleteStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	if err := db.GetDB().Delete(&Student{}, id).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "deleted",
	})
}

func GetStudents(c echo.Context) error {
	var students []Student
	result := db.GetDB().Preload("Course").Find(&students)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to get students",
		})
	}
	return c.JSON(http.StatusOK, students)
}

func getStudentByID(id int) (*Student, error) {
	student := new(Student)
	if err := db.GetDB().Preload("Course").First(student, id).Error; err != nil {
		return nil, err
	}

	return student, nil
}
