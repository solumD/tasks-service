package main

import (
	"context"

	"github.com/solumD/tasks-service/internal/app"
)

func main() {
	ctx := context.Background()
	app.InitAndRun(ctx)
}
