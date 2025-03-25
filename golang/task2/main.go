
package main 

import (
	"fmt"
	"bufio"
	"os"
)

// type data  struct {
// 		name string 
// 		score int
// }



func main() {

	// getting started
			// fmt.Println("hi there")
			// fmt.Println("Hi" + "there")
			// fmt.Print("Hi " + "there")
			// fmt.Println("1 + 1 =", 1+1)

	
	// data types in go 

			// var word string = "I Love Her!" 
			// var num int = 7
			// var decimal float64 = 5.8

			// var isEqual bool = true 

			// fmt.Println("word: ", word)
			// fmt.Println("num: ", num)
			// fmt.Println("decimal: ", decimal)
			// fmt.Println("bool:", isEqual)
	
	// variables in go // :=  the short-hand 
	
			// var a int 
			// var b = 3.3
			// p(a)
			// p(b)
			

    // const 
			// const A = 10
			// // A = 4 // cannot assign to A (neither addressable nor a map index expression
			// p(A)

	// for - loops    // the only loop type we have in go 
 

	// Ways to declear or use for loops 
		// way 1: 
				// i := 0 
				// for i < 3 {
				// 	p(i)
				// 	i += 1 
				// }

		// way 2:  start at 0 and endes at 2 
				// for i :=  range 3 {
				// 	p(i)
				// }

		// way 3: 
				
				// for i := 0; i < 5; i ++ {
				// 	p(i)
				// }

		// way 4: for loop as a while loop 
				// count :=  0 
				// endLoop := false 

		        // for {
				// 	if count == 100{
				// 	    endLoop = true } 

				// 	if endLoop {
				// 		break
				// 	}
				// 	print("running for ", count, "Times.")
				// 	count += 1 
				// }

	// CONDITIONALS  
	// arr := [...]int{ 1, 2, 3, 100: 4, 5, 9}
	// p(arr)


	// var a [5]int
	//    fmt.Println("emp:", a)
	// // We can set a value at an index using the array[index] = value syntax, and get a value with array[index].
	
	//    a[4] = 100
	//    fmt.Println("set:", a)
	//    fmt.Println("get:", a[4])
	// // The builtin len returns /the length of an array.
	
	//    fmt.Println("len:", len(a))
	// // Use this syntax to declare and initialize an array in one line.
	
	//    b := [5]int{1, 2, 3, 4, 5}
	//    fmt.Println("dcl:", b)
	// // You can also have the compiler count the number of elements for you with ...
	
	//    b = [...]int{1, 2, 3, 4, 5}
	//    fmt.Println("dcl:", b)
	// // If you specify the index with :, the elements in between will be zeroed.
	
	//    b = [...]int{100, 3: 400, 500}
	//    fmt.Println("idx:", b)
	// // Array types are one-dimensional, but you can compose types to build multi-dimensional data structures.
	
	//    var twoD [2][3]int
	//    for i := 0; i < 2; i++ {
	// 	   for j := 0; j < 3; j++ {
	// 		   twoD[i][j] = i + j
	// 	   }
	//    }
	//    fmt.Println("2d: ", twoD)
	// // You can create and initialize multi-dimensional arrays at once too.
	
	//    twoD = [2][3]int{
	// 	   {1, 2, 3},
	// 	   {1, 2, 3},
	//    }
	//    fmt.Println("2d: ", twoD)


    
	

	// gradeReport()

    //  gradeReport()

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

		if isPalindorme(sentence) {
				fmt.Println(sentence, "is palindrome")
		} else {
			fmt.Println(sentence, "is not palindrome")
		}
	}











}
func p(x any) {

	fmt.Println(x)
}


// func gradeReport() {

//     var student string 
// 	var subjects int 


// 	fmt.Scanln(&student)
// 	fmt.Scanln(&subjects)
    

// 	var arr []data
// 	for i :=0; i <subjects; i ++ {

// 		var info data 
		
// 		fmt.Scanln(&info.name)
// 		fmt.Scanln(&info.score)

// 		arr = append(arr, info)


// 	} 
	

// 	p(arr)
// }


