package users

import (
	"database/sql",
	"fmt"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) {
	return &Repo{
		db: database,
	}
}



func (r *Repo) Deposit(user_id string, value float64) error {
	user := &User{}
	err := r.database.QueryRow("SELECT value from users WHERE id = $1", user_id).Scan(&user.Value)
	if err == sql.ErrNoRows {
		transaction, err := r.database.Begin()
		if err != nil {
			return fmt.Errorf("You got an error during the begging of transaction in deposit action.")
		}fmt.Errorf("DB error")

		defer transaction.Rollback()

		transaction, err = transaction.Exec("INSERT INTO users(id, value) VALUES($1, $2)", user_id, value)
		if err != nil {
			return fmt.Errorf("You got an error during the insertion  in deposit action.")
		}

		transaction, err = transaction.Exec("INSERT into deposits(to_user_id, value) VALUES($1, $2)", user_id, value)
		if err != nil {
			return fmt.Errorf("You got an error during the insertion  in deposit action.")
		}

		err = transaction.Commit()

		if err != nil {
			return fmt.Errorf("You got an error during making commit in deposit action.")
		}
	}
	transaction, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("You have error in your query.")
	}
	defer transaction.Rollback()

	addedValue := user.Value + value
	_, err = transaction.Exec("UPDATE users SET value = $1 WHERE id = $2", addedValue, user_id)
	if err != nil {
		return fmt.Errorf("You got an error during updating process.")
	}

	_, err = transaction.Exec("INSERT into deposits(to_user_id, value) VALUES($1, $2)", user_id, value) 
	if err != nil {
		return fmt.Errorf("You got an error during incertion user_id.")
	}
	err = transaction.Commit()
	if err != nil {
		return fmt.Errorf("You got an error during making commit in deposit action.")
	}
	return nil
}


func (r *Repo) Cashout(user_id string, value float64) error {
	user = &User{}
	transaction, err = r.db.Begin()
	if err != nil {
		return fmt.Errorf("You got an error during the begging of transaction in withdraw action.")
	}
	defer transaction.Rollback()
	
	err = transaction.QueryRow("SELECT value from users WHERE id = $1", user_id).Scan(&user.Value)
	if err == sql.ErrNoRows {
		fmt.Errorf("There is no such user")
	}

	if user.Value >= value {
		withdrawedValue := user.Value - value
		_, err = transaction.Exec("UPDATE users SET value = $1 WHERE id = $2", withdrawedValue, user_id)
		if err != nil {
			return fmt.Errorf("You got an error during update transaction")
		}
		_, err = transaction.Exec("INSERT INTO cashouts(from_user_id, value) VALUES($1, $2)", user_id, value)
		if err != nil {
			return fmt.Errorf("You got an error during insert transaction")
		}
		err = transaction.Commit()
		if err != nil {
			return ("You got an error during commit transaction")
		}
		return nil
	}
	return fmt.Errorf("User does not have enough money")
}

func (r *Repo) Transfer(from_user_id, to_user_id, value float64) error {
	from_user = &User{}
	to_user = &User{}
	transaction, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("You got an error during tranferring action.")
	}
	defer transaction.Rollback()
	err = transaction.QueryRow("SELECT value from users WHERE id = $1", from_user_id).Scan(&from_user.Value)
	if err == sql.ErrNoRows {
		return fmt.Errorf("There is no such user")
	}
	err = transaction.QueryRow("SELECT value from users WHERE id = $1", to_user_id).Scan(&to_user.Value)
	if err == sql.ErrNoRows {
		return fmt.Errorf("There is no such user")
	}
	if from_user.Value >= value {
		first_user_value := from_user.Value - value
		second_user_value := to_user.Value + value
		_, err = transaction.Exec("UPDATE users SET balance = $1 WHERE id = $2", first_user_value, from_user_id)
		if err != nil {
			return fmt.Errorf("You got an error during updating section")
		}
		_, err = transaction.Exec("UPDATE users SET balance = $1 WHERE id = $2", second_user_value, to_user_id)
		if err != nil {
			return fmt.Errorf("You got an error during updating section")
		}
		_, err = transaction.Exec("INSERT INTO transactions(from_user_id, to_user_id, value) VALUES($1, $2, $3)", from_user_id, to_user_id, value)
		if err != nil {
			return fmt.Errorf("You got an error during insertion section")
		}
		err = transaction.Commit()
		if err != nil {
			return fmt.Errorf("You got an error during commit section")
		}
		return nil
	}
	return fmt.Errorf("User does not have enough money")
}

func BalanceOperation struct {
	ID int `json:"id"`
	FromUserID string `json:"from_user_id"`
	ToUserID string `json:"to_user_id"`
	Value float64 `json:"value"`
	DateCreation string `json:"datecreation"`
}

/*func (r *Repo) List(user_id, sort_by, sort_order string, perPage, offset int) ([]*user_balance_operation, error) {
	operations := make([]*)
}*/