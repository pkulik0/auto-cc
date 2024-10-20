package helpers

import (
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

// ReadPb reads a protobuf message from the request body.
func ReadPb[T proto.Message](r *http.Request, v T) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return proto.Unmarshal(data, v)
}

// WritePb writes a protobuf message to the response body.
func WritePb[T proto.Message](w http.ResponseWriter, v T) {
	data, err := proto.Marshal(v)
	if err != nil {
		ErrLog(w, err, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	WriteOrLog(w, data)
}

// ErrLog logs an error and writes an error message to the response.
func ErrLog(w http.ResponseWriter, err error, message string, status int) {
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

// WriteOrLog writes data to the response or logs an error.
func WriteOrLog(w http.ResponseWriter, data []byte) {
	_, err := w.Write(data)
	if err != nil {
		log.Error().Err(err).Msg("failed to write response")
	}
}
