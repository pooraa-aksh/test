package account

import (
	"github.com/gin-gonic/gin"
	"github.com/test/lib"
	"github.com/test/services/account/model"
	"github.com/test/services/account/util"
)

//InitServices : Add Services
func InitServices(r *gin.Engine) {
	account := r.Group("account")
	{
		v1 := account.Group("v1")
		{
			v1.POST("/activity", AccountActivity)
			v1.PUT("/expired", UpdateExpiredCredit)
			v1.GET("/credit/log", GetAccountCreditLog)
		}
	}
}

//AccountActivity : account activity
func AccountActivity(c *gin.Context) {
	var (
		reqPayload model.AccountActivityReq
		data       interface{}
		err        error
	)
	if reqPayload, err = util.VAccountActivity(c); err == nil {
		data, err = util.PAccountActivity(reqPayload)
	}
	lib.GinResponse(c, data, err)
}

//UpdateExpiredCredit : update expired credits
//This api will be used via a cron that will be setup to run at a particular time interval
func UpdateExpiredCredit(c *gin.Context) {
	var (
		data interface{}
		err  error
	)
	data, err = util.PUpdateExpiredCredit()
	lib.GinResponse(c, data, err)
}

//GetAccountCreditLog : get account credit log
func GetAccountCreditLog(c *gin.Context) {
	var (
		data   model.CreditLogList
		err    error
		userID int
	)
	if userID, err = util.VGetAccountCreditLog(c); err == nil {
		data, err = util.PGetAccountCreditLog(userID)
	}
	lib.GinResponse(c, data, err)
}
