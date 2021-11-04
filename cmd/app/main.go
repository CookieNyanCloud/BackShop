package main

import (
	"flag"
	"github.com/cookienyancloud/back/internal/app"
	_ "github.com/lib/pq"
)

const configsDir = "configs"

func main() {
	var local bool
	flag.BoolVar(&local, "local", false, "хост")
	flag.Parse()
	app.Run(configsDir, local)
}
