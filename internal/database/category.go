package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/LucianTavares/comunicacao_entre_sistemas/graphql/graph/model"
	"github.com/google/uuid"
)

type Category struct {
	db          *sql.DB
	ID          string
	Name        string
	Description *string
}

type Categories interface {
	GetAllCategories(ctx context.Context, categoryIDs []string) []*model.Category
}

type MemoryStorage struct {
	db          *sql.DB
	ID          string
	Name        string
	Description string
	categories  map[string]*model.Category
	courses     map[string]*model.Course
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		categories: make(map[string]*model.Category),
		courses:    make(map[string]*model.Course),
	}
}

func NewCategory(db *sql.DB) *Category {
	return &Category{db: db}
}

func (m *MemoryStorage) GetAllCategories(ctx context.Context, categoryIDs []string) []*model.Category {
	fmt.Printf("%v", categoryIDs)

	slice := make([]interface{}, len(categoryIDs))
	for i, v := range categoryIDs {
		slice[i] = v
	}

	rowsCategory, err := m.db.Query("select * from categories")
	if err != nil {
		return nil
	}
	defer rowsCategory.Close()

	output := make([]*model.Category, 0, len(categoryIDs))
	for rowsCategory.Next() {
		var id string
		if err := rowsCategory.Scan(&id); err != nil {
			return nil
		}
		for _, id := range categoryIDs {
			if ctg, ok := m.categories[id]; ok {
				output = append(output, ctg)
			}
		}
	}

	return output

}

func (c *Category) Create(name string, description string) (Category, error) {
	id := uuid.New().String()
	_, err := c.db.Exec("INSERT INTO categories (id, name, description) VALUES ($1, $2, $3)",
		id, name, description)
	if err != nil {
		return Category{}, err
	}
	return Category{ID: id, Name: name, Description: &description}, nil
}

func (c *Category) FindAll() ([]Category, error) {
	rows, err := c.db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	categories := []Category{}
	for rows.Next() {
		var id, name, description string
		if err := rows.Scan(&id, &name, &description); err != nil {
			return nil, err
		}
		categories = append(categories, Category{ID: id, Name: name, Description: &description})
	}
	return categories, nil
}

func (c *Category) FindByCourseID(courseID string) (Category, error) {
	var id, name, description string
	err := c.db.QueryRow("SELECT c.id, c.name, c.description FROM categories c JOIN courses co ON c.id = co.category_id WHERE co.id = $1", courseID).
		Scan(&id, &name, &description)
	if err != nil {
		return Category{}, err
	}
	return Category{ID: id, Name: name, Description: &description}, nil
}

func (c *Category) FindByIds(ids []string) ([]Category, error) {
	slice := make([]interface{}, len(ids))
	for i, v := range ids {
		slice[i] = v
	}

	rows, error := c.db.Query("SELECT id, name, description FROM categories WHERE id IN (?"+strings.Repeat(",?", len(slice)-1)+")", slice...)

	if error != nil {
		return nil, error
	}
	defer rows.Close()
	categories := []Category{}
	for rows.Next() {
		var id, name string
		var description sql.NullString

		if err := rows.Scan(&id, &name, &description); err != nil {
			return nil, err
		}
		var descriptionPtr *string

		if description.Valid {
			descriptionPtr = &description.String
		} else {
			descriptionPtr = nil
		}

		categories = append(categories, Category{ID: id, Name: name, Description: descriptionPtr})
	}
	return categories, nil
}
