// Copyright 2017 The go-ego Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-ego/ego/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package main

import (
	"github.com/go-ego/ego"
)

var DB = make(map[string]string)

func main() {
	// Disable Console Color
	// ego.DisableConsoleColor()
	r := ego.Default()

	// Ping test
	r.GET("/ping", func(c *ego.Context) {
		c.String(200, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *ego.Context) {
		user := c.Params.ByName("name")
		value, ok := DB[user]
		if ok {
			c.JSON(200, ego.Map{"user": user, "value": value})
		} else {
			c.JSON(200, ego.Map{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses ego.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(ego.BasicAuth(ego.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", ego.BasicAuth(ego.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	authorized.POST("admin", func(c *ego.Context) {
		user := c.MustGet(ego.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			DB[user] = json.Value
			c.JSON(200, ego.Map{"status": "ok"})
		}
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
