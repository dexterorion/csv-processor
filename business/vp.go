package business

import (
	"context"
	"fmt"

	"github.com/soap-parser/model"
	"github.com/soap-parser/mongo"
)

const (
	VALID = iota
	DEVIATION
	INVALID
)

type VP interface {
	Process(ctx context.Context) error
}

type vpImpl struct {
	dbAcess  *mongo.DB
	file     []byte
	filetype string
	parking  model.Parking
}

func NewAuconMonza(dbAcess *mongo.DB, file []byte, filetype string, parking model.Parking) VP {
	return &vpImpl{
		dbAcess:  dbAcess,
		file:     file,
		filetype: filetype,
		parking:  parking,
	}
}

func (s *vpImpl) Process(ctx context.Context) error {
	switch s.filetype {
	case "transacoes":
		return s.transactionsProcess(ctx)
	default:
		return fmt.Errorf("filetype [%s] does not exists for Monza parking", s.filetype)
	}
}

func (s *vpImpl) transactionsProcess(ctx context.Context) error {

	return nil
}
