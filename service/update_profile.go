package service

import (
	"encoding/json"
	"net/http"
)

type updateProfileRequest struct {
	UserId      string
	UserProfile `json:"profile"`
}

type updateProfileResponse struct {
}

func (nsr updateProfileResponse) Render(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(http.StatusOK)

	return nil
}

// UpdateProfileMiddleware middleware to get a user from the repo from the request parameterss
func UpdateProfileMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var reqUser updateProfileRequest
		err := decoder.Decode(&reqUser.UserProfile)
		if err != nil {
			RenderResponse(w, r, NewBadRequestErr("invalid request body"))
			return
		}

		if reqUser.UserId == "" {
			RenderResponse(w, r, NewBadRequestErr("userId is required"))
			return
		}

		userRepo, ok := r.Context().Value("repo").(UserRepository)

		if !ok {
			RenderResponse(w, r, NewInternalServerErr("internal error"))
			return
		}

		err = userRepo.UpdateProfile(
			reqUser.UserId,
			reqUser.UserProfile,
		)

		if err != nil {
			RenderResponse(w, r, NewInternalServerErr("internal error"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// UpdateProfile renders the response to the profile update request.
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	RenderResponse(w, r, updateProfileResponse{})
}
