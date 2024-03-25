package cmd

import (
	"fmt"
	"testing"
)

func TestAnalysisKeyFile(t *testing.T) {
	crypt := AnalysisKeyFile("../app.key")
	fmt.Printf("%+v\n", crypt)
}
