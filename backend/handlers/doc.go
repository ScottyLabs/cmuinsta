// Package handlers implements the HTTP handlers for the CMU Insta REST API.
//
// Each file corresponds to a resource group from RFC 2 and implements several
// endpoints. The endpoints rely on JWT for authentication and use a function
// from [middleware.Authenticate] to determine an AndrewID from this JWT. They
// use GORM to manipulate data entried directly after authentication.
//
// These handlers are utilized by the [router] to make the REST API with Gin.

package handlers
