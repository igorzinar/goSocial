package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/igorzinar/goSocial/internal/store"
	"log"
	"net/http"
	"strconv"
)

type userKey string

const userCtx userKey = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	log.Println("Extracted user:", user)
	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

type FollowUser struct {
	UserID string `json:"user_id"`
}

//func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
//	followerUser := getUserFromContext(r.Context())
//
//	// Revert back to auth userID from ctx
//	var payload FollowUser
//
//	if err := readJSON(w, r, &payload); err != nil {
//		app.badRequestResponse(w, r, err)
//		return
//	}
//
//	ctx := r.Context()
//	app.store.Users.Follow(ctx, followerUser.ID, payload.UserID)
//	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
//		app.internalServerError(w, r, err)
//		return
//	}
//}
//
//func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
//	user := getUserFromContext(r.Context())
//	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
//		app.internalServerError(w, r, err)
//		return
//	}
//}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "userID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, id)
		log.Println("Extracted user:", user)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func getUserFromContext(ctx context.Context) *store.User {
	user, _ := ctx.Value(userCtx).(*store.User)
	return user
}
