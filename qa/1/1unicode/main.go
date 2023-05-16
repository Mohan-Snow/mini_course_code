package main

import "fmt"

/*
1
вопрос насчёт строк и рун: т.е. то что мультибайтная строка -
будет итерироваться в рейджне в размере инт32 - это магия самого range ? т.е.
почему обращение по индксу и рейндж - имеют разное поведение*/

/*
2
Как компилятор "угадывает" размер в байтах руны?
Кажется, что все символы строки станут по размеру самой большой руны и тогда угадывать нечего?
Где хранит размер рун компилятор, если структура StringHeader всего два поля имеет?
*/

// Кажется, что все символы строки станут по размеру самой большой руны и тогда угадывать нечего?
// abcde - 5 bytes, abcde - 20 bytes
func main() {
	s := "Hello, 世界, прив"

	for i := 0; i < len(s); i++ {
		fmt.Printf("[idx: %d, byte: %x, string: %s]\n ", i, s[i], string(s[i]))
	}
	fmt.Println()
	fmt.Println()
	fmt.Println()

	for i, r := range s {
		fmt.Printf("[%d: %x, %c]\n", i, r, r)
	}
	fmt.Println()

	// '0' - 1byte, it's ASCII
	// '110'   - 2 bytes
	// '1110'  - 3 bytes
	// '11110' - 4 bytes

	// "h" = \x68 = 0b01101000 - начинается с 0

	// "П" = \xd0\x9f
	// \xd0 = 0b11010000
	// \x9f = 0b10011111

	// "世"= \xe4\xb8\x96.
	// \xe4\xb8\x96
	// \xe4 = 11100100
	// \xb8 = 10111000
	// \x96 = 10010110

	// look at: https://www.branah.com/unicode-converter
	// \u0048\u0065\u006c\u006c\u006f\u002c \u4e16\u754c\u002c \u043f\u0440\u0438\u0432
	// no u: 00480065006c006c006f002c 4e16754c002c 043f044004380432

	s = "Hello, 世界, прив"
	for _, r := range s {
		fmt.Printf("%08b\n", r)
	}
}
