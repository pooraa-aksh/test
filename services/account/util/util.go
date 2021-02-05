package util

import (
	"errors"
	"time"

	"github.com/go-pg/pg"
	db "github.com/test/database"
	"github.com/test/lib"
	"github.com/test/services/account/model"
)

//checkUserExits : check user exists in database
func checkUserExits(userID int) (err error) {
	var (
		conn  *pg.DB
		count int
	)
	selectQry := `SELECT COUNT(*) FROM user_account WHERE user_id = ?;`

	if conn, err = lib.Connect(); err == nil {
		if _, err = conn.Query(&count, selectQry, userID); err == nil {
			if count == 0 {
				err = errors.New("Please Provide a Valid User ID.")
			}
		}
	}
	return
}

//checkCurrentBalance : check user current balance
func checkCurrentBalance(payload model.ActivityDetail) (err error) {
	var (
		conn    *pg.DB
		balance float64
	)
	selectQry := `SELECT balance FROM user_account WHERE user_id = ?;`

	if conn, err = lib.Connect(); err == nil {
		if _, err = conn.Query(&balance, selectQry, payload.UserID); err == nil {
			if payload.Amount > balance {
				err = errors.New("Insufficient Balance in User Account.")
			}
		}
	}
	return
}

//insertUserActivity : insert user activity
func insertUserActivity(tx *pg.Tx, reqPayload model.AccountActivityReq) (err error) {
	var accActivity db.UserAccountActivity
	accActivity.FkUserID = reqPayload.Payload.UserID
	accActivity.Amount = reqPayload.Payload.Amount
	accActivity.Priority = reqPayload.Payload.Priority
	accActivity.Type = reqPayload.Payload.Type
	accActivity.TransactionType = reqPayload.ActivityType
	accActivity.BalanceAmt = reqPayload.Payload.Amount
	accActivity.ExpDate = time.Unix(int64(reqPayload.Payload.Expiry), 0)
	err = tx.Insert(&accActivity)
	return
}

//updateCreditTransactionStatus : update credit transaction status
func updateCreditTransactionStatus(tx *pg.Tx, userID int, amtToDebit float64) (err error) {
	var (
		txIDAmountMap map[int]float64
		sortedTxIDs   []int
		status        string
		remainingAmt  float64
	)

	selectQry := `
	SELECT JSONB_AGG(transaction_id ORDER BY priority desc), JSONB_OBJECT_AGG(transaction_id, balance_amt)
	FROM user_account_activity 
	WHERE fk_user_id = ? AND transaction_type = ? AND status = ?;`

	if _, err = tx.Query(pg.Scan(&sortedTxIDs, &txIDAmountMap), selectQry, userID, "credit", "active"); err == nil {
		for _, curTxID := range sortedTxIDs {
			if amt, exists := txIDAmountMap[curTxID]; exists {
				if amt <= amtToDebit {
					amtToDebit -= amt
					status = "debited"
					remainingAmt = 0.0
				} else {
					remainingAmt = amt - amtToDebit
					amtToDebit -= amtToDebit
				}
				err = updateAccountActivity(tx, curTxID, status, remainingAmt)
				if err != nil || amtToDebit == 0.0 {
					break
				}
			} else {
				err = errors.New("Something Went Wrong.")
				break
			}
		}
	}
	return
}

//updateAccountActivity : update account activity
func updateAccountActivity(tx *pg.Tx, txID int, status string, remainingAmt float64) (err error) {
	updateQry := `UPDATE user_account_activity SET balance_amt = ?`
	qryParam := []interface{}{remainingAmt}
	if status != "" {
		updateQry += `, status = ?`
		qryParam = append(qryParam, status)
	}
	updateQry += ` WHERE transaction_id = ?;`
	qryParam = append(qryParam, txID)
	_, err = tx.Exec(updateQry, qryParam...)
	return
}

//updateUserBalance : update user balance
func updateUserBalance(tx *pg.Tx, userID int, amt float64) (err error) {
	upsertQry := `
	INSERT INTO user_account(user_id, balance) VALUES(?, ?)
	ON CONFLICT (user_id) 
	DO UPDATE SET balance = user_account.balance + EXCLUDED.balance;`

	_, err = tx.Exec(upsertQry, userID, amt)
	return
}

//getUserIDs :: get user ids
func getUserIDs(tx *pg.Tx) (userIDs []int, err error) {
	selectQry := `SELECT user_id FROM user_account;`

	_, err = tx.Query(&userIDs, selectQry)
	return
}

//lockAccountActivityRows : lock account activity rows
func lockAccountActivityRows(tx *pg.Tx, userID int) (err error) {
	var lock []int

	selectQry := `SELECT 1 FROM user_account_activity WHERE fk_user_id = ? FOR UPDATE;`
	_, err = tx.Query(&lock, selectQry, userID)
	return
}

//updateStatusAndRemainingAmt : update status and remaining amt
func updateStatusAndRemainingAmt(tx *pg.Tx, userID int) (amt float64, err error) {
	var remainingAmt []float64

	updateQry := `
	UPDATE user_account_activity SET status = ?, balance_amt = ?
	WHERE fk_user_id = ? AND transaction_type = ? AND status = ? AND exp_date <= NOW()
	RETURNING balance_amt;`

	if _, err = tx.Query(&remainingAmt, updateQry, "expired", 0.0, userID, "credit", "active"); err == nil {
		for _, curAmt := range remainingAmt {
			amt += curAmt
		}
		amt = amt * -1
	}
	return
}

//getAccountLogQry : get account log qry
func getAccountLogQry(userID int) (selectQry string, qryParam []interface{}) {
	selectQry = `
	SELECT fk_user_id, type, status, exp_date, amount, priority 
	FROM user_account_activity WHERE fk_user_id = ? AND transaction_type = ?;`

	qryParam = append(qryParam, userID, "credit")
	return
}
