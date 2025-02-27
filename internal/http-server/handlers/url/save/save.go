package save

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/slg"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}
type Response struct {
	resp.Response
	Alias string `json:"alias"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=UrlSaver --name=AnalyticsServ
type UrlSaver interface {
	SaveUrl(urlToSave, alias string) (int64, error)
}
type AnalyticsServ interface {
	SendEvent(ctx context.Context, name string, date *timestamppb.Timestamp) error
}

func responseOk(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}

const aliasLength = 6

func New(log *slog.Logger, urlSaver UrlSaver, analyticsServ AnalyticsServ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		err := analyticsServ.SendEvent(context.Background(), "SaveUrl", timestamppb.New(time.Now()))
		if err != nil {
			log.Error("failed on send to analytics send event", slg.Err(err))
		}
		var req Request
		err = render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		if err != nil {
			log.Error("failed to decode json", slg.Err(err))
			render.JSON(w, r, resp.Error("failed to decode json"))
			return
		}
		log.Info("request body decoded", slog.Any("req", req))
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", slg.Err(err))
			render.JSON(w, r, resp.Error(validateErr.Error()))
			return
		}
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}
		id, err := urlSaver.SaveUrl(req.URL, alias)
		if errors.Is(err, storage.ErrUrlExist) {
			log.Info("such url is already exist", slog.String("url", req.URL))
			render.JSON(w, r, resp.Error("url already exists"))
			return
		}
		if err != nil {
			log.Error("failed to add url", slg.Err(err))
			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}
		log.Info("url added", slog.Int64("id", id))
		responseOk(w, r, alias)
	}
}
