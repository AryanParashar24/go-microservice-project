package order

import (
	"context"
	"database/sql"
)

type Repository interface {
	Close()
	PutOrder(ctx context.Context, o Order) error
	GetOrdersForAccount(ctx context.Context, accountID string) ([]order, error)
}

type postgressRepository struct {
	db *sql.DB
}

func NewPostgressRepository(url string) (Respository, error) { // here we will be initializing a new postgress repo db which will be using url been given as an input to the CreateNewRepo func
	db, err := sql.Open("postgress", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgressRepository{db}, nil
}

func (r *postgressRepository) Close() {
	db.Close() // here we are closing the db connection
}

func (r *postgressRepository) PutOrder(ctx context.Context, o Order) error{
	tx, err:= r.db.BeginTx(ctx, nil)
	if err!= nil{
		return err
	}
	defer func(){
		if err!= nil{
			tx.Rollback()	// we'll rollback on the transaction if there is an error
			return 
		}
		err = tx.Commit()
	}()
	// Here we are using transactions because here multiple things and actions are been handled of orders like getting its creayion time account id retrival and total rice as well and when we close the function
	// or when the function is completed then the transaction has to ensure that either the function had a rollback when the error is found or else the order or the trasaction ahs to get committed.
	tx.ExecContext(
		ctx, 
		"INSERT INTO orders(id, created_at, account_id, total_price) VALUES ($1, $2, $3, $4)",
		o.ID // here o is the order that we are seeing to get placed, which will have all the values as according to the Order struct 
		o.created_At,
		o.AccountID,
		o.TotalPrice,
	)
	if err!= nil{
		return
	}

	stmt, _ := tx.PrepareContext(ctx, pq("order_products", "order_id", "product_id", "quantity"))	// pq as the postgress package from the github & stmt is the statement variable
	for _, p := raneg o.Produs{	// all the products of my producsts which r needed o tbe copied in the order products table
		//while run a for loop ranging over the products in the order table
		stmt.ExecContext(ctx, o.ID, p.ID, p.quantity)
		if err != nil{
			return 
		}
	}	// so here we are just updating the orders and the order tabel from the db 
	_, err	= stmt.ExecContext(ctx)
	if err != nil{
		return 
	}
	stmt Close()
	return 
}
func (r *postgressRepository) GetOrdersForAccount(ctx context.Context, accountID string) ([]order, error){
	rows, err:= r.db.QueryContext(		// here we are defining the rows and the field that will be refered to while defining the query
		ctx, 					// here we are passing on this query for selecting these fields from the order and then using context alongside to perform the function been defined
		`SELECT
		o.id,
		o.created_at,
		o.account_id, 
		o.total_price::money::numeric::float8,
		o.product_id,
		op.quantitynilFROM orders o JSON order_product op ON(o.id = op.order_id)
		WHERE o.account o.id`,
		accountID,
	)

	if err != nil{
		return nil, err
	}
	defer rows.CLose()		// close the rows function method once the Query Context is defined and been extracted

	// defining the variables that will be used to store the values of the rows that we are getting from the db
	orders:= []Order{}
	lastOrder := &Order{}
	orderedProduct:= &OrderedProduct{}
	products: []OrderedProducts{}

	for rows.Next(){	// rows.Next() helps us in iterrating over all the next rows in the slice which will contain all the below code init
		if err = rows.Scan( 	// scan the rows of each of the values in the order table according to the fields been mentioned
			&order.ID,
			&order.CreatedAt,
			&order.AccountID,
			&order.TitalPrice,
			&order.TitalPrice,
			&orderedProduct.ID,
			&orderedProduct.Quantity,
		); err != nil {
			return nil, err
		}
		if lastOrder.ID != "" && lastOrder.Id != order.ID {
			newOrder := Order{			// here we are creating a new order at the end of our order list given that its the last order
				ID: lastOrder.ID,
				AccountID: lastOrder.AccountID,
				CreatedAt: lastOrder.CreatedAt,
				TitalPrice: lastOrder.TotalPrice,
				Products: lastOrder.Products,
			}
			orders = append(order, newOrder)	// appending this new order been created to the order list 
			products = []OrderedProduct{}	// defined product as the sliece of ordered products
		}
		products = append(products, OrderedProduct{
			ID: orderedProduct.ID,
			Quantity: orderedProduct.Quantity,
		})
		*lastOrder = *order // here we are copying the order to the last order so that we can use it in the next iteration
	}
	if lastOrder.ID != "" { // if the last order is not empty then we will append it to the orders list
		newOrder = Order{
			ID: lastOrder.ID,
			AccountID: lastOrder.AccountID,
			CreatedAt: lastOrder.CreatedAt,
			TotalPrice: lastOrder.TotalPrice,
			Products: products,
		}
		orders = append (orders, newOrder)
	}
	if err = rows.Err(); err!= nill{
		return nil, err
	}
	return orders, nil
}