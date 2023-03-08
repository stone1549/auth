package service

import (
	"context"
	"encoding/json"
	"net/http"
)

type newSessionRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type newSessionResponse struct {
	Token string `json:"token"`
}

func (nsr newSessionResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// NewSessionMiddleware middleware to authenticate a user from the request parameters
func NewSessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var reqUser newSessionRequest
		err := decoder.Decode(&reqUser)
		if err != nil {
			RenderResponse(w, r, NewBadRequestErr("request body invalid"))
			return
		}

		if reqUser.Email == "" {
			RenderResponse(w, r, NewBadRequestErr("email is required"))
			return
		}

		if reqUser.Password == "" {
			RenderResponse(w, r, NewBadRequestErr("password is required"))
			return
		}

		userRepo, ok := r.Context().Value("repo").(UserRepository)

		if !ok {
			RenderResponse(w, r, NewInternalServerErr("repo not found"))
			return
		}

		user, err := userRepo.Authenticate(reqUser.Email, reqUser.Password)

		if err != nil {
			RenderResponse(w, r, NewInternalServerErr("repo error"))
			return
		} else if user.Id == "" {
			RenderResponse(w, r, NewUnauthorizedErr("login failed"))
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

// NewSession responds to authentication request with jwt token or appropriate error
func NewSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token, ok := ctx.Value("token").(string)

	if !ok {
		RenderResponse(w, r, NewInternalServerErr("internal error"))
		return
	}

	RenderResponse(w, r, newSessionResponse{token})
}
