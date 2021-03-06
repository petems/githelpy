package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/petems/githelpy/githelpy"
)

func TestCheckRangeCommitWithErrors(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()

	if err != nil {
		logrus.Fatal(err)
	}

	path, err := os.Getwd()

	if err != nil {
		logrus.Fatal(err)
	}

	exitError = func() {
		panic(1)
	}

	exitSuccess = func() {
		panic(0)
	}

	var errc error

	failure = func(err error) {
		errc = err
	}

	arguments := [][]string{
		[]string{
			"check",
			"range",
		},
		[]string{
			"check",
			"range",
			"master~2",
		},
		[]string{
			"check",
			"range",
			"master~1",
			"master~2",
			"test",
			"whatever",
		},
		[]string{
			"check",
			"range",
			"master~1",
			"master~2",
			"whatever",
		},
		[]string{
			"check",
			"range",
			"master~1",
			"master~2",
			"check.go",
		},
		[]string{
			"check",
			"range",
			"whatever",
			"master",
			"test/",
		},
		[]string{
			"check",
			"range",
			"master~2",
			"master~1",
			"test/",
		},
		[]string{
			"check",
			"range",
			"master~2",
			"master~1",
			"test/",
		},
	}

	errors := []error{
		fmt.Errorf("Two arguments required : origin commit and end commit"),
		fmt.Errorf("Two arguments required : origin commit and end commit"),
		fmt.Errorf("3 arguments must be provided at most"),
		fmt.Errorf(`Ensure "whatever" directory exists`),
		fmt.Errorf(`"check.go" must be a directory`),
		fmt.Errorf(`Can't find reference "whatever"`),
		fmt.Errorf(`At least one matcher must be defined`),
		fmt.Errorf(`At least one example must be defined`),
	}

	configs := []string{
		path + "/../features/.githelpy.toml",
		path + "/../features/.githelpy.toml",
		path + "/../features/.githelpy.toml",
		path + "/../features/.githelpy.toml",
		path + "/../features/.githelpy.toml",
		path + "/../features/.githelpy.toml",
		path + "/../features/.githelpy-no-matchers.toml",
		path + "/../features/.githelpy-no-examples.toml",
	}

	for i, a := range arguments {
		var w sync.WaitGroup

		w.Add(1)

		go func() {
			defer func() {
				if r := recover(); r != nil && r.(int) == 0 {
					errc = nil
				}

				w.Done()
			}()

			os.Args = []string{"", "--config", configs[i]}
			os.Args = append(os.Args, a...)

			_ = RootCmd.Execute()
		}()

		w.Wait()

		assert.Error(t, errc, "Must return an error")
		assert.EqualError(t, errc, errors[i].Error(), "Must return an error : "+errors[i].Error())
	}
}

func TestCheckRangeCommitWithBadCommitMessage(t *testing.T) {
	path, err := os.Getwd()

	if err != nil {
		logrus.Fatal(err)
	}

	for _, filename := range []string{"../features/repo.sh", "../features/bad-commit.sh"} {

		err := exec.Command(filename).Run()

		if err != nil {
			logrus.Fatal(err)
		}
	}

	exitError = func() {
		panic(1)
	}

	exitSuccess = func() {
		panic(0)
	}

	var w sync.WaitGroup

	var matchings *[]*githelpy.Matching
	var examples map[string]string

	renderMatchings = func(m *[]*githelpy.Matching) {
		matchings = m
	}

	renderExamples = func(e map[string]string) {
		examples = e
	}

	w.Add(1)

	go func() {
		defer func() {
			if r := recover(); r != nil && r.(int) == 0 {
				matchings = &[]*githelpy.Matching{}
				examples = map[string]string{}
			}

			w.Done()
		}()

		os.Args = []string{"", "--config", path + "/../features/.githelpy.toml", "check", "range", "master~3", "master", path + "/test"}

		Execute()
	}()

	w.Wait()

	assert.Len(t, *matchings, 1, "Must return 1 commits")
	assert.Len(t, examples, 3, "Must return 3 examples")
}

func TestCheckRangeCommitWithNoErrors(t *testing.T) {
	path, err := os.Getwd()

	if err != nil {
		logrus.Fatal(err)
	}

	for _, filename := range []string{"../features/repo.sh"} {

		err := exec.Command(filename).Run()

		if err != nil {
			logrus.Fatal(err)
		}
	}

	var code int
	var message string
	var w sync.WaitGroup

	success = func(msg string) {
		message = msg
	}

	exitError = func() {
		panic(1)
	}

	exitSuccess = func() {
		panic(0)
	}

	w.Add(1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				code = r.(int)
			}

			w.Done()
		}()

		os.Args = []string{"", "--config", path + "/../features/.githelpy.toml", "check", "range", "master~2", "master", path + "/test"}

		Execute()
	}()

	w.Wait()

	assert.EqualValues(t, 0, code, "Must exit without errors (exit 0)")
	assert.EqualValues(t, "Everything is ok", message, "Must return a message to inform everything is ok")
}
