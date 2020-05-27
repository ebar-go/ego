package secure


// Panic only panic when err not nil
func Panic(err error)  {
	if err != nil {
		panic(err)
	}
}
