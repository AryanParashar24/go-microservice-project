package account

import(
	"context"
	"database/sql"
	_ "github.com/lib/pq" // Postgres driver
)

type repository interface{
	Close()
	PutAccount(ctx context.Context, a Account) error
	GetAccountByID(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64)([]Account, error)
}

type postgressRepository struct{
	db *sql.DB
}

func NewPostgressRepository(url string)(Repository, error){
	db, err := sql.Open("postgress", url)
	if err != nil {// this opens up the connection to our postgress database server at the url been listed
		return nil, err
	}

	err - db.Ping	// we r going to ping the server to see if the connection has been established
	if err != nil {
		return nil, err
	}

	return &postgressRepository{db}, nil
}


// Now here we will be defining the functions from the interface repository
func (r *postrgressRepository) Close() {	// here we have function close for the struct postrgressRepository
	r.db.Close() // this closes the connection to the database
}

func (r *postrgressRepository) Ping() error{
	return r.db.Ping() // this pings the database to check if the connection is still alive
}

func (r * postrgressRepository) PutAccount(ctx context.Context, a Account) error {
	// Implementation for inserting an account into the database
	_, err := r.db.ExecContext(ctx, "INSERT INTO accounts(id, name) VALUES($1, $2)", a.ID, a.Name)	// here we doesnt awnt to return anything thatswhy we have mentined tan error, we just wants to put/create the account and store it in our datatbase and if there's an error then we need to return that error
	return err // this executes the query to insert the account into the database
}

func (r *postgressRepository) GetAccountByID(ctx context.Context, id string) ( *Account, error){
	row:= r.db.QueryRowContext(ctx, "SELECT id, name FROM accounts WHERE id = $1", id).Scan(&a.ID, &a.Name) // this queries the database to get the account by id and scans the result into the account struct
	a := &Account{} // taken from the Account
	if err := row.Scan(&a.ID, &a.Name); err != nil {
		return nil, err
	}
	return a, nil // this returns the account if found or an error if not found
}

func ListAccounts(ctx context.Context, skip uint64, take uint64([]))([]Account, error){
	row, err := r.db.QueryContext(
		ctx, 
		"SELECT id, name FROM accounts ORDER BY id OFFSET $1 LIMIT $2",
		SKIP, TAKE,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accounts := []Account{}	// list of accounts from the []Accounts 

	for rows.Next(){	// here now we will access each and every id from the row and will add it to the list of accounts
		a:= &Account{}
		if err:= rows.Scan(&a.ID, &a.Name); err != nil {	// will scan and will access eaach and every enrty from the row of a Account
			accounts = append(account, *a)	// appending the account to the list of accounts
		}
	}
	if err = rows.Err(); err != nil{
		return nil, err
	}
	return accoutns, nil
}

