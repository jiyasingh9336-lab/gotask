package auth

import "net/http"

func GetUserID(r *http.Request) int64 {
	if id, ok := r.Context().Value("userID").(int64); ok {
		return id
	}
	return 0
}

func GetUserRole(r *http.Request) string {
	if role, ok := r.Context().Value("role").(string); ok {
		return role
	}

	return ""
}
