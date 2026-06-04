package main

func main() {
	app, err := wireApp()
	if err != nil {
		panic(err)
	}
	app.Run()
}
