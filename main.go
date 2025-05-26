package main

import (
	"context"
	"os"

	"github.com/drizzleent/wallet/app"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		os.Exit(1)
	}

	a.Run()
}
