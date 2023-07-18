package cpick_test

import (
	"fmt"
	"testing"

	"github.com/ethanbaker/cpick"
)

func Test_Start(t *testing.T) {
	_, err := cpick.Start(true)
	//_, err := cpick.Start(false, false) // For developing

	if err != nil {
		fmt.Printf("cpick is not working in testing mode. If you think this is a bug, please report an issue on the gitlab page.\nError: %v", err)
	} else {
		fmt.Printf("# cpick is working in testing mode!\n")
	}
}
