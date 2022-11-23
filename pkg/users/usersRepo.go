package users

import (
	"database/sql"
	"fmt"
)

type User struct {
	Id string `json:"ID"`
	Value float64 `json:"balance"`
	Currency string `json:"value"`
	DateOfCreation string `json:"date_of_creation"`
}

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}

var (
	ErrorNotExistedUser = fmt.Errorf("There is no such user")
	ErrorNotEnoughMoney = fmt.Errorf("Does not have enough money")
	ErrorDatabase = fmt.Errorf("Database error")
)

func (r *Repo) Deposit(user_id string, value float64) error {
	user := &User{}
	err := r.db.QueryRow("SELECT value from users WHERE id = $1", user_id).Scan(&user.Value)
	if err == sql.ErrNoRows {
		transaction, err := r.db.Begin()
		if err != nil {
			return ErrorDatabase
		}

		defer transaction.Rollback()

		_, err = transaction.Exec("INSERT INTO users(id, value) VALUES($1, $2)", user_id, value)
		if err != nil {
			return ErrorDatabase
		}

		_, err = transaction.Exec("INSERT into deposits(to_user_id, value) VALUES($1, $2)", user_id, value)
		if err != nil {
			return ErrorDatabase
		}

		err = transaction.Commit()

		if err != nil {
			return ErrorDatabase
		}
	}
	transaction, err := r.db.Begin()
	if err != nil {
		return ErrorDatabase
	}
	defer transaction.Rollback()

	addedValue := user.Value + value
	_, err = transaction.Exec("UPDATE users SET value = $1 WHERE id = $2", addedValue, user_id)
	if err != nil {
		return ErrorDatabase
	}

	_, err = transaction.Exec("INSERT into deposits(to_user_id, value) VALUES($1, $2)", user_id, value) 
	if err != nil {
		return ErrorDatabase
	}
	err = transaction.Commit()
	if err != nil {
		return ErrorDatabase
	}
	return nil
}


func (r *Repo) Cashout(user_id string, value float64) error {
	user := &User{}
	transaction, err := r.db.Begin()
	if err != nil {
		return ErrorDatabase
	}
	defer transaction.Rollback()
	
	err = transaction.QueryRow("SELECT value from users WHERE id = $1", user_id).Scan(&user.Value)
	if err == sql.ErrNoRows {
		return ErrorNotExistedUser
	}

	if user.Value >= value {
		withdrawedValue := user.Value - value
		_, err = transaction.Exec("UPDATE users SET value = $1 WHERE id = $2", withdrawedValue, user_id)
		if err != nil {
			return ErrorDatabase
		}
		_, err = transaction.Exec("INSERT INTO cashouts(from_user_id, value) VALUES($1, $2)", user_id, value)
		if err != nil {
			return ErrorDatabase
		} 	
		err = transaction.Commit()
		if err != nil {
			return ErrorDatabase
		}
		return nil
	}
	return ErrorNotEnoughMoney
}

func (r *Repo) Transfer(from_user_id, to_user_id string, value float64) error {
	from_user := &User{}
	to_user := &User{}
	transaction, err := r.db.Begin()
	if err != nil {
		return ErrorDatabase
	}
	defer transaction.Rollback()
	err = transaction.QueryRow("SELECT value from users WHERE id = $1", from_user_id).Scan(&from_user.Value)
	if err == sql.ErrNoRows {
		return ErrorNotExistedUser
	}
	err = transaction.QueryRow("SELECT value from users WHERE id = $1", to_user_id).Scan(&to_user.Value)
	if err == sql.ErrNoRows {
		return ErrorNotExistedUser
	}
	if from_user.Value >= value {
		first_user_value := from_user.Value - value
		second_user_value := to_user.Value + value
		_, err = transaction.Exec("UPDATE users SET balance = $1 WHERE id = $2", first_user_value, from_user_id)
		if err != nil {
			return ErrorDatabase
		}
		_, err = transaction.Exec("UPDATE users SET balance = $1 WHERE id = $2", second_user_value, to_user_id)
		if err != nil {
			return ErrorDatabase
		}
		_, err = transaction.Exec("INSERT INTO transactions(from_user_id, to_user_id, value) VALUES($1, $2, $3)", from_user_id, to_user_id, value)
		if err != nil {
			return ErrorDatabase
		}
		err = transaction.Commit()
		if err != nil {
			return ErrorDatabase
		}
		return nil
	}
	return ErrorDatabase
}

func (r *Repo) Value(user_id string) (float64, error) {
	user := &User{}
	err := r.db.QueryRow("SELECT value FROM users WHERE id = $1", user_id).Scan(&user.Value)

	if err == sql.ErrNoRows {
		return -1, ErrorNotExistedUser
	}
	return user.Value, nil
}

type ValueOperation struct {
	ID int `json:"id"`
	FromUserID string `json:"from_user_id"`
	ToUserID string `json:"to_user_id"`
	Value float64 `json:"value"`
	DateCreation string `json:"datecreation"`
}

func (r *Repo) List(user_id, sort_by, sort_order string, perPage, offset int) ([]*ValueOperation, error) {
	operations := make([]*ValueOperation, 0, 10)
	sort := fmt.Sprintf(" ORDER BY %s %s", sort_by, sort_order)
	limitation := fmt.Sprintf(" LIMIT %d OFFSET %d ", perPage, offset)
	
	rows, err := r.db.Query(`SELECT id, from_user_id, to_user_id, value, date_of_creation FROM deposits
	WHERE to_user_id = $1
	UNION ALL SELECT id, from_user_id, to_user_id, value, date_of_creation FROM cashouts
	WHERE from_user_id = $1
	UNION ALL SELECT id, from_user_id, to_user_id, value, date_of_creation FROM transactions 
	WHERE from_user_id = $1
	UNION ALL SELECT id, from_user_id, to_user_id, value, date_of_creation FROM transactions
	WHERE to_user_id = $1` + sort + limitation, user_id)

	if err != nil {
		fmt.Println(err)
		return nil, ErrorDatabase
	}

	defer rows.Close()

	for rows.Next() {
		item := &ValueOperation{}
		err := rows.Scan(&item.ID, &item.FromUserID, &item.ToUserID, &item.Value, &item.DateCreation)
		if err != nil {
			return nil, ErrorDatabase
		}
		operations = append(operations, item)
	}
	return operations, nil
}