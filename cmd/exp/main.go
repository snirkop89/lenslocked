package main

import (
	stdctx "context"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/snirkop89/lenslocked/context"
	"github.com/snirkop89/lenslocked/models"
)

type ctxKey string

const (
	favoriteColorKey ctxKey = "favorite-color"
)

func main() {
	ctx := stdctx.Background()
	user := models.User{
		Email: "john@example.com",
	}
	ctx = context.WithUser(ctx, &user)

	u := context.User(ctx)
	fmt.Println(u)
}
