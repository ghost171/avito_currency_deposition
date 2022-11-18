package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/Serj1c/user-balance/pkg/users"
	"github.com/Serj1c/user-balance/pkg/util"
)

func to_json() {}

type UserHandler struct {
	r *users.Repo
}

func NewUserHandler(r *users.Repo) *UserHandler {
	return &UserHandler{r}
}