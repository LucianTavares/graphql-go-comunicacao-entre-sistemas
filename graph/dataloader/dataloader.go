package dataloader

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/LucianTavares/comunicacao_entre_sistemas/graphql/graph/model"
	"github.com/LucianTavares/comunicacao_entre_sistemas/graphql/internal/database"
	"github.com/graph-gophers/dataloader"
	gopher_dataloader "github.com/graph-gophers/dataloader"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

type DataLoader struct {
	categoryLoader *dataloader.Loader
}

type categoryBatcher struct {
	db database.Categories
}

func (i *DataLoader) GetCategory(ctx context.Context, categoryID string) (*model.Category, error) {
	
	// fmt.Printf("result = %v", )
	thunk := i.categoryLoader.Load(ctx, gopher_dataloader.StringKey(categoryID))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	return result.(*model.Category), nil
}

func NewDataLoader(db database.Categories) *DataLoader {
	categories := &categoryBatcher{db: db}
	return &DataLoader{
		categoryLoader: dataloader.NewBatchedLoader(categories.get),
	}
}

func Middleware(loader *DataLoader, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCtx := context.WithValue(r.Context(), loadersKey, loader)
		r = r.WithContext(nextCtx)
		next.ServeHTTP(w, r)
	})
}

func For(ctx context.Context) *DataLoader {
	return ctx.Value(loadersKey).(*DataLoader)
}

func (c *categoryBatcher) get(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	fmt.Printf("dataloader.categoryBatcher.get, category: [%s]\n", strings.Join(keys.Keys(), ","))
	fmt.Printf("IDS %v", keys.Keys())

	keyOrder := make(map[string]int, len(keys))

	var categoryIDs []string

	for ix, key := range keys {
		categoryIDs = append(categoryIDs, key.String())
		keyOrder[key.String()] = ix
	}

	fmt.Printf("%v", keys)

	dbRecords := c.db.GetAllCategories(ctx, categoryIDs)
	if dbRecords != nil {
		return []*dataloader.Result{{Data: nil}}
	}

	results := make([]*dataloader.Result, len(keys))

	for _, record := range dbRecords {
		ix, ok := keyOrder[record.ID]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, record.ID)
		}
	}

	for categoryID, ix := range keyOrder {
		err := fmt.Errorf("category not found %s", categoryID)
		results[ix] = &dataloader.Result{Data: nil, Error: err}
	}

	return results
}