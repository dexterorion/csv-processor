package business

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/soap-parser/model"
	"github.com/soap-parser/mongo"
)

type AuconMonza interface {
	Process(ctx context.Context) error
}

type auconMonzaImpl struct {
	dbAcess *mongo.DB
	file    []byte
}

func NewAuconMonza(dbAcess *mongo.DB, file []byte) AuconMonza {
	return &auconMonzaImpl{
		dbAcess: dbAcess,
		file:    file,
	}
}

func (am *auconMonzaImpl) Process(ctx context.Context) error {
	var v model.AuconEnvelope
	err := xml.Unmarshal(am.file, &v)

	if err != nil {
		return err
	}

	for _, saida := range v.Body.GetResponse.GetResult.Saidas {
		fmt.Println(saida)
	}

	return nil
}
