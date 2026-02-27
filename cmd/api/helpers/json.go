package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Envelope map[string]interface{}

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	// Set the header too accept json
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

const maxBytes = 1_048_576

func ReadJson(req *http.Request, dst interface{}) error {
	req.Body = http.MaxBytesReader(nil, req.Body, int64(maxBytes))

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var maxBytesError *http.MaxBytesError

		switch {
		// Malformed JSON
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body continues badly-formed JSON (at character %d)", syntaxError.Offset)

		// Wrong type for a field
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// Unknown Field
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// Empty Body
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// Payload is too large
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		default:
			return err
		}
	}

	if decoder.More() {
		return errors.New("body must only contains a single JSON object")
	}

	return nil
}

func ErrorJson(w http.ResponseWriter, status int, message string) {
	WriteJson(w, status, Envelope{"error": message})
}

func ValidationErrorJson(w http.ResponseWriter, errors map[string]string) {
	WriteJson(w, http.StatusUnprocessableEntity, Envelope{"error": errors})
}
