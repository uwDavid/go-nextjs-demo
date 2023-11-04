package api

import (
	"context"
	"database/sql"
	"net/http"
	db "nextjs/backend/db/sqlc"
	"nextjs/backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Account struct {
	server *Server
}

func (a Account) router(server *Server) {
	a.server = server

	serverGroup := server.router.Group("/account", AuthenticatedMiddleWare())
	serverGroup.POST("create", a.createAccount)
	serverGroup.GET("", a.getUserAccounts)
	serverGroup.POST("transfer", a.transfer)
}

type AccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (a Account) createAccount(ctx *gin.Context) {
	userId, err := utils.GetActiveUser(ctx)
	if err != nil {
		return
	}

	acc := new(AccountRequest)

	if err := ctx.ShouldBindJSON(&acc); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := db.CreateAccountParams{
		UserID:   int32(userId),
		Currency: acc.Currency,
		Balance:  0,
	}

	account, err := a.server.queries.CreateAccount(context.Background(), arg)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Account already have this currency"})
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

func (a *Account) getUserAccounts(c *gin.Context) {
	userId, err := utils.GetActiveUser(c)
	if err != nil {
		return
	}

	accounts, err := a.server.queries.GetAccountByUserID(context.Background(), int32(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

type TransferRequest struct {
	ToAccountID   int32   `json:"to_account_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	FromAccountID int32   `json:"from_account_id" binding:"required"`
}

func (a *Account) transfer(c *gin.Context) {
	userId, err := utils.GetActiveUser(c)
	if err != nil {
		return
	}

	tr := new(TransferRequest)

	if err := c.ShouldBindJSON(&tr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := a.server.queries.GetAccountByID(context.Background(), int64(tr.FromAccountID))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't get account"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if account.UserID != int32(userId) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	toAccount, err := a.server.queries.GetAccountByID(context.Background(), int64(tr.ToAccountID))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't get account"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if toAccount.Currency != account.Currency {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Currencies do not match"})
		return
	}

	if account.Balance < tr.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You don't have enough balance"})
		return
	}

	txArg := db.CreateTransferParams{
		FromAccountID: tr.FromAccountID,
		ToAccountID:   tr.ToAccountID,
		Amount:        tr.Amount,
	}

	tx, err := a.server.queries.TransferTx(context.Background(), txArg)
	if toAccount.Currency != account.Currency {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Encountered issue with transaction."})
		return
	}
	c.JSON(http.StatusCreated, tx)
}
