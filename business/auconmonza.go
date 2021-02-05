package business

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"github.com/soap-parser/model"
	"github.com/soap-parser/mongo"
)

const (
	VALID = iota
	DEVIATION
	INVALID
)

type AuconMonza interface {
	Process(ctx context.Context) error
}

type auconMonzaImpl struct {
	dbAcess  *mongo.DB
	file     []byte
	filetype string
	parking  model.Parking
}

func NewAuconMonza(dbAcess *mongo.DB, file []byte, filetype string, parking model.Parking) AuconMonza {
	return &auconMonzaImpl{
		dbAcess:  dbAcess,
		file:     file,
		filetype: filetype,
		parking:  parking,
	}
}

func (am *auconMonzaImpl) Process(ctx context.Context) error {
	switch am.filetype {
	case "pagamentos":
		return am.pagamentosProcess(ctx)
	case "saidas":
		return am.saidasProcess(ctx)
	case "credenciados":
		return am.credenciaisProcess(ctx)
	default:
		return fmt.Errorf("filetype [%s] does not exists for Monza parking", am.filetype)
	}
}

func (am *auconMonzaImpl) credenciaisProcess(ctx context.Context) error {
	var v model.AuconCredenciadoEnvelope
	err := xml.Unmarshal(am.file, &v)

	if err != nil {
		return err
	}

	for _, cred := range v.Body.GetResponse.GetResult.Credenciados {
		transactions, err := am.dbAcess.TransactionCollection.GetAllByMatricula(ctx, am.parking.ID, cred.MATRICULA)

		if err != nil {
			return err
		}

		for _, transaction := range transactions {
			transaction.Categoria = cred.CATEGORIA
			transaction.UseType = cred.CATEGORIA

			transaction.Status = resolveStatus(&transaction)
			transaction.IsValid = resolveValid(&transaction)

			_, err = am.dbAcess.TransactionCollection.Update(ctx, &transaction)

			if err != nil {
				fmt.Fprintf(os.Stderr, "error changing transaction: [%s]\n", err.Error())
			}
		}
	}

	return nil
}

func (am *auconMonzaImpl) pagamentosProcess(ctx context.Context) error {
	var v model.AuconPagamentoEnvelope
	err := xml.Unmarshal(am.file, &v)

	if err != nil {
		return err
	}

	for _, pagto := range v.Body.GetResponse.GetResult.Diffgram.DocElement.Pagtos {
		transaction, err := am.dbAcess.TransactionCollection.GetByTicketAndMatricula(ctx, pagto.Ticket, am.parking.ID, pagto.Matricula)

		pd, err := time.Parse("2006-01-02T15:04:05-03:00", pagto.Data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing checking date: [%s]\n", err.Error())
			continue
		}

		if err != nil || transaction.PaymentDate.Equal(pd) {
			transaction = &model.Transaction{}
		}

		transaction.UseType = getUseTypePagto(pagto)
		transaction.Sequence = pagto.Ticket
		transaction.FareAmount = pagto.Valor
		transaction.PaidAmount = pagto.Valor
		transaction.Discount = pagto.Desconto
		transaction.ParkingInfo = am.parking
		transaction.PaymentMethod = getPaymentMethod(pagto.TpPagamento)
		transaction.OfferType = "On-demand"

		transaction.Matricula = pagto.Matricula
		transaction.PaymentDate = pd

		transaction.Status = resolveStatus(transaction)
		transaction.IsValid = resolveValid(transaction)

		if transaction.ID.IsZero() {
			transaction, err = am.dbAcess.TransactionCollection.Create(ctx, transaction)
		} else {
			transaction, err = am.dbAcess.TransactionCollection.Update(ctx, transaction)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "error changing transaction: [%s]\n", err.Error())
		}
	}

	return nil
}

func resolveValid(transaction *model.Transaction) bool {
	if transaction.Status == INVALID || transaction.Status == DEVIATION {
		return false
	}

	if transaction.PaymentMethod == "Cancelado" {
		return false
	}

	return true
}

/*
	{
	"_id" : ObjectId("60086c095b8f09069061d147"),
	"time_interval_hour" : [
		ISODate("2020-10-01T09:00:00Z"),
		ISODate("2020-10-01T10:00:00Z"),
		ISODate("2020-10-01T11:00:00Z"),
		ISODate("2020-10-01T12:00:00Z")
	],
	"checkin_date" : ISODate("2020-10-01T08:45:46Z"),
	"checkout_date" : ISODate("2020-10-01T12:16:24Z"),
	"payment_date" : ISODate("0001-01-01T00:00:00Z"),
	"fare_amount" : 0,
	"fare_name" : "",
	"paid_amount" : 0,
	"discount" : 0,
	"payment_data" : "",
	"payment_method" : "",
	"use_type" : "Proprietário",
	"offer_type" : "On-demand",
	"duration" : 4,
	"is_valid" : false,
	"parking_info" : {
		"id" : NumberLong(6),
		"name" : "Monza",
		"slug" : "monza"
	},
	"status" : 0,
	"cash_register_id" : NumberLong(0),
	"sequence" : "0",
	"fiscal" : "",
	"partial" : "",
	"matricula" : "982",
	"categoria" : "Proprietário",
	"version" : 2,
	"schema" : 1,
	"created_at" : ISODate("2021-01-20T17:44:41.838Z"),
	"updated_at" : ISODate("2021-01-20T20:18:18.748Z")
}


*/

