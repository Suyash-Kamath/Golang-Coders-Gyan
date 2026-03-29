package main

import "fmt"

func ascii_unicode_utf8() {
	s := "A😀"

	fmt.Println(len(s))         // 5 — bytes (1 + 4)
	fmt.Println(len([]rune(s))) // 2 — actual characters

	// Byte level — raw bytes
	for i := 0; i < len(s); i++ {
		fmt.Println(s[i]) // raw byte values — emoji gets broken
	}

	// Rune level — correct characters
	for i, c := range s {
		fmt.Println(i, string(c)) // A and 😀 properly
	}

}
