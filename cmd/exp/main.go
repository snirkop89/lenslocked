package main

import (
	"context"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type ctxKey string

const (
	favoriteColorKey ctxKey = "favorite-color"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, favoriteColorKey, "blue")

	value := ctx.Value(favoriteColorKey)
	strValue := value.(string)

	fmt.Println(value)
	fmt.Println(strings.HasPrefix(strValue, "b"))
}
