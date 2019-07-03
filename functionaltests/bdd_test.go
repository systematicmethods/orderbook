package functionaltests

import (
	"flag"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/colors"
	"os"
	"testing"
)

var opt = godog.Options{Output: colors.Colored(os.Stdout)}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func Test_Features(t *testing.T) {
	flag.Parse()
	opt.Paths = flag.Args()

	stat := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContextLimitOrder(s)
	}, opt)

	if stat != 0 {
		t.Errorf("func tests failed")
	}
}
