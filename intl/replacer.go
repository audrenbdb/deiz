package intl

import "strings"

var month = strings.NewReplacer(
	"January", "Janvier",
	"February", "Février",
	"Mars", "Mars",
	"April", "Avril",
	"May", "Mai",
	"June", "Juin",
	"July", "Juillet",
	"August", "Août",
	"September", "Septembre",
	"October", "Octobre",
	"November", "Novembre",
	"December", "Décembre",
)

var weekday = strings.NewReplacer(
	"Monday", "Lundi",
	"Tuesday", "Mardi",
	"Wednesday", "Mercredi",
	"Thursday", "Jeudi",
	"Friday", "Vendredi",
	"Saturday", "Samedi",
	"Sunday", "Dimanche",
)

func ToFrench(timeString string) string {
	return month.Replace(weekday.Replace(timeString))
}
