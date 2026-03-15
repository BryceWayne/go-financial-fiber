// main.go
package main

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

// FinanceData holds all inputs and calculated outputs
type FinanceData struct {
	// --- INPUTS ---
	// Income
	Salary         float64 `form:"salary"`
	Bonus          float64 `form:"bonus"`
	RentalProperty float64 `form:"rental_property"`
	InterestIncome float64 `form:"interest_income"`

	// Location
	State string `form:"state"`
	City  string `form:"city"`

	// Deductions & Rates (Hardcoded in calculation for now, but could be inputs)
	Monthly401kContribution float64 `form:"monthly_401k"`

	// Responsible Expenses
	Mortgage        float64 `form:"mortgage"`
	Groceries       float64 `form:"groceries"`
	Utilities       float64 `form:"utilities"`
	StudentLoans    float64 `form:"student_loans"`
	CreditCardDebt  float64 `form:"credit_card_debt"`
	Investments     float64 `form:"investments"`
	Laundry         float64 `form:"laundry"`
	Car             float64 `form:"car"`
	PublicTransport float64 `form:"public_transport"`

	// Bonus Spending
	EatingOut float64 `form:"eating_out"`
	Travel    float64 `form:"travel"`
	Gym       float64 `form:"gym"`
	Hair      float64 `form:"hair"`
	Skincare  float64 `form:"skincare"`
	Clothes   float64 `form:"clothes"`
	Coffee    float64 `form:"coffee"`
	OtherSubs float64 `form:"other_subs"`

	// Irresponsible Stuff
	FunRandom float64 `form:"fun_random"`
	Alcohol   float64 `form:"alcohol"`
	Ubers     float64 `form:"ubers"`

	// Savings
	Savings401k      float64 `form:"savings_401k"`
	Cash             float64 `form:"cash"`
	OtherInvestments float64 `form:"other_investments"`
	Brokerage        float64 `form:"brokerage"`
	HomeEquity       float64 `form:"home_equity"`

	// --- CALCULATED OUTPUTS ---
	TotalPreTaxAnnual float64

	PreTaxMonthly   float64
	MonthlyTaxable  float64
	MonthlyFedTax   float64
	MonthlyStateTax float64
	MonthlyCityTax  float64
	PostTaxTakeHome float64

	BonusTaxable         float64
	BonusTaxes           float64
	TotalBonusTakeHome   float64
	TotalAnnualTakeHome  float64
	TotalMonthlyTakeHome float64

	TotalResponsible   float64
	TotalBonusSpending float64
	TotalIrresponsible float64
	TotalSavings       float64
}

// formatCurrency adds commas and a dollar sign (e.g., 1050000 -> $1,050,000)
func formatCurrency(n float64) string {
	in := strconv.FormatFloat(n, 'f', 0, 64)
	var out []byte
	for i, c := range in {
		if i > 0 && (len(in)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, byte(c))
	}
	return "$" + string(out)
}

// formatPercent converts a float to a percentage string
func formatPercent(n float64) string {
	return strconv.FormatFloat(n*100, 'f', 0, 64) + "%"
}

