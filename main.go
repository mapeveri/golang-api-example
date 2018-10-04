package main

func main() {
	a := App{}
	a.Initialize(
		"postgres",
		"Polgara8",
		"fejipba_development")

	a.Run(":8080")
}
