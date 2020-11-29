package model

import "encoding/xml"

type AuconSaida struct {
	XMLName     xml.Name
	DataEnt     string  `xml:"DATA_ENT" bson:"data_ent"`
	DataSai     string  `xml:"DATA_SAI" bson:"data_sai"`
	TpAcesso    string  `xml:"TP_ACESSO" bson:"tp_acesso"`
	Ticket      int64   `xml:"TICKET" bson:"ticket"`
	Card        string  `xml:"CARD" bson:"card"`
	Matricula   int64   `xml:"MATRICULA" bson:"matricula"`
	Placa       string  `xml:"PLACA" bson:"placa"`
	Modelo      string  `xml:"MODELO" bson:"modelo"`
	Cor         string  `xml:"COR" bson:"cor"`
	Tabela      int64   `xml:"TABELA" bson:"table"`
	Debito      float64 `xml:"DEBITO" bson:"debito"`
	MaquinaEnt  int64   `xml:"MAQUINA_ENT" bson:"maquina_ent"`
	MaquinaSai  int64   `xml:"MAQUINA_SAI" bson:"maquina_sai"`
	OperadorEnt string  `xml:"OPERADOR_ENT" bson:"operador_ent"`
	OperadorSai string  `xml:"OPERADOR_SAI" bson:"operador_sai"`
	Fx1         string  `xml:"FX1" bson:"fx1"`
	Fx2         string  `xml:"FX2" bson:"fx2"`
	Filial      int64   `xml:"FILIAL" bson:"filial"`
	TSND        string  `xml:"TSND" bson:"tnsd"`
	TSND2       string  `xml:"TSND2" bson:"tnsd2"`
	Clockset    int64   `xml:"CLOCKSET" bson:"clockset"`
	Moto        string  `xml:"MOTO" bson:"moto"`
}

type AuconEnvelope struct {
	XMLName xml.Name
	Body    AuconBody
}

type AuconBody struct {
	XMLName     xml.Name
	GetResponse AuconGetSaidasResponse `xml:"GetSaidasResponse"`
}

type AuconGetSaidasResponse struct {
	XMLName   xml.Name
	GetResult AuconGetSaidasResult `xml:"GetSaidasResult"`
}

type AuconGetSaidasResult struct {
	XMLName xml.Name
	Saidas  []AuconSaida `xml:"SAIDA""`
}
