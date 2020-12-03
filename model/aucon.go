package model

import (
	"encoding/xml"
	"time"
)

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

type AuconPagamento struct {
	XMLName       xml.Name
	Data          time.Time `xml:"DATA" bson:"data"`
	Ticket        int64     `xml:"TICKET" bson:"ticket"`
	Tp            string    `xml:"TP" bson:"tp"`
	Matricula     int64     `xml:"MATRICULA" bson:"matricula"`
	TpPagamento   string    `xml:"TP_PAGTO" bson:"tp_pagto"`
	Convenio      string    `xml:"CONVENIO" bson:"convenio"`
	TempoCob      int64     `xml:"TEMPO_COB" bson:"tempo_cob"`
	Valor         float64   `xml:"VALOR" bson:"valor"`
	Debito        float64   `xml:"DEBITO" bson:"debito"`
	Desconto      float64   `xml:"DESCONTO" bson:"desconto"`
	Operador      string    `xml:"OPERADOR" bson:"operador"`
	Motivo        string    `xml:"MOTIVO" bson:"motivo"`
	Filial        int64     `xml:"FILIAL" bson:"filial"`
	QtSelos       int64     `xml:"QT_SELOS" bson:"qt_selos"`
	TpSelos       int64     `xml:"TP_SELOS" bson:"tp_selos"`
	Nomeconvenio  string    `xml:"NOMECONVENIO" bson:"nomeconvenio"`
	Fechamento    int64     `xml:"FECHAMENTO" bson:"fechamento"`
	ConvCodBarras int64     `xml:"CONVENIO_COD_BARRAS" bson:"convenio_cod_barras"`
}

type AuconCredenciado struct {
	XMLName    xml.Name
	MATRICULA  int64  `xml:"MATRICULA"`
	NOME       string `xml:"NOME"`
	ENDERECO   string `xml:"ENDERECO"`
	CIDADE     string `xml:"CIDADE"`
	BAIRRO     string `xml:"BAIRRO"`
	UF         string `xml:"UF"`
	CEP        string `xml:"CEP"`
	FONE       string `xml:"FONE"`
	FONE2      string `xml:"FONE2"`
	TP_PESSOA  string `xml:"TP_PESSOA"`
	CPF        string `xml:"CPF"`
	CNPJ       string `xml:"CNPJ"`
	DT_NASC    string `xml:"DT_NASC"`
	TIPO       string `xml:"TIPO"`
	CATEGORIA  string `xml:"CATEGORIA"`
	ESPECIAL   string `xml:"ESPECIAL"`
	VAGAS      string `xml:"VAGAS"`
	VALOR_MES  string `xml:"VALOR_MES"`
	TOLERANCIA string `xml:"TOLERANCIA"`
	INSCRICAO  string `xml:"INSCRICAO"`
	OBS        string `xml:"OBS"`
	BOX_OCUP   string `xml:"BOX_OCUP"`
	MVENC      string `xml:"MVENC"`
	ATIVO      string `xml:"ATIVO"`
	DT_COMUTA  string `xml:"DT_COMUTA"`
	CREDITO    string `xml:"CREDITO"`
	DEBITO     string `xml:"DEBITO"`
	FILIAL     string `xml:"FILIAL"`
}

type AuconSaidaEnvelope struct {
	XMLName xml.Name
	Body    AuconSaidaBody
}

type AuconSaidaBody struct {
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

type AuconPagamentoEnvelope struct {
	XMLName xml.Name
	Body    AuconPagamentoBody
}

type AuconPagamentoBody struct {
	XMLName     xml.Name
	GetResponse AuconGetPagamentosResponse `xml:"Get_PagamentosResponse"`
}

type AuconGetPagamentosResponse struct {
	XMLName   xml.Name
	GetResult AuconGetPagamentosResult `xml:"Get_PagamentosResult"`
}

type AuconGetPagamentosResult struct {
	XMLName  xml.Name
	Diffgram AuconGetPagamentosDiffgram `xml:"diffgram"`
}

type AuconGetPagamentosDiffgram struct {
	XMLName    xml.Name
	DocElement AuconGetPagamentosDocElement `xml:"DocumentElement"`
}

type AuconGetPagamentosDocElement struct {
	XMLName xml.Name
	Pagtos  []AuconPagamento `xml:"PAGTO"`
}

type AuconCredenciadoEnvelope struct {
	XMLName xml.Name
	Body    AuconCredenciadoBody
}

type AuconCredenciadoBody struct {
	XMLName     xml.Name
	GetResponse AuconGetCredenciadosResponse `xml:"Get_CredenciadosResponse"`
}

type AuconGetCredenciadosResponse struct {
	XMLName   xml.Name
	GetResult AuconGetCredenciadosResult `xml:"Get_CredenciadosResult"`
}

type AuconGetCredenciadosResult struct {
	XMLName      xml.Name
	Credenciados []AuconCredenciado `xml:"WS_CRED_LOCAL""`
}
