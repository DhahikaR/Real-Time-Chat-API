// Package main provides API documentation for Real Time Chat API
//
// # Documentation of Real Time Chat API
//
// Terms of service: http://swagger.io/terms/
//
// Schemes: http, https
// Version: 1.0
// License: MIT http://opensource.org/licenses/MIT
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// SecurityDefinitions: bearer
//
//	Authorization: (Bearer) true
//	token: (JWT Access Token) true
//
// swagger:meta
package main

// @title           Real Time Chat API + WebSocket
// @version         1.0
// @description     Documentation API for Real Time Chat using Fiber + WebSocket
// @host            localhost:8080
// @schemes         http https
// @basePath        /
// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com
// @license.name   Apache 2.0
// @license.url    http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description       Type "Bearer" and JWT token value of the logged in user to access this API.
