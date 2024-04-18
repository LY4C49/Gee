package main

import (
	"gee/gee"
	"net/http"
)

//func testMethod1(w http.ResponseWriter, r *http.Request) {
//	fmt.Println("Current Time is: ", time.Now())
//}

func main() {
	r := gee.New()
	//r.GET("/", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Println("Hello, I am a server")
	//})
	//r.GET("/time", testMethod1)

	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello, I am %s\n", "World")
	})

	r.GET("/hello", func(c *gee.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run("localhost:8888")
}
