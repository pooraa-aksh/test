package util

import (
	"github.com/go-pg/pg"
	"github.com/test/lib"
	"github.com/test/services/account/model"
)

//PAccountActivity : process account activity
func PAccountActivity(reqPayload model.AccountActivityReq) (resp interface{}, err error) {
	var tx *pg.Tx
	if tx, err = lib.BeginTx(); err == nil {
		if err = insertUserActivity(tx, reqPayload); err == nil {
			amt := reqPayload.Payload.Amount
			if reqPayload.ActivityType == "debit" {
				err = updateCreditTransactionStatus(tx, reqPayload.Payload.UserID, amt)
				amt = amt * -1
			}
			err = updateUserBalance(tx, reqPayload.Payload.UserID, amt)
		}
		lib.CompleteTransaction(tx, err)
	}
	return
}

//PUpdateExpiredCredit : process update expired credit
func PUpdateExpiredCredit() (resp interface{}, err error) {
	var (
		userIDs []int
		tx      *pg.Tx
		amt     float64
	)
	if tx, err = lib.BeginTx(); err == nil {
		if userIDs, err = getUserIDs(tx); err == nil && len(userIDs) > 0 {
			for _, curID := range userIDs {
				if amt, err = updateStatusAndRemainingAmt(tx, curID); err == nil {
					err = updateUserBalance(tx, curID, amt)
				}
				if err != nil {
					break
				}
			}
		}
		lib.CompleteTransaction(tx, err)
	}
	return
}

//PGetAccountCreditLog : process get account credit log
func PGetAccountCreditLog(userID int) (resp model.CreditLogList, err error) {
	var conn *pg.DB
	selectQry, qryParam := getAccountLogQry(userID)

	if conn, err = lib.Connect(); err == nil {
		if _, err = conn.Query(&resp.CreditLog, selectQry, qryParam...); err == nil {
			if len(resp.CreditLog) == 0 {
				resp.CreditLog = make([]model.CreditLog, 0)
			}
		}
	}
	return
}
