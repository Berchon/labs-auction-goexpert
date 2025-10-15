package fixtures

var (
	ValidAuction = map[string]interface{}{
		"product_name": "Mola maluca",
		"category":     "Brinquedo",
		"description":  "Você vai adorar",
		"condition":    1,
	}

	MissingField = map[string]interface{}{
		"category":    "Brinquedo",
		"description": "Você vai adorar",
		"condition":   1,
	}

	InvalidType = map[string]interface{}{
		"product_name": "Mola maluca",
		"category":     "Brinquedo",
		"description":  "Você vai adorar",
		"condition":    "usado",
	}

	InvalidCondition = map[string]interface{}{
		"product_name": "Mola maluca",
		"category":     "Brinquedo",
		"description":  "Você vai adorar",
		"condition":    99,
	}
)
