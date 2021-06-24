package auth

type LoginData struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Verify   string `json:"verify"`
}
