// SPDX-License-Identifier: LGPL-3.0-or-later

package main

import "testing"

func TestCalculatorBasicOperation(t *testing.T) {
	calculator := newCalculator()
	calculator.digit("1")
	calculator.digit("2")
	calculator.choose("+")
	calculator.digit("3")
	calculator.equals()
	if calculator.display != "15" || calculator.expression != "12 + 3 =" {
		t.Fatalf("calculation = %q, %q", calculator.display, calculator.expression)
	}
}

func TestCalculatorChainsOperations(t *testing.T) {
	calculator := newCalculator()
	calculator.digit("8")
	calculator.choose("×")
	calculator.digit("5")
	calculator.choose("−")
	calculator.digit("4")
	calculator.equals()
	if calculator.display != "36" {
		t.Fatalf("chained result = %q", calculator.display)
	}
}

func TestCalculatorDecimalSignPercentAndBackspace(t *testing.T) {
	calculator := newCalculator()
	calculator.digit("5")
	calculator.decimal()
	calculator.digit("2")
	calculator.digit("5")
	calculator.backspace()
	calculator.sign()
	calculator.percent()
	if calculator.display != "-0.052" {
		t.Fatalf("transformed number = %q", calculator.display)
	}
}

func TestCalculatorDivisionByZeroAndClear(t *testing.T) {
	calculator := newCalculator()
	calculator.digit("9")
	calculator.choose("÷")
	calculator.digit("0")
	calculator.equals()
	if !calculator.error || calculator.display != "Cannot divide by zero" {
		t.Fatalf("division error = %#v", calculator)
	}
	calculator.clear()
	if calculator.error || calculator.display != "0" || calculator.operator != "" {
		t.Fatalf("cleared calculator = %#v", calculator)
	}
}
