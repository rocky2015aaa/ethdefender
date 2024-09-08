package main

import (
	"context"
	"fmt"

	"github.com/rocky2015aaa/ethdefender/internal/config"
	"github.com/rocky2015aaa/ethdefender/internal/services/preventer"
)

func main() {
	fmt.Printf("Build Date: %s\nBuild Version: %s\nBuild: %s\n\n", config.Date, config.Version, config.Build)
	preventer.NewApp().Run(context.Background())
}