func main() {
	engine := html.New("./views", ".html")
	engine.AddFunc("formatCurrency", formatCurrency)
	engine.AddFunc("formatPercent", formatPercent)

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Initial page load
	app.Get("/", func(c *fiber.Ctx) error {
		// Provide default values matching the screenshot
		defaultData := FinanceData{
			Salary: 350000, Bonus: 700000, RentalProperty: 12000, InterestIncome: 3000,
			State: "NY", City: "NYC",
			Monthly401kContribution: 1950,
			Mortgage:                6500, Groceries: 900, Utilities: 375, Investments: 32000, Laundry: 50,
		}
		calculateTotals(&defaultData)
		return c.Render("index", defaultData)
	})

	// HTMX endpoint
	app.Post("/calculate", func(c *fiber.Ctx) error {
		data := new(FinanceData)
		if err := c.BodyParser(data); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		calculateTotals(data)
		// We re-render ONLY the dashboard partial
		return c.Render("dashboard", data)
	})

	log.Println("Server starting on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}

type TaxBracket struct {
	Limit float64
	Rate  float64
}

// calcTax computes exact marginal tax
func calcTax(income float64, brackets []TaxBracket) float64 {
	tax := 0.0
	prev := 0.0
	for _, b := range brackets {
		if income > prev {
			taxable := income - prev
			if b.Limit > 0 && income > b.Limit {
				taxable = b.Limit - prev
			}
			tax += taxable * b.Rate
			prev = b.Limit
		} else {
			break
		}
	}
	return tax
}

// 2024 Single Filer Brackets
var fedBraces = []TaxBracket{{11600, 0.10}, {47150, 0.12}, {100525, 0.22}, {191950, 0.24}, {243725, 0.32}, {609350, 0.35}, {0, 0.37}}
var nyBraces = []TaxBracket{{8500, 0.04}, {11700, 0.045}, {13900, 0.0525}, {80650, 0.0585}, {215400, 0.0597}, {1077550, 0.0685}, {5000000, 0.0965}, {0, 0.103}}
var caBraces = []TaxBracket{{10412, 0.01}, {24684, 0.02}, {38959, 0.04}, {54081, 0.06}, {68350, 0.08}, {349137, 0.093}, {418961, 0.103}, {698271, 0.113}, {0, 0.123}}
var nycBraces = []TaxBracket{{12000, 0.03078}, {25000, 0.03762}, {50000, 0.03819}, {0, 0.03876}}

func calculateTotals(d *FinanceData) {
	// 1. Income Math
	d.TotalPreTaxAnnual = d.Salary + d.Bonus + d.RentalProperty + d.InterestIncome
	annualTaxable := d.TotalPreTaxAnnual - (d.Monthly401kContribution * 12)

	// Calculate True Effective Rates using Marginal Brackets
	annualFedTax := calcTax(annualTaxable, fedBraces)
	
	annualStateTax := 0.0
	switch d.State {
	case "NY": annualStateTax = calcTax(annualTaxable, nyBraces)
	case "CA": annualStateTax = calcTax(annualTaxable, caBraces)
	case "IL": annualStateTax = annualTaxable * 0.0495 // flat
	case "Other": annualStateTax = annualTaxable * 0.05 // proxy
	}

	annualCityTax := 0.0
	switch d.City {
	case "NYC": annualCityTax = calcTax(annualTaxable, nycBraces)
	case "SF": annualCityTax = annualTaxable * 0.015 // generic local tax proxy
	case "Other": annualCityTax = annualTaxable * 0.01
	}

	effFed := 0.0
	effState := 0.0
	effCity := 0.0
	if annualTaxable > 0 {
		effFed = annualFedTax / annualTaxable
		effState = annualStateTax / annualTaxable
		effCity = annualCityTax / annualTaxable
	}

	// 2. Post-Tax Monthly Math (Excluding Bonus)
	d.PreTaxMonthly = (d.Salary + d.RentalProperty + d.InterestIncome) / 12
	d.MonthlyTaxable = d.PreTaxMonthly - d.Monthly401kContribution

	// Apply effective tax rates to regular paycheck
	d.MonthlyFedTax = d.MonthlyTaxable * effFed
	d.MonthlyStateTax = d.MonthlyTaxable * effState
	d.MonthlyCityTax = d.MonthlyTaxable * effCity

	d.PostTaxTakeHome = d.PreTaxMonthly - d.Monthly401kContribution - d.MonthlyFedTax - d.MonthlyStateTax - d.MonthlyCityTax

	// 3. Bonus Math
	d.BonusTaxable = d.Bonus
	// Use the effective rates for the bonus so that the total tax matches the annual bracket math exactly.
	d.BonusTaxes = d.BonusTaxable * (effFed + effState + effCity)
	d.TotalBonusTakeHome = d.BonusTaxable - d.BonusTaxes

	// 4. Grand Totals
	d.TotalAnnualTakeHome = (d.PostTaxTakeHome * 12) + d.TotalBonusTakeHome
	d.TotalMonthlyTakeHome = d.TotalAnnualTakeHome / 12

	// 5. Expense & Savings Totals
	d.TotalResponsible = d.Mortgage + d.Groceries + d.Utilities + d.StudentLoans + d.CreditCardDebt + d.Investments + d.Laundry + d.Car + d.PublicTransport
	d.TotalBonusSpending = d.EatingOut + d.Travel + d.Gym + d.Hair + d.Skincare + d.Clothes + d.Coffee + d.OtherSubs
	d.TotalIrresponsible = d.FunRandom + d.Alcohol + d.Ubers
	d.TotalSavings = d.Savings401k + d.Cash + d.OtherInvestments + d.Brokerage + d.HomeEquity
}
