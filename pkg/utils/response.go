// Package utils provides shared HTTP helpers.
package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/projuktisheba/pse-api-v1/internal/model"
)

// readJSON read json from request body into data. It accepts a single JSON of 1MB max size value in the body
func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 //maximum allowable bytes is 1MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})

	if err != io.EOF {
		return errors.New("body must only have a single JSON value")
	}

	return nil
}

// writeJSON writes arbitrary data out as json
func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	//add the headers if exists
	if len(headers) > 0 {
		for i, v := range headers[0] {
			w.Header()[i] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(out)
	return nil
}

// ---- Standard Response Helpers ----

// OK sends a successful response with data.
func OK(w http.ResponseWriter, message string, data interface{}) {
	WriteJSON(w, http.StatusOK, model.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a 201 Created response.
func Created(w http.ResponseWriter, message string, data interface{}) {
	WriteJSON(w, http.StatusCreated, model.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// BadRequest sends a 400 error response.
func BadRequest(w http.ResponseWriter, err error) {
	WriteJSON(w, http.StatusBadRequest, model.APIResponse{
		Success: false,
		Message: "bad request",
		Error:   sanitizeError(err),
	})
}

// Unauthorized sends a 401 error response.
func Unauthorized(w http.ResponseWriter, err error) {
	WriteJSON(w, http.StatusUnauthorized, model.APIResponse{
		Success: false,
		Message: "unauthorized",
		Error:   sanitizeError(err),
	})
}

// Forbidden sends a 403 error response.
func Forbidden(w http.ResponseWriter, err error) {
	WriteJSON(w, http.StatusForbidden, model.APIResponse{
		Success: false,
		Message: "forbidden",
		Error:   sanitizeError(err),
	})
}

// NotFound sends a 404 error response.
func NotFound(w http.ResponseWriter, err error) {
	WriteJSON(w, http.StatusNotFound, model.APIResponse{
		Success: false,
		Message: "not found",
		Error:   sanitizeError(err),
	})
}

// ValidationError sends a 422 error response.
func ValidationError(w http.ResponseWriter, err error) {
	WriteJSON(w, http.StatusUnprocessableEntity, model.APIResponse{
		Success: false,
		Message: "validation error",
		Error:   sanitizeError(err),
	})
}

// ServerError sends a 500 error response.
func ServerError(w http.ResponseWriter, err error) {
	WriteJSON(w, http.StatusInternalServerError, model.APIResponse{
		Success: false,
		Message: "internal server error",
		Error:   sanitizeError(err),
	})
}

// sanitizeError prevents leaking sensitive error details.
func sanitizeError(err error) string {
	if err == nil {
		return ""
	}

	msg := err.Error()

	// Avoid leaking SQL or internal stack info
	banned := []string{
		"sql:", "pq:", "pgx:", "redis:", "syntax error",
	}

	for _, b := range banned {
		if strings.Contains(strings.ToLower(msg), b) {
			return "internal error"
		}
	}

	return msg
}
