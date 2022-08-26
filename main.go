package main

func main() {
	config, db := read_conf()

	run_bot(config, *db)
}
