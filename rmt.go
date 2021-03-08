package main

import (
	"github.com/ratemyteam/rmt/api"
	"github.com/ratemyteam/rmt/common"
)

func main() {
	api.Run(common.CreateServerContext())
}