func resolveStatus(transaction *model.Transaction) int {
	if transaction.Sequence == "0" && transaction.Matricula != "0" && transaction.FareAmount == 0 {
		return VALID
	}

	if transaction.Sequence == "0" && transaction.Matricula != "0" && transaction.FareAmount != 0 {
		return DEVIATION
	}

	if transaction.Sequence == "0" && transaction.Matricula == "0" {
		return DEVIATION
	}

	if transaction.Sequence != "0" && transaction.Matricula != "0" && transaction.FareAmount == 0 {
		return VALID
	}

	if transaction.Sequence != "0" && transaction.Matricula == "0" && transaction.FareAmount == 0 && transaction.PaymentMethod != "Isento" {
		return DEVIATION
	}

	if transaction.Sequence != "0" && transaction.Matricula == "0" && transaction.FareAmount == 0 && transaction.PaymentMethod == "Isento" {
		return VALID
	}

	if transaction.Sequence != "0" && transaction.Matricula == "0" && transaction.FareAmount != 0 {
		return VALID
	}

	// apply fare rules

	return INVALID
}

func getPagtoValid(value string) bool {
	return value != "Other" && value != "C"
}

func getUseTypePagto(pagto model.AuconPagamento) string {
	if pagto.Ticket != "0" && pagto.Matricula == "0" {
		return "Avulso"
	}

	return ""
}

func getUseTypeSaida(saida model.AuconSaida) string {
	if saida.Ticket != "0" && saida.Matricula == "0" {
		return "Avulso"
	}

	return ""
}

func getPaymentMethod(value string) string {
	switch value {
	case "CC":
		return "Creditcard"
	case "CD":
		return "Debitcard"
	case "DI":
		return "Dinheiro"
	case "CA":
		return "Cancelado"
	case "IS":
		return "Isento"
	case "CH":
		return "Cheque"
	default:
		return "Outro"
	}
}

func (am *auconMonzaImpl) saidasProcess(ctx context.Context) error {
	var v model.AuconSaidaEnvelope
	err := xml.Unmarshal(am.file, &v)

	if err != nil {
		return err
	}

	for _, saida := range v.Body.GetResponse.GetResult.Saidas {
		transaction, err := am.dbAcess.TransactionCollection.GetByTicketAndMatricula(ctx, saida.Ticket, am.parking.ID, saida.Matricula)

		ci, err := time.Parse("2006-01-02T15:04:05", saida.DataEnt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing checking date: [%s]\n", err.Error())
			continue
		}

		co, err := time.Parse("2006-01-02T15:04:05", saida.DataSai)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing checkout date: [%s]\n", err.Error())
			continue
		}

		if err != nil || !transaction.CheckinDate.Equal(ci) || !transaction.CheckoutDate.Equal(co) {
			transaction = &model.Transaction{}
		}

		transaction.CheckinDate = ci
		transaction.CheckoutDate = co
		transaction.Matricula = saida.Matricula

		co = time.Date(co.Year(), co.Month(), co.Day(), co.Hour(), 0, 0, 0, co.Location())

		if ci.Minute() > 0 || ci.Second() > 0 {
			ci = time.Date(ci.Year(), ci.Month(), ci.Day(), ci.Hour()+1, 0, 0, 0, ci.Location())
		}

		transaction.TimeIntervalHour = []time.Time{}
		transaction.Duration = 0

		for {
			if ci.After(co) {
				break
			}

			transaction.Duration = transaction.Duration + 1
			transaction.TimeIntervalHour = append(transaction.TimeIntervalHour, ci)
			ci = ci.Add(1 * time.Hour)
		}

		transaction.OfferType = "On-demand"
		transaction.Sequence = saida.Ticket
		transaction.ParkingInfo = am.parking

		transaction.UseType = getUseTypeSaida(saida)
		transaction.Status = resolveStatus(transaction)
		transaction.IsValid = resolveValid(transaction)

		if transaction.ID.IsZero() {
			transaction, err = am.dbAcess.TransactionCollection.Create(ctx, transaction)
		} else {
			transaction, err = am.dbAcess.TransactionCollection.Update(ctx, transaction)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "error changing transaction: [%s]\n", err.Error())
		}
	}

	return nil
}
