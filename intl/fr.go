package intl

import (
	"strings"
	"time"
)

var month = strings.NewReplacer(
	"Jan", "jan.",
	"Feb", "fév.",
	"Mar", "mars",
	"Apr", "avril",
	"May", "mai",
	"Jun", "juin",
	"Jul", "juil.",
	"Aug", "août",
	"September", "sept.",
	"Oct", "oct.",
	"Nov", "nov.",
	"Dec", "déc.",
)

var weekday = strings.NewReplacer(
	"Monday", "lundi",
	"Tuesday", "mardi",
	"Wednesday", "mercredi",
	"Thursday", "jeudi",
	"Friday", "vendredi",
	"Saturday", "samedi",
	"Sunday", "dimanche",
)

func (fr *Fr) FmtMMMEEEEd(date time.Time) string {
	return month.Replace(weekday.Replace(date.In(fr.loc).Format("le Monday 02 Jan à 15h04")))
}

func (fr *Fr) FmtyMd(date time.Time) string {
	return date.In(fr.loc).Format("02/01/2006")
}

type Fr struct {
	loc *time.Location
}

type Parser struct {
	Fr Fr
}

func NewIntlParser(region string, loc *time.Location) *Parser {
	switch region {
	case "Fr":
		return &Parser{
			Fr: Fr{loc: loc},
		}
	}
	return &Parser{}
}
