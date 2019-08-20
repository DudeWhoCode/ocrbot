package main

func main() {
	f, err := newMultiLogger()
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to create log file \n %s", err)
	}

	startServer()
}
