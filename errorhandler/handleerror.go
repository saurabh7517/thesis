package errorhandler

/*
This function handles error
*/
func Check(e error) {
	if e != nil {
		panic(e)
	}
}
