package redirect

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/slg"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=UrlGetter --name=AnalyticsServ
type UrlGetter interface {
	GetUrl(alias string) (string, error)
}
type AnalyticsServ interface {
	SendEvent(ctx context.Context, name string, date *timestamppb.Timestamp) error
}

func New(log *slog.Logger, urlGetter UrlGetter, analyticsServ AnalyticsServ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"
		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		err := analyticsServ.SendEvent(context.Background(), "Redirect", timestamppb.New(time.Now()))
		if err != nil {
			log.Error("failed on send to analytics send event", slg.Err(err))
		}
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("alias not found"))
			return
		}
		resUrl, err := urlGetter.GetUrl(alias)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("such url not found", "alias:", alias)
			render.JSON(w, r, resp.Error("url not found"))
			return
		}
		if err != nil {
			log.Error("failed to get url", slg.Err(err))
			render.JSON(w, r, resp.Error("failed to get url"))
			return
		}
		log.Info("got url", slog.String("url", resUrl))
		http.Redirect(w, r, resUrl, http.StatusFound)

	}
}
