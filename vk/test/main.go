package main

import (
	"log"
	"net/http"

	"github.com/suborbital/vektor/vk"
	"github.com/suborbital/vektor/vlog"
)

type testMeta struct {
	Version string `json:"version"`
}

func main() {
	logger := vlog.Default(
		vlog.Level(vlog.LogLevelTrace),
		vlog.AppMeta(testMeta{Version: "v0.1.1"}),
		vlog.ToFile("/Users/cohix-16/.op/logfile.log"),
	)

	server := vk.New(
		vk.UseAppName("vk tester"),
		vk.UseLogger(logger),
		vk.UseEnvPrefix("APP"),
		vk.UseInsecureHTTP(8080),
	)

	server.GET("/f", HandleFound)
	server.POST("/f", HandleFound)
	server.GET("/nf", HandleNotFound)

	v1 := vk.Group("/v1").Before(denyMiddleware, headerMiddleware)
	v1.GET("/me", HandleMe)
	v1.GET("/me/hack", HandleMe)

	v2 := vk.Group("/v2").Before(setScopeMiddleware).After(logAfter)
	v2.GET("/you", HandleYou)
	v2.GET("/mistake", HandleBadMistake)

	api := vk.Group("/api").After(getSetLogAfterware)
	api.AddGroup(v1)
	api.AddGroup(v2)

	server.AddGroup(api)

	server.HandleHTTP(http.MethodGet, "/http", HandleHTTP)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
