package main

func main() {
	LoggerInit()

	wg.Add(1)

	go DeleteNotices()
	RunBot()
}
