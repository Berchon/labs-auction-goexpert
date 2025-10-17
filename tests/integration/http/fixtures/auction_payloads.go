package fixtures

var (
	ValidAuction = map[string]interface{}{
		"product_name": "Mola maluca",
		"category":     "Brinquedo",
		"description":  "Você vai adorar",
		"condition":    1, // 1 = novo, 2 = usado, 3 = recondicionado
	}

	ValidAuction2 = map[string]interface{}{
		"product_name": "Skate Profissional",
		"category":     "Esporte",
		"description":  "Skate de madeira canadense com rolamento ABEC-9",
		"condition":    2, // 1 = novo, 2 = usado, 3 = recondicionado
	}

	ValidAuction3 = map[string]interface{}{
		"product_name": "Violão Clássico Yamaha C40",
		"category":     "Instrumentos Musicais",
		"description":  "Violão de nylon, perfeito para iniciantes e músicos experientes",
		"condition":    3, // 1 = novo, 2 = usado, 3 = recondicionado
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

	InvalidProductName = map[string]interface{}{
		"product_name": "",
		"category":     "Brinquedo",
		"description":  "Você vai adorar esse brinquedo incrível",
		"condition":    1,
	}

	InvalidCategory = map[string]interface{}{
		"product_name": "Mola maluca",
		"category":     "AB",
		"description":  "Descrição válida e completa",
		"condition":    1,
	}

	InvalidDescriptionAndCondition = map[string]interface{}{
		"product_name": "Mola maluca",
		"category":     "Brinquedo",
		"description":  "Curto",
		"condition":    1,
	}

	ValidShortDescription = map[string]interface{}{
		"product_name": "Mola maluca",
		"category":     "Brinquedo",
		"description":  "Curto",
		"condition":    1,
	}

	MultipleInvalidFields = map[string]interface{}{
		"product_name": "A",
		"category":     "AB",
		"description":  "Curto",
		"condition":    1,
	}
)
