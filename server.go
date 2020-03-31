package main

import (
	api "github.com/sys520084/namenode/api"
)

func main() {
	r := api.SetupRouter()
	r.Run(":8080")
}
