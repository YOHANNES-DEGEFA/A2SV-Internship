
package main

import (
	"fmt"
)

type data  struct {
		name string 
		score int
}




func gradeReport() {

    var student string 
	var subjects int 

	var info1 data
	var arr []data

	fmt.Scanln(&student)
	fmt.Scanln(&subjects)
    

	info1.name = student
	info1.score = subjects
	
	arr = append(arr, info1)

	// fmt.Print(arr, "lets see the behavior")
	for i := 0; i < subjects; i ++ {

		var info data 
		
		fmt.Scanln(&info.name)
		fmt.Scanln(&info.score)

		arr = append(arr, info)


	} 

	displayGradeReport(arr)
	




}




func displayGradeReport(datas []data) {
	  average := findAverage(datas)

	  fmt.Println("Student Name:", datas[0].name) 
	  
	  fmt.Println("Number of subjects:", datas[0].score)

	  fmt.Print("\n\n\n") 

	  fmt.Println("\tName\t\tscore")
	  for _, data := range datas[1:] {
		fmt.Println("\t", data.name, "\t", data.score)
	  }


	  fmt.Print("\n\n\n")
	  fmt.Println("Total Subjects:", datas[0].score)
	
	  fmt.Println("Total Score:    ",  float64(len(datas)-1)*average)
	  fmt.Println("Average Score:", average)
        

}


func findAverage( arr []data) float64{
	total := 0 
	for _, info := range arr[1:] {
	
		 
		total += info.score
	}

	return float64(total)/float64(len(arr)-1)
}

func pp(item any){

	fmt.Println(item)
}