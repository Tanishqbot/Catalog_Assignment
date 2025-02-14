package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"strconv" // Import the strconv package
)

type Point struct {
	Base  string `json:"base"`
	Value string `json:"value"`
}

type TestCase struct {
	Keys struct {
		N int `json:"n"`
		K int `json:"k"`
	} `json:"keys"`
	Points map[string]Point `json:"-"`
}

func main() {
	testCase1 := `{
		"keys": {
			"n": 4,
			"k": 3
		},
		"1": {
			"base": "10",
			"value": "4"
		},
		"2": {
			"base": "2",
			"value": "111"
		},
		"3": {
			"base": "10",
			"value": "12"
		},
		"6": {
			"base": "4",
			"value": "213"
		}
	}`

	testCase2 := `{
		"keys": {
			"n": 10,
			"k": 7
		},
		"1": {
			"base": "6",
			"value": "13444211440455345511"
		},
		"2": {
			"base": "15",
			"value": "aed7015a346d63"
		},
		"3": {
			"base": "15",
			"value": "6aeeb69631c227c"
		},
		"4": {
			"base": "16",
			"value": "e1b5e05623d881f"
		},
		"5": {
			"base": "8",
			"value": "316034514573652620673"
		},
		"6": {
			"base": "3",
			"value": "2122212201122002221120200210011020220200"
		},
		"7": {
			"base": "3",
			"value": "20120221122211000100210021102001201112121"
		},
		"8": {
			"base": "6",
			"value": "20220554335330240002224253"
		},
		"9": {
			"base": "12",
			"value": "45153788322a1255483"
		},
		"10": {
			"base": "7",
			"value": "1101613130313526312514143"
		}
	}`

	secret1 := processTestCase(testCase1)
	secret2 := processTestCase(testCase2)

	fmt.Println("Secret for test case 1:", secret1)
	fmt.Println("Secret for test case 2:", secret2)
}

func processTestCase(testCase string) *big.Int {
	var data map[string]json.RawMessage
	if err := json.Unmarshal([]byte(testCase), &data); err != nil {
		panic(err)
	}

	var tc TestCase
	if err := json.Unmarshal(data["keys"], &tc.Keys); err != nil {
		panic(err)
	}

	tc.Points = make(map[string]Point)
	for key, value := range data {
		if key == "keys" {
			continue
		}
		var p Point
		if err := json.Unmarshal(value, &p); err != nil {
			panic(err)
		}
		tc.Points[key] = p
	}

	var pointsList []struct {
		X *big.Int
		Y *big.Int
	}

	for keyStr, p := range tc.Points {
		x := new(big.Int)
		if _, ok := x.SetString(keyStr, 10); !ok {
			panic(fmt.Errorf("invalid x key: %s", keyStr))
		}

		base, err := strconv.Atoi(p.Base) // Use strconv.Atoi to convert base to int
		if err != nil {
			panic(fmt.Errorf("invalid base %s: %v", p.Base, err))
		}

		valueStr := strings.ToLower(p.Value)
		y := new(big.Int)
		if _, ok := y.SetString(valueStr, base); !ok {
			panic(fmt.Errorf("invalid value %s in base %d", p.Value, base))
		}

		pointsList = append(pointsList, struct {
			X *big.Int
			Y *big.Int
		}{X: x, Y: y})
	}

	sort.Slice(pointsList, func(i, j int) bool {
		return pointsList[i].X.Cmp(pointsList[j].X) < 0
	})

	k := tc.Keys.K
	if k > len(pointsList) {
		panic("insufficient points")
	}
	selected := pointsList[:k]

	prodNums := make([]*big.Int, k)
	prodDens := make([]*big.Int, k)

	for i := 0; i < k; i++ {
		xi := selected[i].X
		prodNum := big.NewInt(1)
		prodDen := big.NewInt(1)

		for j := 0; j < k; j++ {
			if j == i {
				continue
			}
			xj := selected[j].X

			termNum := new(big.Int).Neg(xj)
			prodNum.Mul(prodNum, termNum)

			termDen := new(big.Int).Sub(xi, xj)
			prodDen.Mul(prodDen, termDen)
		}

		prodNums[i] = prodNum
		prodDens[i] = prodDen
	}

	productDenominators := big.NewInt(1)
	for _, den := range prodDens {
		productDenominators.Mul(productDenominators, den)
	}

	sumNum := big.NewInt(0)
	for i := 0; i < k; i++ {
		y := selected[i].Y
		num := prodNums[i]
		den := prodDens[i]

		div := new(big.Int).Div(productDenominators, den)
		term := new(big.Int).Mul(y, num)
		term.Mul(term, div)

		sumNum.Add(sumNum, term)
	}

	c := new(big.Int).Div(sumNum, productDenominators)
	return c
}