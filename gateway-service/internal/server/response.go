package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"gateway-service/internal/pkg/requestid"

	kratoserrors "github.com/go-kratos/kratos/v2/errors"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

type responseEnvelope struct {
	Code       int32           `json:"code"`
	Message    string          `json:"message"`
	RequestID  string          `json:"request_id"`
	ServerTime int64           `json:"server_time"`
	Data       json.RawMessage `json:"data"`
}

func requestIDFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(requestid.Header)
		if rid == "" {
			rid = newRequestID()
		}
		w.Header().Set(requestid.Header, rid)
		next.ServeHTTP(w, r.WithContext(requestid.NewContext(r.Context(), rid)))
	})
}

func unifiedResponseEncoder(w http.ResponseWriter, r *http.Request, v any) error {
	if v == nil {
		return nil
	}
	codec, _ := khttp.CodecForRequest(r, "Accept")
	data, err := codec.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(responseEnvelope{
		Code:       0,
		Message:    "ok",
		RequestID:  requestid.FromContext(r.Context()),
		ServerTime: time.Now().UnixMilli(),
		Data:       data,
	})
}

func unifiedErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	se := kratoserrors.FromError(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(se.Code))
	_ = json.NewEncoder(w).Encode(responseEnvelope{
		Code:       se.Code,
		Message:    se.Message,
		RequestID:  requestid.FromContext(r.Context()),
		ServerTime: time.Now().UnixMilli(),
		Data:       nil,
	})
}

func newRequestID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return time.Now().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(b[:])
}
