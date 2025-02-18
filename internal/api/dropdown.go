package api

import (
	"net/http"
)

type Dropdown struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var types []Dropdown
var expenses []Dropdown
var incomes []Dropdown

func init() {
	expenses = append(expenses, Dropdown{Key: "1", Value: "household"})
	expenses = append(expenses, Dropdown{Key: "2", Value: "food"})
	expenses = append(expenses, Dropdown{Key: "3", Value: "health"})
	expenses = append(expenses, Dropdown{Key: "4", Value: "beauty"})
	expenses = append(expenses, Dropdown{Key: "5", Value: "transport"})
	expenses = append(expenses, Dropdown{Key: "6", Value: "fashion"})
	expenses = append(expenses, Dropdown{Key: "7", Value: "social"})
	expenses = append(expenses, Dropdown{Key: "8", Value: "education"})
	expenses = append(expenses, Dropdown{Key: "9", Value: "other"})

	incomes = append(incomes, Dropdown{Key: "1", Value: "salary"})
	incomes = append(incomes, Dropdown{Key: "2", Value: "petty cash"})
	incomes = append(incomes, Dropdown{Key: "3", Value: "allowance"})
	incomes = append(incomes, Dropdown{Key: "4", Value: "other"})

	types = append(types, Dropdown{Key: "1", Value: "income"})
	types = append(types, Dropdown{Key: "2", Value: "expense"})
}

func GetExpenses() ([]Dropdown, int) {
	return expenses, http.StatusOK
}

func ValidExpense(input string) bool {
	for _, value := range expenses {
		if value.Value == input {
			return true
		}
	}
	return false;
}

func GetIncomes() ([]Dropdown, int) {
	return incomes, http.StatusOK
} 

func ValidIncome(input string) bool {
	for _, value := range incomes {
		if value.Value == input {
			return true
		}
	}
	return false;
}

func GetTypes() ([]Dropdown, int) {
	return types, http.StatusOK
}

func ValidType(input string) bool {
	for _, value := range types {
		if value.Value == input {
			return true
		}
	}
	return false;
}