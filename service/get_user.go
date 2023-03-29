package service

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type getUserRequest struct {
	Id string `json:"id"`
}

type getUserResponse struct {
	User `json:"user"`
}

func (nsr getUserResponse) Render(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(http.StatusOK)

	return nil
}

// GetUserMiddleware middleware to retrieve a user from the repo
func GetUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqUser getUserRequest

		reqUser.Id = chi.URLParam(r, "id")

		if reqUser.Id == "" {
			RenderResponse(w, r, NewBadRequestErr("id is required in path"))
			return
		}

		userRepo, ok := r.Context().Value("repo").(UserRepository)

		if !ok {
			RenderResponse(w, r, NewInternalServerErr("internal error"))
			return
		}

		user, err := userRepo.GetUser(reqUser.Id)

		if err != nil {
			RenderResponse(w, r, NewInternalServerErr("internal error"))
			return
		}

		ctx := context.WithValue(r.Context(), "user", &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUser renders the response to the get user request.
func GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value("user").(*User)

	if !ok {
		RenderResponse(w, r, NewInternalServerErr("internal error"))
		return
	}

	RenderResponse(w, r, getUserResponse{*user})
}
