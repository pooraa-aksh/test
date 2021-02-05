package services

import (
	"github.com/gin-gonic/gin"
	"github.com/test/services/account"
)

//InitServices : Initialise All Services
func InitServices(r *gin.Engine) {
	account.InitServices(r)
}
