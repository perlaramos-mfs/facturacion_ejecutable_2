package Schemas

import (
	"encoding/xml"
)

type BillingEvents []LoanTransaction

type LoanTransaction struct {
	EventType       string          `json:"type"`
	LoanAccount     LoanAccount     `json:"loanAccount"`
	InvoiceGenerals InvoiceGenerals `json:"invoiceGenerals"`
	LineItems       []LineItem      `json:"lineItems"`
}

type Client struct {
	Id         string `json:"id"`
	EncodedKey string `json:"encodedKey"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
}

type LoanAccount struct {
	Id               string `json:"id"`
	EncodedKey       string `json:"encodedKey"`
	AccountHolderKey string `json:"accountHolderKey"`
}

type LoanAccountResponse struct {
	Id               string `json:"id"`
	EncodedKey       string `json:"encodedKey"`
	AccountHolderKey string `json:"accountHolderKey"`
}

type InvoiceGenerals struct {
	Amount     float64 `json:"amount"`
	EncodedKey string  `json:"encodedKey"`
	Timestamp  string  `json:"timestamp"`
	Condition  string  `json:"condition"`
}

type LineItem struct {
	Item          string  `json:"item"`
	TotalMambu    float64 `json:"totalMambu"`
	SubtotalMambu float64 `json:"subtotalMambu"`
	TaxMambu      float64 `json:"taxMambu"`
	TaxRate       float64 `json:"taxRate"`
	TotalSET      float64 `json:"totalSET"`
	SubtotalSET   float64 `json:"subtotalSET"`
	TaxSET        float64 `json:"taxSET"`
}

type LoanDynamo struct {
	LoanID        string `json:"LoanID"`
	Client        string `json:"Client"`
	EncodedKey    string `json:"EncodedKey"`
	LoanAccountId string `json:"LoanAccountId"`
}

type ClientDynamo struct {
	FirstName      string `json:"firstName"`      //client - FirstName
	LastName       string `json:"lastName"`       //client - LastName
	DocumentNumber string `json:"documentNumber"` //client - DocumentNumber
	DocumentType   string `json:"documentType"`   //client - DocumentType
}

type RUC struct {
	TaxpayerId string `json:"taxpayerId"`
	Ruc        string `json:"ruc"`
	RucDigit   string `json:"rucDigit"`
	PersonType string `json:"personType"`
}

type BillingData struct {
	RucEmisor                 string `json:"rucEmisor"`                 //RUC del Emisor
	VerificationDigitEmisor   string `json:"verificationDigitEmisor"`   //Digito verificador del emisor
	BillingType               string `json:"billingType"`               // Tipo de Factura: 1=Factura Electronica, 2=Factura electrónica de exportación. 3=Factura electrónica de importación, 4=Autofactura electrónica, 5=Nota de crédito electrónica, 6=Nota de débito electrónica, 7=Nota de remisión electrónica, 8=Comprobante de retención electrónico
	OperationType             string `json:"operationType"`             //Tipo Operación: 1=B2B, 2=B2C, 3=B2G
	TransactionType           string `json:"transactionType"`           // Tipo de transacción para el emisor: 1=Venta de mercadería, 2=Prestación de servicios, 3=Mixto (Venta de mercadería y servicios), 4=Venta de activo fijo, 5=Venta de divisas, 6=Compra de divisas,7=Promociones o entrega de muestras, 8=Donaciones, 9=Anticipos, 10=Compra de productos, 11=Compra de servicios, 12=Venta de crédito fiscal
	TaxType                   string `json:"taxType"`                   // Tipo Impuesto Consumo: 1=IVA, 2=ISC
	PresenceCode              string `json:"presenceCode"`              // Código indicador de presencia: 1=Operación presencial, 2=Operación por internet, 3=Operación telemarketing, 4=Venta a domicilio
	PaymentType               string `json:"paymentType"`               //Tipo de pago: 1=Efectivo, 2=Cheque, 3=Tarjeta de Crédito, 4=Tarjeta de Débito, 5=Transferencia, 6=Giros, 7=Billetera Electrónica, 8=Tarjetas Empresariales, 9=Vales
	DenominationCard          string `json:"denominationCard"`          //Denominación TC: 1=Visa, 2=Mastercard, 3=American Express, 4=Maestro, 5=Panal, 6=Cabal, 9=Otro
	ProcessingPaymentType     string `json:"processingPaymentType"`     //Forma de procesamiento de pago. Obligatorio si iTiPago = 3. Opcional1 POS2 Pago Electrónico (Ejemplo: compras por Internet)
	TimbNumber                string `json:"timbNumber"`                // Número de timbrado
	EstablishmentNumber       string `json:"establishmentPointNumber"`  // Establecimiento
	ExpeditionPointNumber     string `json:"expeditionPointNumber"`     // Punto de expedición
	CardProcessorBusinessName string `json:"cardProcessorBusinessName"` //Razón social de la procesadora de tarjeta
	ReturnKuDE                string `json:"returnKuDE"`
	ReturnSignedXML           string `json:"returnSignedXML"`
	TemplateKuDE              string `json:"templateKuDE"`
	CalculationsValidation    string `json:"calculationsValidation"`
	BillingCycle              string `json:"billingCycle"`
	ReentryForce              string `json:"reentryForce"`
	NotifyUpdateStatus        string `json:"notifyUpdateStatus"`
	URLCallback               string `json:"URLCallback"`
}

type TransactionForBilling struct {
	LoanTransaction
	ClientDynamo
	BillingData

	RUC
}

type XMLForBilling struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	Text    string   `xml:",chardata"`
	Soapenv string   `xml:"xmlns:soapenv,attr"`
	Ws      string   `xml:"xmlns:ws,attr"`
	Header  string   `xml:"soapenv:Header"`
	Body    Body     `xml:"soapenv:Body"`
}

type Body struct {
	Text                string              `xml:",chardata"`
	ProcesarLoteRequest ProcesarLoteRequest `xml:"ws:procesarLoteRequest"`
}

type ProcesarLoteRequest struct {
	Text              string              `xml:",chardata"`
	ProcesarDocumento []ProcesarDocumento `xml:"procesarDocumento"`
}

type ProcesarDocumento struct {
	Text                    string                  `xml:",chardata"`
	Xmlns                   string                  `xml:"xmlns,attr"`
	RDE                     RDE                     `xml:"rDE"`
	ParametrosProcesamiento ParametrosProcesamiento `xml:"parametrosProcesamiento"`
}

type RDE struct {
	Text string `xml:",chardata"`
	DE   DE     `xml:"DE"`
}

type DE struct {
	Text        string      `xml:",chardata"`
	GOpeDE      GOpeDE      `xml:"gOpeDE"`
	GTimb       GTimb       `xml:"gTimb"`
	GDatGralOpe GDatGralOpe `xml:"gDatGralOpe"`
	GDtipDE     GDtipDE     `xml:"gDtipDE"`
	GTotSub     GTotSub     `xml:"gTotSub"`
}

type GOpeDE struct {
	Text    string `xml:",chardata"`
	ITipEmi string `xml:"iTipEmi"`
}

type GTimb struct {
	Text    string `xml:",chardata"`
	ITiDE   string `xml:"iTiDE"`
	DNumTim string `xml:"dNumTim"`
	DEst    string `xml:"dEst"`
	DPunExp string `xml:"dPunExp"`
	DNumDoc string `xml:"dNumDoc"`
}

type GDatGralOpe struct {
	Text     string  `xml:",chardata"`
	DFeEmiDE string  `xml:"dFeEmiDE"`
	GOpeCom  GOpeCom `xml:"gOpeCom"`
	GEmis    GEmis   `xml:"gEmis"`
	GDatRec  GDatRec `xml:"gDatRec"`
}

type GOpeCom struct {
	Text     string `xml:",chardata"`
	ITipTra  string `xml:"iTipTra"`
	ITImp    string `xml:"iTImp"`
	CMoneOpe string `xml:"cMoneOpe"`
}

type GEmis struct {
	Text   string `xml:",chardata"`
	DRucEm string `xml:"dRucEm"`
	DDVEmi string `xml:"dDVEmi"`
}

type GDtipDE struct {
	Text     string    `xml:",chardata"`
	GCamFE   GCamFE    `xml:"gCamFE"`
	GCamCond GCamCond  `xml:"gCamCond"`
	GCamItem GCamItem  `xml:"gCamItem"`
	GCamEsp  *GCamEsp  //Factura crédito B2C Contribuyente y no contribuyente
	GCamNCDE *GCamNCDE //Nota de crédito B2C Contribuyente y no contribuyente
}

type GCamFE struct {
	Text     string `xml:",chardata"`
	IIndPres string `xml:"iIndPres"`
}

type GCamNCDE struct {
	XMLName xml.Name `xml:"gCamNCDE,omitempty"`
	Text    string   `xml:",chardata"`
	IMotEmi string   `xml:"iMotEmi,omitempty"`
}

type GCamCond struct {
	Text       string     `xml:",chardata"`
	ICondOpe   string     `xml:"iCondOpe"`
	GPaConEIni GPaConEIni `xml:"gPaConEIni"`
	GPagCred   *GPagCred
}

type GCamEsp struct {
	XMLName  xml.Name `xml:"gCamEsp,omitempty"`
	Text     string   `xml:",chardata"`
	GGrupAdi *GGrupAdi
}

type GGrupAdi struct {
	XMLName   xml.Name `xml:"gGrupAdi,omitempty"`
	Text      string   `xml:",chardata"`
	DCiclo    string   `xml:"dCiclo,omitempty"`
	DFecIniC  string   `xml:"dFecIniC,omitempty"`
	DFecFinC  string   `xml:"dFecFinC,omitempty"`
	DVencPag  string   `xml:"dVencPag,omitempty"`
	DContrato string   `xml:"dContrato,omitempty"`
}

type GPagCred struct {
	XMLName   xml.Name `xml:"gPagCred,omitempty"`
	ICondCred string   `xml:"iCondCred"`
	DPlazoCre string   `xml:"dPlazoCre"`
}

type GPaConEIni struct {
	Text       string  `xml:",chardata"`
	ITiPago    string  `xml:"iTiPago"`
	DMonTiPag  float64 `xml:"dMonTiPag"`
	CMoneTiPag string  `xml:"cMoneTiPag"`
	GPagTarCD  *GPagTarCD
}

type GPagTarCD struct {
	XMLName   xml.Name `xml:"gPagTarCD,omitempty"`
	IDenTarj  string   `xml:"ws:iDenTarj,omitempty"`
	DRSProTar string   `xml:"ws:dRSProTar,omitempty"`
	IForProPa string   `xml:"ws:iForProPa,omitempty"`
}

type GTotSub struct {
	Text         string  `xml:",chardata"`
	DSubExe      float64 `xml:"dSubExe"`
	DSub10       float64 `xml:"dSub10"`
	DTotOpe      float64 `xml:"dTotOpe"`
	DTotDesc     float64 `xml:"dTotDesc"`
	DIVA10       float64 `xml:"dIVA10"`
	DLiqTotIVA10 float64 `xml:"dLiqTotIVA10"`
	DTotIVA      float64 `xml:"dTotIVA"`
	DBaseGrav10  float64 `xml:"dBaseGrav10"`
	DTBasGraIVA  float64 `xml:"dTBasGraIVA"`
}

type ParametrosProcesamiento struct {
	Text                         string      `xml:",chardata"`
	RetornarKuDE                 string      `xml:"retornarKuDE"`
	RetornarXmlFirmado           string      `xml:"retornarXmlFirmado"`
	TemplateKuDE                 string      `xml:"templateKuDE"`
	ValidarCalculos              string      `xml:"validarCalculos"`
	CicloFacturacion             string      `xml:"cicloFacturacion"`
	ForzarReingreso              string      `xml:"forzarReingreso"`
	NotificarActualizacionEstado string      `xml:"notificarActualizacionEstado"`
	URLNotificacion              string      `xml:"urlNotificacion"`
	DocAsociado                  DocAsociado `xml:"docAsociado,omitempty"`
}

type DocAsociado struct { // PARA NOTA DE CRÉDITO CONTRIBUYENTE
	Text            string `xml:",chardata"`
	Tipo            string `xml:"tipo,omitempty"`
	Timbrado        string `xml:"timbrado,omitempty"`
	Establecimiento string `xml:"establecimiento,omitempty"`
	PuntoExpedicion string `xml:"puntoExpedicion,omitempty"`
	Numero          string `xml:"numero,omitempty"`
}

type GDatRec struct {
	Text       string `xml:",chardata"`
	INatRec    string `xml:"iNatRec"`
	ITiOpe     string `xml:"iTiOpe"`
	CPaisRec   string `xml:"cPaisRec"`
	ITiContRec string `xml:"iTiContRec"`
	DRucRec    string `xml:"dRucRec"`
	DDVRec     string `xml:"dDVRec"`
	DNomRec    string `xml:"dNomRec"`
	ITipIDRec  string `xml:"iTipIDRec,omitempty"` //Factura contado
	DNumIDRec  string `xml:"dNumIDRec,omitempty"` //Factura contado
}

type GCamItem struct {
	Text        string     `xml:",chardata"`
	DCodInt     string     `xml:"dCodInt"`
	DDesProSer  string     `xml:"dDesProSer"`
	DCantProSer string     `xml:"dCantProSer"`
	GValorItem  GValorItem `xml:"gValorItem"`
	GCamIVA     GCamIVA    `xml:"gCamIVA"`
}

type GValorItem struct {
	Text            string          `xml:",chardata"`
	DPUniProSer     float64         `xml:"dPUniProSer"`
	DTotBruOpeItem  float64         `xml:"dTotBruOpeItem"`
	GValorRestaItem GValorRestaItem `xml:"gValorRestaItem"`
}

type GValorRestaItem struct {
	Text            string  `xml:",chardata"`
	DDescItem       float64 `xml:"dDescItem"`
	DAntGloPreUniIt float64 `xml:"dAntGloPreUniIt"`
	DTotOpeItem     float64 `xml:"dTotOpeItem"`
}

type GCamIVA struct {
	Text        string  `xml:",chardata"`
	IAfecIVA    string  `xml:"iAfecIVA"`
	DTasaIVA    float64 `xml:"dTasaIVA"`
	DBasGravIVA float64 `xml:"dBasGravIVA"`
	DLiqIVAItem float64 `xml:"dLiqIVAItem"`
}
