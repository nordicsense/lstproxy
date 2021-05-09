package fail

import "log"

func OnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
