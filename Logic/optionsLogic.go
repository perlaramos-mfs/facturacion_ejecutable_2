package Logic

func SetTipoContribuyente(value string) string {
	var result string
	switch value {
	case "F": //Física
		result = "1"
		break
	case "J": //Jurídica
		result = "2"
		break
	default:
		break
	}
	return result
}

func SetCondition(value string) string {
	var result string
	switch value {
	case "contado":
		result = "1"
		break
	case "credito":
		result = "2"
		break
	case "nota_credito":
		result = "3"
		break
	default:
		break
	}
	return result
}

func SetNatReceptor(value bool) string {
	var result string
	switch value {
	case false:
		result = "2"
		break
	case true:
		result = "1"
		break
	default:
		break
	}
	return result
}

func SetCodeAndDescription(value string) (string, string) {
	var code string
	var description string
	switch value {
	case "adminFee":
		description = "Gastos administrativos"
		code = "REV0001"
		break
	case "interestActiveAccount":
		description = "Intereses de préstamo activo"
		code = "REV0002"
		break
	case "moratoryInterest":
		description = "Intereses moratorios"
		code = "REV0003"
		break
	case "penalty":
		description = "Intereses punitorios"
		code = "REV0004"
		break
	default:
		break
	}
	return code, description
}
