package database

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Course struct {
	db          *sql.DB
	ID          string
	Name        string
	Description *string
	CategoryID  string
	Category    Category
}

func NewCourse(db *sql.DB) *Course {
	return &Course{db: db}
}

func (c *Course) Create(name, description, categoryID string) (*Course, error) {
	id := uuid.New().String()
	_, err := c.db.Exec("INSERT INTO courses (id, name, description, category_id) VALUES ($1, $2, $3, $4)",
		id, name, description, categoryID)
	if err != nil {
		return nil, err
	}
	return &Course{
		ID:          id,
		Name:        name,
		Description: &description,
		CategoryID:  categoryID,
	}, nil
}

func (c *Course) FindAll() ([]Course, error) {
	rows, err := c.db.Query("SELECT id, name, description, category_id FROM courses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	courses := []Course{}
	for rows.Next() {
		var id, name, categoryID string
		var description sql.NullString
		if err := rows.Scan(&id, &name, &description, &categoryID); err != nil {
			return nil, err
		}
		var descriptionPtr *string

		if description.Valid {
			descriptionPtr = &description.String
		} else {
			descriptionPtr = nil
		}
		courses = append(courses, Course{ID: id, Name: name, Description: descriptionPtr, CategoryID: categoryID})
	}
	return courses, nil
}

func (c *Course) FindByCategoryID(categoryID string) ([]Course, error) {
	rows, err := c.db.Query("SELECT id, name, description, category_id FROM courses WHERE category_id = $1", categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	courses := []Course{}
	for rows.Next() {
		var id, name, categoryID string
		var description sql.NullString
		if err := rows.Scan(&id, &name, &description, &categoryID); err != nil {
			return nil, err
		}
		var descriptionPtr *string

		if description.Valid {
			descriptionPtr = &description.String
		} else {
			descriptionPtr = nil
		}
		courses = append(courses, Course{ID: id, Name: name, Description: descriptionPtr, CategoryID: categoryID})
	}
	return courses, nil
}

func (c *Course) FindCategoryInCourse(categoryID string) ([]Course, error) {
	rowsCourse, err := c.db.Query("SELECT * FROM courses")
	if err != nil {
		return nil, err
	}
	defer rowsCourse.Close()
	courses := []Course{}
	categoryIDs := []string{}
	for rowsCourse.Next() {
		var id, name, description, categoryID string
		if err := rowsCourse.Scan(&id, &name, &description, &categoryID); err != nil {
			return nil, err
		}
		courses = append(courses, Course{ID: id, Name: name, Description: &description, CategoryID: categoryID})

		for _, c := range courses {
			if c.CategoryID != categoryID {
				categoryIDs = append(categoryIDs, categoryID)
			}
		}
	}

	rowsCategory, err := c.db.Query("SELECT * FROM categories WHERE categories.id IN (?)", categoryID)
	if err != nil {
		fmt.Println(err)
	}
	defer rowsCategory.Close()
	categories := []Category{}
	for _, course := range courses {
		for rowsCategory.Next() {
			var id, name, description string
			if err := rowsCategory.Scan(&id, &name, &description); err != nil {
				return nil, err
			}
			if id == course.CategoryID {
				categories = append(categories, Category{ID: id, Name: name, Description: &description})
				courses = append(courses, Course{ID: id, Name: name, Description: &description, Category: Category{ID: id, Name: name, Description: &description}})
			}
		}
	}

	return courses, nil
}
