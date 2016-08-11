package main

import "strings"

func TranslateKey(identificator string) string {
	replacer := strings.NewReplacer(
		"druh", "class",
		"vyrobce", "manufacturer",
		"typ", "type",
		"barva", "color",
		"spz", "regno",
		"mpz", "rpw",
		"motor", "engine",
		"rokvyroby", "productionyear",
		"nahlaseno", "stolendate",
	)
	return replacer.Replace(identificator)
}

func StandardizeDate(date string) string {
	replacer := strings.NewReplacer(
		"ledna", "1.",
		"února", "2.",
		"března", "3.",
		"dubna", "4.",
		"května", "5.",
		"června", "6.",
		"července", "7.",
		"srpna", "8.",
		"září", "9.",
		"října", "10.",
		"listopadu", "11.",
		"prosince", "12.",
		" ", "")
	return replacer.Replace(date)
}
