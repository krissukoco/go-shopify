package shopify

import (
	"net/http"
	"time"
)

type NonceGenerateFunc func() (string, error)

type AppOption func(a *App)

type App struct {
	cfg           ShopifyAppConfig
	httpClient    *http.Client
	nonceGenerate NonceGenerateFunc
	accessMode    AccessMode
}

type ShopifyAppConfig struct {
	ClientId     string
	ClientSecret string
	Scopes       string
	RedirectUrl  string
}

func NewApp(config ShopifyAppConfig, opts ...AppOption) *App {
	app := &App{
		cfg: config,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		nonceGenerate: DefaultNonceGenerate,
		accessMode:    OfflineAccessMode,
	}
	for _, opt := range opts {
		opt(app)
	}
	return app
}

func WithNonceGeneratorFunction(fn NonceGenerateFunc) AppOption {
	return func(a *App) {
		a.nonceGenerate = fn
	}
}

func WithHttpClient(c *http.Client) AppOption {
	return func(a *App) {
		a.httpClient = c
	}
}
