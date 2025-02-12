package response

type ErrorResponse struct {
	Errors string `json:"errors"`
}

//func SendErrorResponse(w http.ResponseWriter, statusCode int, errMsg string) {
//	w.WriteHeader(statusCode)
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(ErrorResponse{Errors: errMsg})
//}
