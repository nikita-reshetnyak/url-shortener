package delete

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

type UrlDelete interface {
	DeleteUrl(alias string) (int64, error)
}
type AnalyticsServ interface {
	SendEvent(ctx context.Context, name string, date *timestamppb.Timestamp) error
}
type Response struct {
	resp.Response
	Alias string
}

func New(log *slog.Logger, urlDelete UrlDelete, analyticsServ AnalyticsServ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete.New"
		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		err := analyticsServ.SendEvent(context.Background(), "Delete", timestamppb.New(time.Now()))
		if err != nil {
			log.Error("failed on send to analytics send event", slg.Err(err))
		}
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("alias not found"))
			return
		}
		deletedId, err := urlDelete.DeleteUrl(alias)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("alias to delete not found", "alias:", alias)
			render.JSON(w, r, resp.Error("alias to delete not found"))
			return
		}
		if err != nil {
			log.Error("failed to delete alias", slg.Err(err))
			render.JSON(w, r, resp.Error("failed to delete alias"))
			return
		}
		log.Info("alias deleted success", slog.Int64("deleted_alias:", deletedId))
		render.JSON(w, r, Response{Response: resp.OK(), Alias: alias})
	}
}
