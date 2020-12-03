package business

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/soap-parser/model"
	"github.com/soap-parser/mongo"
	"os"
	"strconv"
	"time"
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
		transaction, err := am.dbAcess.TransactionCollection.GetByMatricula(ctx, am.parking.ID, cred.MATRICULA)

		if err != nil {
			transaction = &model.Transaction{}
		}

		transaction.Matricula = cred.MATRICULA
		transaction.Categoria = cred.CATEGORIA

		if transaction.ID.IsZero() {
			transaction, err = am.dbAcess.TransactionCollection.Create(ctx, transaction)
		} else {
			transaction.IsValid = true
			transaction, err = am.dbAcess.TransactionCollection.Update(ctx, transaction)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "error changing transaction: [%s]\n", err.Error())
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
		transaction, err := am.dbAcess.TransactionCollection.GetByTicketOrMatricula(ctx, pagto.Ticket, am.parking.ID, pagto.Matricula)

		if err != nil {
			transaction = &model.Transaction{}
		}

		transaction.UseType = getUseType(pagto.Tp)
		transaction.OfferType = "On-demand"
		transaction.Sequence = strconv.Itoa(int(pagto.Ticket))
		transaction.FareAmount = pagto.Valor
		transaction.PaidAmount = pagto.Valor
		transaction.Discount = pagto.Desconto
		transaction.ParkingInfo = am.parking
		transaction.PaymentMethod = getPaymentMethod(pagto.TpPagamento)
		transaction.Status = "valid"
		transaction.Matricula = pagto.Matricula

		if transaction.ID.IsZero() {
			transaction, err = am.dbAcess.TransactionCollection.Create(ctx, transaction)
		} else {
			transaction.IsValid = true
			transaction, err = am.dbAcess.TransactionCollection.Update(ctx, transaction)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "error changing transaction: [%s]\n", err.Error())
		}
	}

	return nil
}

func getPagtoValid(value string) bool {
	return value != "Other" && value != "C"
}

func getUseType(value string) string {
	switch value {
	case "A":
		return "Rotativo"
	case "L":
		return "Local"
	case "R":
		return "Mensalista"
	default:
		return "Outro"
	}
}

func getPaymentMethod(value string) string {
	switch value {
	case "CC":
		return "creditcard"
	case "CD":
		return "debitcard"
	case "DI":
		return "DINHEIRO"
	case "CA":
		return "Cancelado"
	case "IS":
		return "Isento"
	case "CH":
		return "cheque"
	default:
		return "Other"
	}
}

func (am *auconMonzaImpl) saidasProcess(ctx context.Context) error {
	var v model.AuconSaidaEnvelope
	err := xml.Unmarshal(am.file, &v)

	if err != nil {
		return err
	}

	for _, saida := range v.Body.GetResponse.GetResult.Saidas {
		transaction, err := am.dbAcess.TransactionCollection.GetByTicketOrMatricula(ctx, saida.Ticket, am.parking.ID, saida.Matricula)

		if err != nil {
			transaction = &model.Transaction{}
		}

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

		transaction.CheckinDate = ci
		transaction.CheckoutDate = co
		transaction.Matricula = saida.Matricula

		if co.Minute() > 30 {
			co = time.Date(co.Year(), co.Month(), co.Day(), co.Hour()+1, 0, 0, 0, co.Location())
		} else {
			if co.Minute() > 1 {
				co = time.Date(co.Year(), co.Month(), co.Day(), co.Hour(), 30, 0, 0, co.Location())
			} else {
				co = time.Date(co.Year(), co.Month(), co.Day(), co.Hour(), 0, 0, 0, co.Location())
			}
		}

		if ci.Minute() > 30 {
			ci = time.Date(ci.Year(), ci.Month(), ci.Day(), ci.Hour()+1, 0, 0, 0, ci.Location())
		} else {
			ci = time.Date(ci.Year(), ci.Month(), ci.Day(), ci.Hour(), 30, 0, 0, ci.Location())
		}

		transaction.TimeIntervalHour = []time.Time{}
		transaction.Duration = 0

		for {
			if ci.After(co) {
				break
			}

			transaction.Duration = transaction.Duration + 0.5
			transaction.TimeIntervalHour = append(transaction.TimeIntervalHour, ci)
			ci = ci.Add(30 * time.Minute)
		}

		transaction.Status = "valid"
		transaction.OfferType = "On-demand"
		transaction.Sequence = strconv.Itoa(int(saida.Ticket))
		transaction.ParkingInfo = am.parking

		if transaction.ID.IsZero() {
			transaction, err = am.dbAcess.TransactionCollection.Create(ctx, transaction)
		} else {
			transaction.IsValid = true
			transaction, err = am.dbAcess.TransactionCollection.Update(ctx, transaction)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "error changing transaction: [%s]\n", err.Error())
		}
	}

	return nil
}
