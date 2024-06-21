package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ПРЕОБРАЗОВАНИЕ РИМСКИХ ЧИСЕЛ В АРАБСКИЕ
var RomArab = map[string]int{
	"I":    1,
	"II":   2,
	"III":  3,
	"IV":   4,
	"V":    5,
	"VI":   6,
	"VII":  7,
	"VIII": 8,
	"IX":   9,
	"X":    10,
}

// Преобразование арабских чисел в римские по индексу
var ArabRom = []string{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X"}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Ошибка:", r)
			fmt.Println("Программа завершена из-за некорректного ввода.")
		}
	}()

	for {
		fmt.Println("Введите выражение (или 'ex' для выхода):")
		var input string
		fmt.Scanln(&input)
		if input == "ex" {
			break
		}

		result, err := calculate(input)
		if err != nil {
			fmt.Println("Ошибка:", err)
		} else {
			fmt.Println("Результат:", result)
		}
	}
}

func calculate(input string) (string, error) {
	// Удаляем пробелы из ввода
	input = strings.ReplaceAll(input, " ", "")

	// Проверяем минимальную длину ввода
	if len(input) < 3 {
		panic("некорректный формат ввода")
	}

	// Инициализация переменных
	var aStr string
	var id int

	// Извлекаем первое число
	for id = 0; id < len(input); id++ {
		if input[id] == '+' || input[id] == '-' || input[id] == '*' || input[id] == '/' {
			break
		}
		aStr += string(input[id])
	}

	// Проверка на корректность
	if len(aStr) == 0 {
		panic("некорректный формат ввода")
	}

	// Извлекаем оператор
	if id == len(input) {
		panic("не найден оператор")
	}
	operator := string(input[id])

	// Извлекаем второе число
	var bStr string
	for id++; id < len(input); id++ {
		bStr += string(input[id])
	}

	if len(bStr) == 0 {
		panic("некорректный формат ввода")
	}

	// Проверяем, является ли первое число римским
	isRomanA := isRomanNumeral(aStr)
	isRomanB := isRomanNumeral(bStr)

	var a, b int
	var errA, errB error

	if isRomanA && isRomanB {
		// Если оба числа римские, конвертируем их в арабские
		a, errA = romanToArabic(aStr)
		b, errB = romanToArabic(bStr)
	} else if !isRomanA && !isRomanB {
		// Если оба числа арабские, преобразуем строки в числа
		a, errA = strconv.Atoi(aStr)
		b, errB = strconv.Atoi(bStr)
	} else {
		panic("невозможно выполнить операцию между римским и арабским числами")
	}

	if errA != nil || errB != nil {
		panic("некорректный ввод чисел или оператора")
	}

	// Проверяем, что числа находятся в диапазоне от 1 до 10 включительно
	if (a < 1 || a > 10) || (b < 1 || b > 10) {
		panic("числа должны быть в диапазоне от 1 до 10 включительно")
	}

	// Выполняем операцию
	result, err := performOperation(a, b, operator)
	if err != nil {
		panic(err.Error())
	}

	// Если оба числа были римскими, конвертируем результат в римские числа
	if isRomanA && isRomanB {
		return arabicToRoman(result)
	}

	return strconv.Itoa(result), nil
}

func performOperation(a, b int, operator string) (int, error) {
	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, errors.New("деление на ноль")
		}
		return a / b, nil
	default:
		return 0, errors.New("некорректный оператор")
	}
}

// Функция для проверки, является ли строка римским числом
func isRomanNumeral(s string) bool {
	_, ok := RomArab[s]
	return ok
}

// Функция для конвертации римских чисел в арабские
func romanToArabic(roman string) (int, error) {
	roman = strings.ToUpper(roman)
	romanMap := map[rune]int{
		'I': 1,
		'V': 5,
		'X': 10,
	}

	arabic := 0
	length := len(roman)

	for i := 0; i < length; i++ {
		value, exists := romanMap[rune(roman[i])]
		if !exists {
			return 0, fmt.Errorf("неверный римский символ: %c", roman[i])
		}

		// Если это последний символ или текущий символ больше или равен следующему, то прибавляем значение
		if i+1 < length && value < romanMap[rune(roman[i+1])] {
			arabic -= value
		} else {
			arabic += value
		}
	}

	return arabic, nil
}

// Функция для конвертации арабских чисел в римские
func arabicToRoman(arabic int) (string, error) {
	if arabic < 1 || arabic > 3999 {
		return "", fmt.Errorf("невозможно представить число %d в римской системе", arabic)
	}

	var result strings.Builder

	// Массив пар чисел и их римских эквивалентов, упорядоченный по убыванию
	numeralMap := []struct {
		Value  int
		Symbol string
	}{ //Отношения ри
		{100, "C"}, {90, "XC"}, {50, "L"}, {40, "XL"},
		{10, "X"}, {9, "IX"}, {5, "V"}, {4, "IV"}, {1, "I"},
	}

	for _, pair := range numeralMap {
		for arabic >= pair.Value {
			result.WriteString(pair.Symbol)
			arabic -= pair.Value
		}
	}

	return result.String(), nil
}

