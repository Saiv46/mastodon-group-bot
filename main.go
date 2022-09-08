package main

func main() {
	wg.Add(1)

	go DeleteNotices()
	RunBot()
}
