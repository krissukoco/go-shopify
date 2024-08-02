package shopify

import (
	"log"
	"net/http"
)

type HttpHandler struct {
	app     *App
	options *HttpHandlerOptions
}

func (a *App) NewHandler(opt *HttpHandlerOptions) *HttpHandler {
	return &HttpHandler{a, opt}
}

func (h *HttpHandler) HandleAuthorization(w http.ResponseWriter, r *http.Request) {
	shopQuery := r.URL.Query().Get("shop")
	if shopQuery == "" {
		http.Error(w, "Query 'shop' is required", http.StatusUnprocessableEntity)
		return
	}
	shop := ParseShop(shopQuery)

	authorizeUrl, err := h.app.GetAuthorizeUrl(shop, h.app.accessMode)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, authorizeUrl, http.StatusFound)
}

func (h *HttpHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	valid, err := h.app.VerifyAuthorizationQuery(query)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !valid {
		http.Error(w, "Error validating request", http.StatusBadRequest)
		return
	}

	shop := ParseShop(query.Get("shop"))
	code := query.Get("code")
	token, err := h.app.GetAccessToken(r.Context(), shop, code)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	log.Printf("Access Token response: %+v", *token)

	if h.options != nil {
		if h.options.AccessTokenGeneratedFunc != nil {
			err = h.options.AccessTokenGeneratedFunc(token)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
	}
}
