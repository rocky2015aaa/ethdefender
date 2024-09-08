package main

import (
	"context"
	"fmt"

	_ "github.com/rocky2015aaa/ethdefender/docs"
	"github.com/rocky2015aaa/ethdefender/internal/config"
	"github.com/rocky2015aaa/ethdefender/internal/services/reporter"
)

// @title   ETH Defender Service
// @version 1.0.0

// @contact.name  Donggeon Lee
// @contact.email rocky2010aaa@gmail.com

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host     localhost
// @BasePath /
func main() {
	fmt.Printf("Build Date: %s\nBuild Version: %s\nBuild: %s\n\n", config.Date, config.Version, config.Build)
	reporter.NewApp().Run(context.Background())
}
