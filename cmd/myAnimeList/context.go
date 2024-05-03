package main

import (
	"context"
	"net/http"

	"final-project/pkg/model"
)

type contextKey string

const userContextKey = contextKey("user")

// contextSetUser returns a new copy of the request with the provided User struct added to the
func (app *application) contextSetUser(r *http.Request, user *model.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *model.User {
	user, ok := r.Context().Value(userContextKey).(*model.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
