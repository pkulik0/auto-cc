package server

import (
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

func readPb[T proto.Message](r *http.Request, v T) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return proto.Unmarshal(data, v)
}

func writePb[T proto.Message](w http.ResponseWriter, v T) {
	data, err := proto.Marshal(v)
	if err != nil {
		errLog(w, err, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	writeOrLog(w, data)
}

func errLog(w http.ResponseWriter, err error, message string, status int) {
	log.Error().Err(err).Msg(message)
	switch status {
	case http.StatusBadRequest:
		http.Error(w, "Bad request", status)
	case http.StatusInternalServerError:
		http.Error(w, "Internal server error", status)
	default:
		log.Error().Msg("invalid status")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func writeOrLog(w http.ResponseWriter, data []byte) {
	_, err := w.Write(data)
	if err != nil {
		log.Error().Err(err).Msg("failed to write response")
	}
}
