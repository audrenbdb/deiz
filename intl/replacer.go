package intl

import "strings"

var month = strings.NewReplacer(
	"Jan", "jan.",
	"Feb", "fév.",
	"Mar", "mars",
	"Apr", "avril",
	"May", "mai",
	"Jun", "juin",
	"Jul", "juil.",
	"Aug", "août",
	"Setember", "sept.",
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

func ToFrench(timeString string) string {
	return month.Replace(weekday.Replace(timeString))
}
