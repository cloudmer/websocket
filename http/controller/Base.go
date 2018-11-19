package controller

import "net/http"

type Basecontroller struct {
	Writer http.ResponseWriter
	Request *http.Request
}