package database

import (
	"context"
	"database/sql"

	"github.com/LucianTavares/comunicacao_entre_sistemas/graphql/graph/model"
	"github.com/google/uuid"
)

type Category struct {
    db *sql.DB
    ID string
    Name string
    Description string
}

type MemoryStorage struct {
	categories map[string]*model.Category
	courses map[string]*model.Course
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		categories: make(map[string]*model.Category),
		courses: make(map[string]*model.Course),
	}
}

func (m *MemoryStorage) GetCategories(ctx context.Context, categoryIDs []string) ([]*model.Category) {
    output := make([]*model.Category, 0, len(categoryIDs))
	for _, id := range categoryIDs {
        if ctg, ok := m.categories[id]; ok {
            output = append(output, ctg)
		}
	}
	return output
}

func NewCategory(db *sql.DB) *Category {
    return &Category{db: db}
}

func (c *Category) Create(name string, description string) (Category, error) {
    id := uuid.New().String()
    _, err := c.db.Exec("INSERT INTO categories (id, name, description) VALUES ($1, $2, $3)", 
                id, name, description)
    if err != nil {
        return Category{}, err
    }
    return Category{ID: id, Name: name, Description: description}, nil
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
        categories = append(categories, Category{ID: id, Name: name, Description: description})
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
    return Category{ID: id, Name: name, Description: description}, nil
}
