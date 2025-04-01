package main

import (
	"fmt"
	"unicode"
)


func isPalindorme(sentence string) bool {

	left := 0 
	right := len(sentence)-1

	fmt.Println(sentence, "adlasf")
	for {
		if left >= right { 
			return true

		} 
		
		

		if unicode.IsLetter(rune(sentence[left])) && unicode.IsLetter(rune(sentence[right])) {
			// fmt.Println("a, b", a, b)
			if sentence[right] != sentence[left] {

				return false
				} 
				
			left += 1 
			right -= 1
		} else if !unicode.IsLetter(rune(sentence[left])) {
			left += 1 

		} else {
			right -= 1 

		}
	
	}
}