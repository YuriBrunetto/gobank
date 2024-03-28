package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Storage interface {
	// CreateAccount creates a new account and returns its ID.
	CreateAccount(*Account) (int, error)

	// DeleteAccount deletes the account with the given ID.
	DeleteAccount(int) error

	// UpdateAccount updates the account information.
	UpdateAccount(*Account) error

	// GetAccounts returns a list of all accounts.
	GetAccounts() ([]*Account, error)

	// GetAccountByID retrieves an account by its ID.
	GetAccountByID(int) (*Account, error)

	// GetAccountByNumber retrieves an account by its number.
	GetAccountByNumber(int) (*Account, error)

	// CreateTransfer creates a transfer from one account to another.
	// fromAccount is the number of the account transferring money.
	// toAccount is the number of the account receiving money.
	// amount is the amount of money to transfer.
	CreateTransfer(fromAccount, toAccount, amount int) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	err := s.CreateAccountTable()
	if err != nil {
		return err
	}

	err = s.CreateTransferTable()
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
		id serial primary key,
		first_name varchar(100),
		last_name varchar(100),
		number serial,
    encrypted_password varchar(100),
		balance serial,
		created_at timestamp
	)`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateTransferTable() error {
	query := `CREATE TABLE IF NOT EXISTS transfer (
    id serial primary key,
    from_account serial,
    to_account serial,
    amount serial,
    created_at timestamp
  )`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) (int, error) {
	query := `INSERT INTO account
  (first_name, last_name, number, encrypted_password, balance, created_at)
  VALUES ($1, $2, $3, $4, $5, $6)
  RETURNING id`

	var id int
	err := s.db.QueryRow(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.EncryptedPassword,
		acc.Balance,
		acc.CreatedAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query("delete from account where id = $1", id)
	return err
}

func (s *PostgresStore) GetAccountByNumber(number int) (*Account, error) {
	rows, err := s.db.Query("select * from account where number = $1", number)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("Account with number [%d] not found", number)
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query("select * from account where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("Account with ID %d not found", id)
}

func (s *PostgresStore) CreateTransfer(fromAccount, toAccount, amount int) error {
	query := `INSERT INTO transfer
  (from_account, to_account, amount, created_at)
  VALUES ($1, $2, $3, $4)`

	_, err := s.db.Query(
		query,
		fromAccount,
		toAccount,
		amount,
		time.Now().UTC(),
	)
	if err != nil {
		return err
	}

	// subtract amount from the sender
	querySubtractAmount := `UPDATE account
  SET balance = balance - $1
  WHERE number = $2
  `
	_, err = s.db.Query(querySubtractAmount, amount, fromAccount)
	if err != nil {
		return err
	}

	// add amount to the receiver
	queryAddAmount := `UPDATE account
  SET balance = balance + $1
  WHERE number = $2`
	_, err = s.db.Query(queryAddAmount, amount, toAccount)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("select * from account")
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	defer rows.Close()

	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.EncryptedPassword,
		&account.Balance,
		&account.CreatedAt)

	if err != nil {
		return nil, err
	}

	return account, err
}
