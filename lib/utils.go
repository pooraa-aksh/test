package lib

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	pg "github.com/go-pg/pg"
)

//Resp : Response Model
type Resp struct {
	Data   interface{} `json:"data"`
	Error  string      `json:"error"`
	Status bool        `json:"status"`
}

//GinResponse : creates the response of gin API
func GinResponse(c *gin.Context, data interface{}, err error) {
	out := &Resp{}

	httpCode := http.StatusBadRequest
	out.Data = data
	if err == nil {
		out.Status = true
		httpCode = http.StatusOK
	} else {
		out.Error = err.Error()
	}
	c.JSONP(httpCode, out)
}

//Connect : connect to database
func Connect() (db *pg.DB, err error) {
	opts := &pg.Options{
		User:     "postgres",
		Password: "Bvd^412",
		Addr:     "172.16.24.188:5432",
		Database: "tachyon",
	}

	if db = pg.Connect(opts); db == nil {
		err = errors.New("Failed To Connect To Database.")
	}
	return
}

//BeginTx : begin transaction
func BeginTx() (tx *pg.Tx, err error) {
	var conn *pg.DB
	if conn, err = Connect(); err == nil {
		if tx, err = conn.Begin(); err == nil && tx == nil {
			err = errors.New("Failed To Initiate Transaction.")
		}
	}
	return
}

//CompleteTransaction : Commits/Rollback Transaction Based On Error
func CompleteTransaction(tx *pg.Tx, err error) {
	if tx != nil {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}
}
