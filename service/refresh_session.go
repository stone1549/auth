package service

import (
	"context"
	"net/http"
)

type refreshSessionResponse struct {
	Token string `json:"token"`
}

func (rsr refreshSessionResponse) Render(res http.ResponseWriter, _ *http.Request) error {
	res.WriteHeader(200)
	return nil
}

// RefreshSessionMiddleware middleware to authenticate a user from the request parameters
func RefreshSessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(User)

		if !ok {
			RenderResponse(w, r, NewInternalServerErr("internal error"))
			return
		}

		tokenFactory, ok := r.Context().Value("tokenFactory").(TokenFactory)

		if !ok {
			RenderResponse(w, r, NewInternalServerErr("internal error"))
			return
		}

		token, err := tokenFactory.NewToken(NewClaims(user.Id, user.Email, user.Username))

		if err != nil {
			RenderResponse(w, r, NewInternalServerErr("internal error"))
			return
		}

		ctx := context.WithValue(r.Context(), "token", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RefreshSession responds to authenticated requests with a new token
func RefreshSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token, ok := ctx.Value("token").(string)

	if !ok {
		RenderResponse(w, r, NewInternalServerErr("internal error"))
		return
	}

	RenderResponse(w, r, refreshSessionResponse{token})
}
