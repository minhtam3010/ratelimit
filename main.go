package main

import (
	"flag"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minhtam3010/ratelimit/api"
)

type user struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

var users = []user{
	{Id: "546", Username: "John"},
	{Id: "894", Username: "Mary"},
	{Id: "326", Username: "Jane"},
}

func getUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, users)
}

func main() {
	var limit api.Limiter
	flag.Float64Var(&limit.RPS, "rps", 2, "Rate limiter request per second")
	flag.IntVar(&limit.Burst, "burst", 4, "Rate limiter burst, server can handle")
	flag.BoolVar(&limit.Enabled, "enabled", true, "Rate limiter enabled")

	flag.Parse()
	router := gin.Default()

	// TODO define routes
	router.GET("/users", api.ExceedRequestsLimit(getUsers, limit))

	router.Run("localhost:8080")
}
