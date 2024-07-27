package helpers

import "runtime"

func UserAgent() string {
	return "ShuffleShowdown/0.1 (https://shuffle-showdown.joelbosch.nl) " + runtime.Version()
}
