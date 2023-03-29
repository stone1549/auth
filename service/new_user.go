package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type newUserRequest struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	UserProfile `json:"profile"`
}

type newUserResponse struct {
	Token string `json:"token"`
}

func (nsr newUserResponse) Render(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(http.StatusOK)

	return nil
}

// NewUserMiddleware middleware to add a new user to the repo from the request parameters
func NewUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var reqUser newUserRequest
		err := decoder.Decode(&reqUser)
		if err != nil {
			RenderResponse(w, r, NewBadRequestErr("invalid request body"))
			return
		}

		if reqUser.Email == "" {
			RenderResponse(w, r, NewBadRequestErr("email is required"))
			return
		}

		if reqUser.Username == "" {
			RenderResponse(w, r, NewBadRequestErr("handle is required"))
			return
		}

		if reqUser.Password == "" {
			RenderResponse(w, r, NewBadRequestErr("password is required"))
			return
		}

		userRepo, ok := r.Context().Value("repo").(UserRepository)

		if !ok {
			RenderResponse(w, r, NewInternalServerErr("internal error"))
			return
		}

		id, err := userRepo.NewUser(
			reqUser.Email,
			reqUser.Username,
			reqUser.Password,
			reqUser.UserProfile.Gender,
			reqUser.UserProfile.Age,
			reqUser.UserProfile.Topics,
		)

		if err != nil {
			RenderResponse(w, r, NewInternalServerErr("internal error"))
			log.Println(err)
			return
		}

		tokenFactory, ok := r.Context().Value("tokenFactory").(TokenFactory)

		if !ok {
			RenderResponse(w, r, NewInternalServerErr("internal error"))
			return
		}

		token, err := tokenFactory.NewToken(NewClaims(id, reqUser.Email, reqUser.Username))

		if err != nil {
			RenderResponse(w, r, NewInternalServerErr("internal error"))
			return
		}

		ctx := context.WithValue(r.Context(), "token", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// NewUser renders the response to the product update request.
func NewUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token, ok := ctx.Value("token").(string)

	if !ok {
		RenderResponse(w, r, NewInternalServerErr("internal error"))
		return
	}

	RenderResponse(w, r, newUserResponse{token})
}
