package business

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/csv-processor/model"
	"github.com/csv-processor/mongo"
	"go.uber.org/zap"
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
	reader   *csv.Reader
	filetype string
	parking  model.Parking
	cn       chan *model.Line
	logger   *zap.Logger
}

func NewVP(dbAcess *mongo.DB, reader *csv.Reader, filetype string, parking model.Parking) VP {
	log, _ := zap.NewProduction()

	return &vpImpl{
		dbAcess:  dbAcess,
		reader:   reader,
		filetype: filetype,
		parking:  parking,
		cn:       make(chan *model.Line),
		logger:   log,
	}
}

func (s *vpImpl) Process(ctx context.Context) error {
	switch s.filetype {
	case "transactions":
		return s.transactionsProcess(ctx)
	default:
		return fmt.Errorf("filetype [%s] does not exists for parking", s.filetype)
	}
}

func (s *vpImpl) transactionsProcess(ctx context.Context) error {

	for {
		line, err := s.reader.Read()

		if err != nil {
			if err == io.EOF {
				s.logger.Info("Done")
				return nil
			} else {
				return err
			}
		}

		cin, err := time.Parse("10/1/20 12:56", line[4])
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(cin)
	}
}
