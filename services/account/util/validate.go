package util

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/test/services/account/model"
)

//VAccountActivity: validate account activity
func VAccountActivity(req *gin.Context) (reqPayload model.AccountActivityReq, err error) {
	if err = req.ShouldBindJSON(&reqPayload); err == nil {
		if reqPayload.Payload.UserID <= 0 {
			err = errors.New("Please Provide a Valid User ID.")
		} else if reqPayload.Payload.Amount <= 0.0 {
			err = errors.New("Please Provide a Valid Amount.")
		} else if reqPayload.ActivityType != "credit" && reqPayload.ActivityType != "debit" {
			err = errors.New("Invalid Activity type.")
		}
		if err == nil && reqPayload.ActivityType == "debit" {
			err = checkCurrentBalance(reqPayload.Payload)
		}
	}

	return
}

//VGetAccountCreditLog: validate get account credit log
func VGetAccountCreditLog(req *gin.Context) (userID int, err error) {
	userIDStr := req.DefaultQuery("user_id", "")
	if userID, err = strconv.Atoi(userIDStr); err == nil {
		if userID <= 0 {
			err = errors.New("Please Provide a Valid User ID.")
		} else {
			err = checkUserExits(userID)
		}
	}
	return
}
