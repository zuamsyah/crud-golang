package handler

import (
	"net/http"
	"strconv"

	"crud-golang/db"

	"github.com/labstack/echo/v4"
)

type Course struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Students  []Student `gorm:"foreignKey:CourseID" json:"students"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}

func CreateCourse(c echo.Context) error {
	course := new(Course)
	if err := c.Bind(course); err != nil {
		return err
	}

	if err := db.GetDB().Create(course).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, course)
}

func GetCourse(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	course, err := getCourseByID(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, course)
}

func UpdateCourse(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	course, err := getCourseByID(id)
	if err != nil {
		return err
	}

	if err := c.Bind(course); err != nil {
		return err
	}

	if err := db.GetDB().Save(course).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, course)
}

func DeleteCourse(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	if err := db.GetDB().Delete(&Course{}, id).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "deleted",
	})
}

func GetCourses(c echo.Context) error {
	var courses []Course
	result := db.GetDB().Preload("Students").Find(&courses)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to get courses",
		})
	}
	return c.JSON(http.StatusOK, courses)
}

func getCourseByID(id int) (*Course, error) {
	course := new(Course)
	if err := db.GetDB().Preload("Students").First(course, id).Error; err != nil {
		return nil, err
	}

	return course, nil
}
