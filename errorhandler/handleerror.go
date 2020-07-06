package errorhandler

import "fmt"

/*
This function handles error
*/
func Check(e error) {
	if e != nil {
		fmt.Println(e)
		return
	}
}
