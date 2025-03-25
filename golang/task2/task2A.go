
/*
var tests int 
var sentences []string

reader := bufio.NewReader(os.Stdin)

fmt.Scan(&tests)

for  i := 0; i < tests; i++ {
	// var sentence string
	// fmt.Scanln(&sentence)

	input, _ := reader.ReadString('\n')

	sentences = append(sentences, input)

}


for _, sentence := range sentences {

	frequencyCounter(sentence)
}

}

 */

package main

import (
	"fmt"
)

func frequencyCounter(sentence string) {

	count  := make(map[string]int)
	words := extractWords(sentence)

	for _, word := range words {
		count[word] += 1 
	}

	fmt.Println("words ", words, count)


}


func extractWords(sentece string) []string {

	length := len(sentece)
	var words []string

	left := 0 
	for {
        
		var cur_word string
		if left >= length {
			break 
		}

		right := left 
		for right < length {
			
			char := sentece[right]
			if char == ' ' {
				right +=1 
				break 
			}

			ascival := int(char)
			if  (65 <= ascival && ascival <= 90) || (97 <= ascival && ascival <= 122) {
				cur_word += string(sentece[right])
			}


			right +=  1 
	
		}

		if cur_word != "" {
			words = append(words, cur_word)
		}

		left = right

	}

	return words

}
