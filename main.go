package main

import (
	"os"

	"log/slog"

	"github.com/gofrs/uuid"
)

const (
	attrError   = "error"
	attrVersion = "version"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	var idv4 uuid.UUID
	var err error
	if idv4, err = uuid.NewV4(); err != nil {
		logger.Error("Error generating V4 UUID", attrError, err)
	}
	logger.Info(idv4.String(), attrVersion, idv4.Version())
	var idv6 uuid.UUID
	if idv6, err = uuid.NewV7(); err != nil {
		logger.Error("Error generating V6 UUID", attrError, err)
	}
	logger.Info(idv6.String(), attrVersion, idv6.Version())

}
