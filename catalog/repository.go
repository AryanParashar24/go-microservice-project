package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	ErrNotFound = errors.New("Entity not found")
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, p Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type elasticRepository struct {
	client *elastic.Client
}

type productDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false), // Disable sniffing which is used to discover other nodes in the cluster
	)
	if err != nil {
		return nil, err
	}
	return &elasticRepository{client}, nil
}

func (r *elasticRepository) Close() {
}

func (r *elasticRepository) PutProduct(ctx context.Context, p Product) error { // so we are usign to PutProduct to put our new product
	_, err := r.client.Index(). // here we are indexing the product to the elasticsearch
					Index("catalog"). // we have index as our catalog
					Type("product").  // type has been set to product
					Id(p.ID).
					BodyJson(productDocument{
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		}).Do(ctx) // we do context return error
	return err
}
func (r *elasticRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	res, err := r.client.Get().
		Index("catalog").
		Type("product").
		Id(id).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	if !res.Found { // if no respose has found then ErrorNotFound otherwise return nil
		return nil, ErrNotFound
	}

	p := productDocument{}
	if err = json.Unmarshal(*res.Source, &p); err != nil { // here we will unmarshal which means converting json to productDocument:p and here we will get error if it is not in correct format
		return nil, err
	}
	return &Product{ // Here we are returning the Product along with its structures to the product field
		ID:          id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}, err
}

func (r *elasticRepository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	res, err := r.client.Search(). // here the products are seached and listed
					Index("catalog").
					Type("product").
					Query(elastic.NewMatchAllQuery()). // here its function is to match all the terms or the products with the query that we sent
					From(int(skip)).Size(int(take)).
					Do(ctx) // we do context retrun error
	if err != nil {
		log.Println(err)
		return nil, ErrNotFound
	}
	products := []Product{}
	for _, hit := range res.Hits.Hits { // Hits are all the values that got matched and here we used Hits twice because the first one is the response and the second one is the collection of all the hits that we got after searching
		// so we are getting all the hits and then we are unmarshalling them to get the product document
		// and then we are appending them to the products slice
		// here we need to range over them and then send the collection of them with a specified document which can be set through the struct
		p := Product{}
		if err = json.Unmarshal(*hit.Source, &p); err == nil {
			products = append(products, Product{
				ID:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return products, err
}

func (r *elasticRepository) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	items := []*elastic.MultiGetItem{}
	for _, id := range ids {
		items = append(items,
			elastic.NewMultiGetItem().
				Index("catalog"). // here we are setting the index as catalogIndex("catalog").
				Type("product").
				Id(id),
		)
	}
	res, err := r.client.MultiGet().
		Add(items...).
		Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	products := []Product{}
	for _, doc := range res.Docs {
		p := productDocument{}
		if err := json.Unmarshal(*doc.Source, &p); err == nil {
			products = append(products, Product{
				ID:          doc.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return products, nil
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error) {
	res, err := r.client.Search().
		Index("catalog").
		Type("product").
		Query(elastic.NewMultiMatchQuery(query, "name", "description")).
		From(int(skip)).Size(int(take)).
		Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	products := []Product{}
	for _, hit := range res.Hits.Hits { // here we are ranging over the hits that we got after searching for any particular product from the catalog
		p := productDocument{} // it builds up the documnet of the products after hitting the search
		if err = json.Unmarshal(*hit.Source, &p); err == nil {
			products = append(products, Product{ // heres the docuemnt which is appended to the products document
				ID:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return products, nil
}
