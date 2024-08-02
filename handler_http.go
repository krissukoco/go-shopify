package shopify

import "net/http"

func (a *App) HandleAuthorization(w http.ResponseWriter, r *http.Request) {
	shopQuery := r.URL.Query().Get("shop")
	if shopQuery == "" {
		http.Error(w, "Query 'shop' is required", http.StatusUnprocessableEntity)
		return
	}
	shop := ParseShop(shopQuery)

	authorizeUrl, err := a.GetAuthorizeUrl(shop, a.accessMode)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, authorizeUrl, http.StatusFound)
}

func (a *App) HandleCallback(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	valid, err := a.VerifyAuthorizationQuery(query)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !valid {
		http.Error(w, "Error validating request", http.StatusBadRequest)
	}

	shop := ParseShop(query.Get("shop"))
	code := query.Get("code")
	a.GetAccessToken(r.Context(), shop, code)
}
