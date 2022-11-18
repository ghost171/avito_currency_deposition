package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"io"
	"net/http"
	"strconv"
	"time"

	"https://github.com/ghost171/avito_currency_deposition/tree/main/pkg/users"
)

type UserHandler struct {
	response *users.Repo
}

func NewUserHandler(response *users.Repo) *UserHandler {
	return &UserHandler{response}
}

func dateFormChange() string {
	time_now := time.Now().Local().AddDate(0, 0, -1)
	formatted := time_now.Format("2022-11-17")
	return formatted
}

func callAPIExchange(currency string) float64 {
	date := dateFormChange()
	
	respond, err := http.Get("http://api.exchangeratesapi.io/v1/" + date + "?access_key=36d585d941651b79dd7d412d57dc66ff&base=EUR&symbols=USD," + currency)
	if err != nil {
		return -1
	}

	defer respond.Body.Close()

	data, err := ioutil.ReadAll(respond.Body)
	result := make(map[string]interface{})
	json.Unmarshal([]byte(data), &result)

	rates := result["rates"].(map[string]interface{})
	curString := fmt.Sprint(rates[currency])
	usdString := fmt.Sprint(rates["USD"])

	cur, _ := strconv.ParseFloat(curString, 64)
	usd, _ := strconv.ParseFloat(usdString, 64)

	exchange_rate  := cur / usd
	return exchange_rate
}

func (uh *UserHandler) GetValue(rw http.ResponseWriter, r *http.Request) {
	parameters := r.URL.Query()
	user_id = parameters["user"][0]
	value, err := uh.response.Value(user_id)
	
	if err == users.ErrorNotExistedUser {
		http.Error(rw, "There is no such user", http.StatusBadRequest)
	} else {
		if _, ok := parameters["currency"]; ok 
		{
			currency := parameters["currency"][0]
			value = value * callAPIExchange(currency)
		}
		respond, err := json.Marshal(value)
		if err != nil {
			http.Error(rw, "Cannot marshal respond", http.StatusInternalServerError)
		}
		rw.Write(respond)
	}

	
}

func (uh *UserHandler)  Deposit(rw http.ResponseWriter, r *Request) {
	parameters := r.URL.Query()
	user_id := parameters["user"][0]
	value, err := strconv.ParseFloat(parameters["value"][0], 64)
	if err != nil {
		http.Error(rw, "Cannot parse value of money", http.StatusBadRequest)
	} 
	if value > 0 {
		err = uh.r.Deposit(user_id, value)
		switch err {
			case nil:
			case users.ErrNoRows:
				http.Error(rw, "User does not exist", http.StatusBadRequest)
			case user.ErrorNotEnoughMoney:
				http.Error(rw, "User does not have enough money", http.StatusBadRequest)
			case users.ErrorDatabase:
				http.Error(rw, "Internal Error", http.StatusInternalServerError)
		}
	} else {
		http.Error(rw, "Withdrawal of only positive sums is allowed", http.StatusBadRequest)
	}
}

func (r *Repo) Cashout(user_id string, value float64) error {

	parameters := r.URL.Query()
	user_id := parameters["users"][0]

	value, err := strconv.ParseFloat(parameters["value"][0], 64)
	if err != nil {
		http.Error(rw, "Cannot parse value of money in user account", http.StatusBadRequest)
	}
	if value > 0 {
		err = uh.r.Cashout(user_id, value)
		switch err {
			case nil:
			case users.ErrNoRows:
				http.Error(rw, "User does not exist", http.StatusBadRequest)
			case user.ErrorNotEnoughMoney:
				http.Error(rw, "User does not have enough money", http.StatusBadRequest)
			case users.ErrorDatabase:
				http.Error(rw, "Internal Error", http.StatusInternalServerError)
		}
	} else {
		http.Error(rw, "Withdrawal of only positive sums is allowed", http.StatusBadRequest)
	}
}

func (uh *UserHandler) Transfer(rw http.ResponseWriter, r *http.Request) {
	parameters := r.URL.Query()
	from_user_id := parameters["from_user"][0]
	to_user_id := parameters["to_user"][0]
	value, err := strconv.ParseFloat(parameters["value"][0], 64)

	if err != nil {
		http.Error(rw, "Cannot parse value of money", http.StatusBadRequest)
	}
	if value > 0 {
		err = uh.r.Transfer(from_user_id, to_user_id, value)
		switch err {
		case nil:
		case users.ErrorNotEnoughMoney:
			http.Error(rw, "User does not have enough money", http.StatusBadRequest)
		case users.ErrNoUser:
			http.Error(rw, "One or both users do not exist", http.StatusBadRequest)
		case users.ErrDBQuery:
			http.Error(rw, "Internal error", http.StatusInternalServerError)
		}
	} else {
		http.Error(rw, "Transfer of only positive sums is allowed", http.StatusBadRequest)
	}
}

func to_json(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

func (uh *UserHandler) ListOperations(rw http.ResponseWriter, r *http.Request) {
	parameters := r.URL.Query()
	user_id := parameters["user"][0]
	sort_by := "date_of_creation"
	sort_order := "asc"
	page := 1
	perPage := 10

	if _, ok := parameters["sort"]; ok {
		sort_by = parameters["sort"][0]
		sort_order = parameters["sort"][1]
	}

	if _, ok := parameters["page"]; ok {
		page, _ = strconv.Atoi(parameters["page"][0])
	}

	if _, ok := parameters["per_page"]; ok {
		per_page, _ = strconv.Atoi(parameters["per_page"][0])
	}

	offset := (page - 1) * per_page
	operations, err := uh.r.List(user_id, sort_by, sort_order, per_page, offset)

	if err != nil {
		http.Error(rw, "Database error", http.StatusInternalServerError)
	}

	err = to_json(operations, rw)
	if err != nil {
		http.Error(rw, "Internal error", http.StatusInternalServerError)
	}
}