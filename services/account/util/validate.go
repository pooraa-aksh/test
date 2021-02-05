package util

import (
	"errors"
	"strconv"
	"time"

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
		if err == nil {
			if reqPayload.ActivityType == "debit" {
				err = checkCurrentBalance(reqPayload.Payload)
			} else {
				expDate := time.Unix(int64(reqPayload.Payload.Expiry), 0)
				if isAfter := time.Now().After(expDate); isAfter {
					err = errors.New("Invalid Expiry Date.")
				}
			}
		}
	}

	return
}

//VGetAccountCreditLog: validate get account credit log
func VGetAccountCreditLog(req *gin.Context) (userID int, err error) {
	userIDStr := req.DefaultQuery("user_id", "0")
	if userID, err = strconv.Atoi(userIDStr); err == nil {
		if userID <= 0 {
			err = errors.New("Please Provide a Valid User ID.")
		} else {
			err = checkUserExits(userID)
		}
	}
	return
}
