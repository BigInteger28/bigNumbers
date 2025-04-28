package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strings"
)

var (
	x1000     = []string{"", "k", "M", "G", "T", "P", "E", "Z", "Y", "R", "Q", "X11", "X12", "X13", "X14", "X15", "X16", "X17", "X18", "X19", "X20", "X21", "X22", "X23", "X24", "X25", "X26", "X27", "X28", "X29", "X30", "X31", "X32", "X33", "X34", "X35", "X36", "X37", "X38", "X39", "X40", "X41"}
	x1000text = []string{"", " Duizend ", " Miljoen ", " Miljard ", " Biljoen ", " Biljard ", " Triljoen ", " Triljard ", " Quadriljoen ", " Quadriljard ", " Quintiljoen ", " Quintiljard ", " Sextiljoen ", " Sextiljard ", " Septiljoen ", " Septiljard ", " Octiljoen ", " Octiljard ", " Noniljoen ", " Noniljard ", " Deciljoen ", " Deciljard ", " Undeciljoen ", " Undeciljard ", " Duodeciljoen ", " Duodeciljard ", " Tredeciljoen ", " Tredeciljard ", " Quattuordeciljoen ", " Quattuordeciljard ", " Quindeciljoen ", " Quindeciljard ", " Sedeciljoen ", " Sedeciljard ", " Septendeciljoen ", " Septendeciljard ", " Octodeciljoen ", " Octodeciljard ", " Novemdeciljoen ", " Novemdeciljard ", " Vigintiljoen ", " Vigintiljard "}
)

const maxIndex = 41 // Maximale index om fouten te voorkomen

type Number struct {
	amount        big.Int
	dig           big.Int
	x1000         int
	originalInput string
	hasSuffix     bool // Houdt bij of er een x1000-suffix in de invoer zat
}

// Geeft de maximale waarde van twee integers terug
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Formatteert een groot getal met spaties tussen elk blok van 3 cijfers
func formatBigNumber(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	parts := []string{}
	for i := n; i > 0; i -= 3 {
		start := max(0, i-3)
		parts = append([]string{s[start:i]}, parts...)
	}
	return strings.Join(parts, " ")
}

// Bepaalt de juiste x1000-index en kapt af tot 3 significante cijfers
func setX1000(value *Number) {
	strAmount := value.amount.String()
	length := len(strAmount)

	// Controleer op ongeldige input
	if length == 0 {
		value.x1000 = 0
		value.dig.SetInt64(0)
		return
	}

	// Bepaal index (elke 3 extra cijfers betekent +1 in x1000-index)
	if strAmount[0] == '-' {
		value.x1000 = (length - 2) / 3
	} else {
		value.x1000 = (length - 1) / 3
	}

	// Voorkom dat x1000 buiten de limiet gaat
	if value.x1000 > maxIndex {
		value.x1000 = maxIndex
	}
	
	truncated := length-(int(value.x1000)*3)
	valueS := value.amount.String()
	valueS = valueS[:truncated]
	value.dig.SetString(valueS, 10)
}

// Zet een geldwaarde met of zonder x1000-suffix om in een groot getal
func setValue(value *Number, input string) {
	value.originalInput = input // Bewaar originele invoer
	value.hasSuffix = false     // Standaard: geen suffix

	// Zoek eerst naar een suffix
	found := false
	var numberPart string
	for i := len(x1000) - 1; i > 0; i-- { // Loop door x1000 vanaf de grootste waarde
		suffix := x1000[i]
		if strings.HasSuffix(input, suffix) { // Controleer of input eindigt op de suffix
			numberPart = strings.TrimSuffix(input, suffix)
			// Verwijder alle spaties, punten en komma's
			numberPart = strings.ReplaceAll(numberPart, " ", "")
			numberPart = strings.ReplaceAll(numberPart, ".", "")
			numberPart = strings.ReplaceAll(numberPart, ",", "")
			// Controleer of numberPart een geldig getal is
			if _, ok := value.amount.SetString(numberPart, 10); !ok {
				fmt.Println("Ongeldige invoer: geen geldig getal")
				return
			}
			multiplier := new(big.Int).Exp(big.NewInt(1000), big.NewInt(int64(i)), nil)
			value.amount.Mul(&value.amount, multiplier)
			value.hasSuffix = true // Markeer dat invoer een suffix had
			found = true
			break
		}
	}

	if !found { // Geen geldige suffix gevonden, behandel het als een normaal getal
		// Verwijder alle spaties, punten en komma's
		cleanInput := strings.ReplaceAll(input, " ", "")
		cleanInput = strings.ReplaceAll(cleanInput, ".", "")
		cleanInput = strings.ReplaceAll(cleanInput, ",", "")
		// Controleer of cleanInput een geldig getal is
		if _, ok := value.amount.SetString(cleanInput, 10); !ok {
			fmt.Println("Ongeldige invoer: geen geldig getal")
			return
		}
	}

	setX1000(value)
}

func show(value Number) {
	fmt.Printf("\n%s\n%s %s%s\n", formatBigNumber(value.amount.String()), formatBigNumber(value.dig.String()), x1000text[value.x1000], x1000[value.x1000])
}

func main() {
	var value Number
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nValue: ")
		scanner.Scan()
		input := scanner.Text()
		setValue(&value, input)
		show(value)
	}
}
