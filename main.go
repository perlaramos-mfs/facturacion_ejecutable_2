package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	Dynamo "main/DynamoDB"
	Logic "main/Logic"
	"main/Schemas"
	"time"
)

const SOAPENV = "http://schemas.xmlsoap.org/soap/envelope/"
const WS = "http://ws.efactura.isaltda.py/"

func main() {
	events, err := ioutil.ReadFile("events.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	data := Schemas.BillingEvents{}
	_ = json.Unmarshal(events, &data)

	billingData := Schemas.BillingData{ //Información común para todas las facturas
		RucEmisor:                 "80021477", //dRucEm: RUC del emisor
		VerificationDigitEmisor:   "3",        //dDVEmi: Digito verificador del emisor
		BillingType:               "1",        //iTiDE: Tipo de Documento Electrónico: 1 Factura Electronica, 2 Factura electrónica de exportación 3 Factura electrónica de importación 4 Autofactura electrónica 5 Nota de crédito electrónica 6 Nota de débito electrónica 7 Nota de remisión electrónica 8 Comprobante de retención electrónico
		OperationType:             "2",        //iTiOpe: Tipo Operación: 1=B2B, 2=B2C, 3=B2G
		TransactionType:           "2",        //iTipTra: Tipo de transacción: 1= Venta de mercadería2= Prestación de servicios3= Mixto (Venta de mercadería y servicios)4= Venta de activo fijo5= Venta de divisas6= Compra de divisas7= Promociones o entrega de muestras8= Donaciones9= Anticipos10= Compra de productos11= Compra de servicios12= Venta de crédito fiscal
		TaxType:                   "1",        //iTImp: Tipo Impuesto afectado Consumo1 IVA, 2 ISC
		PresenceCode:              "2",        //iIndPres: Código indicador de presencia 1= Operación presencial 2= Operación por internet 3= Operación telemarketing 4= Venta a domicilio
		PaymentType:               "1",        //iTiPago: Tipo de pago: 1:Efectivo, 2:Cheque, 3:Tarjeta de Crédito, 4:Tarjeta de Débito, 5:Transferencia, 6:Giros, 7:Billetera Electrónica, 8:Tarjetas Empresariales, 9:Vales
		DenominationCard:          "2",        //iDenTarj: Denominación TC: 1=Visa, 2=Mastercard 3 American Express4 Maestro5 Panal6 Cabal9 Otro
		ProcessingPaymentType:     "2",        //iForProPa: Forma de procesamiento de pago. Obligatorio si iTiPago = 3. Opcional1 POS2 Pago Electrónico (Ejemplo: compras por Internet)
		TimbNumber:                "12559367", // Número de timbrado *CAMBIARÁ POR AMBIENTE
		EstablishmentNumber:       "001",      // Establecimiento *CAMBIARÁ POR AMBIENTE
		ExpeditionPointNumber:     "001",      // Punto de expedición *CAMBIARÁ POR AMBIENTE
		CardProcessorBusinessName: "aaaa",     //dRSProTarRazón social de la procesadora de tarjeta // ToDo: Saber qué valor mandar
		ReturnKuDE:                "true",     //retornarKuDE ToDo: Saber cuándo cambia valor
		ReturnSignedXML:           "true",     //retornarXmlFirmado ToDo: Saber cuándo cambia valor
		TemplateKuDE:              "1",        //templateKuDEToDo: Saber cuándo cambia valor
		CalculationsValidation:    "false",    //validarCalculos ToDo: Saber cuándo cambia valor
		BillingCycle:              "20210707", //cicloFacturacion ToDo: Saber cuándo cambia valor
		ReentryForce:              "",         //forzarReingreso ToDo: Saber qué es y cuándo cambia valor
		NotifyUpdateStatus:        "false",    //notificarActualizacionEstado ToDo: Saber cuándo cambia valor
		URLCallback:               "",
	}

	request := Schemas.XMLForBilling{
		Soapenv: SOAPENV,
		Ws:      WS,
	}

	var procesarDocumento []Schemas.ProcesarDocumento

	for i := 0; i < len(data); i++ {
		loanTransaction := data[i]

		loanDynamo := Dynamo.QueryDB(loanTransaction.LoanAccount.AccountHolderKey, "Loan", "EncodedKey-index", "EncodedKey = :EncodedKey", ":EncodedKey") //MANDAR infoFromMambu.loanAccount.accountHolderKey
		clientDynamo := Dynamo.GetClient(loanDynamo.Client, "Client", "ClientID = :Client", ":Client")                                                    //mandar loanDynamo.Client
		rucRegisterDynamo := Dynamo.QueryRUC(clientDynamo.DocumentNumber, "rucs-dev", "taxpayerId")                                                       //MANDAR  clientDynamo.DocumentNumber

		info := joinInfo(loanTransaction, clientDynamo, billingData, rucRegisterDynamo)
		res := conversionToXML(info)

		procesarDocumento = append(procesarDocumento, res)
		//request.Body.ProcesarLoteRequest.ProcesarDocumento = res

	}

	request.Body.ProcesarLoteRequest.ProcesarDocumento = procesarDocumento

	v, _ := xml.MarshalIndent(request, "", " ")

	fmt.Println(string(v))
}

func joinInfo(loanTransaction Schemas.LoanTransaction, infoClientDynamo Schemas.ClientDynamo, billingData Schemas.BillingData, rucRegisterDynamo Schemas.RUC) Schemas.TransactionForBilling {
	var transactionForBilling Schemas.TransactionForBilling

	transactionForBilling.LoanTransaction = loanTransaction
	transactionForBilling.ClientDynamo = infoClientDynamo
	transactionForBilling.BillingData = billingData
	transactionForBilling.RUC = rucRegisterDynamo

	return transactionForBilling
}

func conversionToXML(transactionForBilling Schemas.TransactionForBilling) Schemas.ProcesarDocumento {
	//ToDo Llamar a Redis usando transactionForBilling.LoanTransaction.InvoiceGenerals.EncodedKey que es el EncodedKey de la transacción que se usará para ver en Redis si ya se facturó esta transacción o no
	dNumDoc := "898989" // Número de factura ToDo: Si no se ha facturado, colocar valor obtenido de Redis, consecutivo generado

	//Constantes:
	const iTipEmi = "1"  //Tipo Emisión: 1=Normal, 2=Contingencia
	const iAfecIVA = "1" //Forma afectación IVA1 Gravado IVA2 Exonerado (Art. 83-Ley 125/91)3 Exento4 Gravado parcial (Grav-Exento)

	const dCantProSer = "1"  //Cantidad de producto o servicio
	const cMoneOpe = "PYG"   // Moneda de operación
	const cPaisRec = "PRY"   //Código de pais PRY
	const cMoneTiPag = "PYG" // Moneda

	procesarDocumento := Schemas.ProcesarDocumento{}

	procesarDocumento.Xmlns = WS
	procesarDocumento.RDE.DE.GOpeDE.ITipEmi = iTipEmi

	iTide := transactionForBilling.BillingData.BillingType

	gTimb := Schemas.GTimb{
		ITiDE:   iTide,
		DNumTim: transactionForBilling.BillingData.TimbNumber,
		DEst:    transactionForBilling.BillingData.EstablishmentNumber,
		DPunExp: transactionForBilling.BillingData.ExpeditionPointNumber,
		DNumDoc: dNumDoc,
	}
	procesarDocumento.RDE.DE.GTimb = gTimb
	procesarDocumento.RDE.DE.GDatGralOpe.DFeEmiDE = GetTimeFormat()

	gOpeCom := Schemas.GOpeCom{
		ITImp:    transactionForBilling.BillingData.TaxType,
		CMoneOpe: cMoneOpe,
	}

	procesarDocumento.RDE.DE.GDatGralOpe.GOpeCom = gOpeCom

	gEmis := Schemas.GEmis{
		DRucEm: transactionForBilling.BillingData.RucEmisor,
		DDVEmi: transactionForBilling.BillingData.VerificationDigitEmisor,
	}

	procesarDocumento.RDE.DE.GDatGralOpe.GEmis = gEmis

	dCodInt, dDesProSer := Logic.SetCodeAndDescription(transactionForBilling.LineItems[0].Item) //Descripción product o o servicio

	var iCondOpe string                                                                            //Condición de operación: 1=Contado, 2=Crédito
	iCondOpe = Logic.SetCondition(transactionForBilling.LoanTransaction.InvoiceGenerals.Condition) //Condición de operación: 1=Contado, 2=Crédito
	iTiOpe := transactionForBilling.BillingData.OperationType                                      //Tipo de Operación 1=B2B, 2=B2C, 3=B2G

	gDatRec := Schemas.GDatRec{
		ITiOpe:   iTiOpe,
		CPaisRec: cPaisRec,
		DNomRec:  transactionForBilling.FirstName + " " + transactionForBilling.LastName,
	}

	var iNatRec string
	if transactionForBilling.RUC.Ruc != "" { //Si es contribuyente se agrega ruc y dígito verificador
		iNatRec = Logic.SetNatReceptor(true)                                           //Naturaleza del receptor: 1=Contribuyente
		iTiContRec := Logic.SetTipoContribuyente(transactionForBilling.RUC.PersonType) //Tipo Contribuyente receptor (si iNatRec == 1)
		gDatRec.ITiContRec = iTiContRec
		gDatRec.DRucRec = transactionForBilling.RUC.Ruc
		gDatRec.DDVRec = transactionForBilling.RUC.RucDigit
	} else { //si no es contribuyente, se agrega tipo de documento de identidad y número de identidad
		iNatRec = Logic.SetNatReceptor(false)                                 //Naturaleza del receptor: 2=No Contribuyente
		gDatRec.ITipIDRec = transactionForBilling.ClientDynamo.DocumentType   //tipo de documento
		gDatRec.DNumIDRec = transactionForBilling.ClientDynamo.DocumentNumber //numero de documento
	}
	gDatRec.INatRec = iNatRec
	procesarDocumento.RDE.DE.GDatGralOpe.GDatRec = gDatRec

	if iTide == "1" || iTide == "4" {
		procesarDocumento.RDE.DE.GDatGralOpe.GOpeCom.ITipTra = transactionForBilling.BillingData.TransactionType
		procesarDocumento.RDE.DE.GDtipDE.GCamCond.ICondOpe = iCondOpe
	}

	if iTide == "1" {
		procesarDocumento.RDE.DE.GDtipDE.GCamFE.IIndPres = transactionForBilling.BillingData.PresenceCode
	}

	if iCondOpe == "2" { //si la condición de operación es 2 (Crédito), se especifican los plazos
		gPagCred := Schemas.GPagCred{
			ICondCred: "1",       //Condición de la operación a crédito: 1=Plazo, 2=Cuota ES UN VALOR FIJO
			DPlazoCre: "30 días", //Descripción de la condición de la operación a crédito ES UN VALOR FIJO
		}
		procesarDocumento.RDE.DE.GDtipDE.GCamCond.GPagCred = &gPagCred
	}

	iTiPago := transactionForBilling.BillingData.PaymentType
	gPaConEIni := Schemas.GPaConEIni{
		ITiPago:    iTiPago,
		DMonTiPag:  transactionForBilling.LineItems[0].TotalSET,
		CMoneTiPag: cMoneTiPag,
	}
	procesarDocumento.RDE.DE.GDtipDE.GCamCond.GPaConEIni = gPaConEIni

	if iTiPago == "3" || iTiPago == "4" {
		gPagTarCD := Schemas.GPagTarCD{
			IDenTarj:  transactionForBilling.BillingData.DenominationCard,
			DRSProTar: transactionForBilling.BillingData.CardProcessorBusinessName,
			IForProPa: transactionForBilling.BillingData.ProcessingPaymentType,
		}
		procesarDocumento.RDE.DE.GDtipDE.GCamCond.GPaConEIni.GPagTarCD = &gPagTarCD
	}

	gValorRestaItem := Schemas.GValorRestaItem{
		DDescItem:       0,                                           //ToDo: Saber valor
		DAntGloPreUniIt: 0,                                           //ToDo: Saber valor
		DTotOpeItem:     transactionForBilling.LineItems[0].TotalSET, //ToDo: Saber si sí mandamos el monto
	}

	gValorItem := Schemas.GValorItem{
		DPUniProSer:     transactionForBilling.LineItems[0].TotalSET, //ToDo: Saber si sí mandamos el monto
		DTotBruOpeItem:  transactionForBilling.LineItems[0].TotalSET, //ToDo: Saber si sí mandamos el monto
		GValorRestaItem: gValorRestaItem,
	}

	gCamIVA := Schemas.GCamIVA{
		IAfecIVA:    iAfecIVA,
		DTasaIVA:    transactionForBilling.LineItems[0].TaxRate,
		DBasGravIVA: transactionForBilling.LineItems[0].SubtotalSET,
		DLiqIVAItem: transactionForBilling.LineItems[0].TaxSET,
	}

	gCamItem := Schemas.GCamItem{
		DCodInt:     dCodInt,
		DDesProSer:  dDesProSer,
		DCantProSer: dCantProSer,
		GValorItem:  gValorItem,
		GCamIVA:     gCamIVA,
	}

	//factCredB2CContribuyente --> operationType B2C = 2, Contribuyente = true, Condition Operation crédito = 2
	/*if transactionForBilling.BillingData.OperationType == "2" && transactionForBilling.RUC.Ruc != "" && iCondOpe == "2" {
		dCiclo := "20220307"     //ToDo: Saber valor
		dFecIniC := "2022-03-07" //ToDo: Saber valor
		dFecFinC := "2022-04-06" //ToDo: Saber valor
		dVencPag := "2022-05-11" //ToDo: Saber valor
		dContrato := "4B9217131" //ToDo: Saber valor

		gGrupAdi := Schemas.GGrupAdi{
			DCiclo:    dCiclo,
			DFecIniC:  dFecIniC,
			DFecFinC:  dFecFinC,
			DVencPag:  dVencPag,
			DContrato: dContrato,
		}
		request.Body.ProcesarLoteRequest.ProcesarDocumento.RDE.DE.GDtipDE.GCamEsp.GGrupAdi = gGrupAdi
	}*/ //SE COMENTA PORQUE POR AHORA NO SE UTILIZARÁ ESTE GRUPO

	procesarDocumento.RDE.DE.GDtipDE.GCamItem = gCamItem

	gTotSub := Schemas.GTotSub{
		DSubExe:      0,                                              //ToDo: Saber valor
		DSub10:       transactionForBilling.LineItems[0].TotalSET,    //ToDo: Saber valor
		DTotOpe:      transactionForBilling.LineItems[0].TotalSET,    //ToDo: Saber valor
		DTotDesc:     0,                                              //ToDo: Saber valor
		DIVA10:       transactionForBilling.LineItems[0].TaxSET,      //ToDo: Saber valor
		DLiqTotIVA10: 0,                                              //ToDo: Saber valor
		DTotIVA:      transactionForBilling.LineItems[0].TaxSET,      //ToDo: Saber valor
		DBaseGrav10:  transactionForBilling.LineItems[0].SubtotalSET, //ToDo: Saber valor
		DTBasGraIVA:  transactionForBilling.LineItems[0].SubtotalSET, //ToDo: Saber valor
	}
	procesarDocumento.RDE.DE.GTotSub = gTotSub

	parametrosProcesamiento := Schemas.ParametrosProcesamiento{
		RetornarKuDE:                 transactionForBilling.BillingData.ReturnKuDE,
		RetornarXmlFirmado:           transactionForBilling.BillingData.ReturnSignedXML,
		TemplateKuDE:                 transactionForBilling.BillingData.TemplateKuDE,
		ValidarCalculos:              transactionForBilling.BillingData.CalculationsValidation,
		CicloFacturacion:             transactionForBilling.BillingData.BillingCycle,
		ForzarReingreso:              transactionForBilling.BillingData.ReentryForce,
		NotificarActualizacionEstado: transactionForBilling.BillingData.NotifyUpdateStatus,
		URLNotificacion:              transactionForBilling.BillingData.URLCallback,
	}
	procesarDocumento.ParametrosProcesamiento = parametrosProcesamiento

	if iTide == "5" || iTide == "6" { //En caso de Nota de Crédito
		gCamNCDE := Schemas.GCamNCDE{
			IMotEmi: "", //ToDo: Ver qué valor irá aquí cuando se implementen notas de crédito
		}
		procesarDocumento.RDE.DE.GDtipDE.GCamNCDE = &gCamNCDE

		tipo := "1"            //Tipo DE: 1=Factura Electronica, 2=Factura electrónica de exportación, 3=Factura electrónica de importación, 4=Autofactura electrónica, 5=Nota de crédito electrónica, 6=Nota de débito electrónica, 7=Nota de remisión electrónica,8=Comprobante de retención electrónico ToDo: Saber cuándo cambia valor cuando se implementen notas de crédito
		timbrado := "12559367" //Timbrado del documento asociado //ToDo: Saber valor cuando se implementen notas de crédito
		establecimiento := "1" //Establecimiento del documento asociado //ToDo: Saber cuándo cambia valor cuando se implementen notas de crédito
		puntoExpedicion := "1" //Punto de Expedición del documento asociado //ToDo: Saber cuándo cambia valor cuando se implementen notas de crédito
		numero := "14"         //Numero del documento asociado //ToDo: Saber valor cuando se implementen notas de crédito

		docAsociado := Schemas.DocAsociado{ //PARA NOTA DE CRÉDITO CONTRIBUYENTE
			Tipo:            tipo,
			Timbrado:        timbrado,
			Establecimiento: establecimiento,
			PuntoExpedicion: puntoExpedicion,
			Numero:          numero,
		}
		procesarDocumento.ParametrosProcesamiento.DocAsociado = docAsociado
	}

	//v, _ := xml.MarshalIndent(procesarDocumento, "", " ")
	return procesarDocumento
}

func GetTimeFormat() string {
	loc, _ := time.LoadLocation("America/Asuncion")
	layout := "2006-01-02T15:04:05"

	//set timezone
	now := time.Now().In(loc).Format(layout)
	return now
}
