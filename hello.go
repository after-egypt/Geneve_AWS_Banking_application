package main
import ("fmt"
	//""time")
)
func main(){
	b := []int{1,2,3,4}
	fmt.Print(b)
	for a  := range b {
	fmt.Print(a)
	}
}
