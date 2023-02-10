package transporter

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"wget/pkg/utils"
)

type Credentials struct {
	URL          string
	FileName     string
	Path         string
	RateLimit    int64
	IsBackground bool
	OutPut       *os.File
}

func NewCredentialsConstructor(URL string) *Credentials {
	return &Credentials{
		URL:          URL,
		FileName:     getFileName(URL),
		Path:         "./",
		IsBackground: false,
		OutPut:       os.Stdout,
	}
}

func getFileName(URL string) string {
	arr := strings.Split(URL, "/")
	return arr[len(arr)-1]
}

var (
	backgroundFlagVal bool
	nameFlagVal       string
	pathFlagVal       string
	rateLimit         string
)

func Switcher() error {
	flag.BoolVar(&backgroundFlagVal, "B", false, "background download")
	flag.StringVar(&nameFlagVal, "O", "tempfile", "give name to saved file")
	flag.StringVar(&pathFlagVal, "P", "./", "path to where you want to save the file")
	flag.StringVar(&rateLimit, "rate-limit", "max", "handle speed limit")
	flag.Parse()

	credentials := NewCredentialsConstructor(flag.Arg(0))

	if err := flagsChecker(credentials); err != nil {
		return err
	}

	if err := Download(credentials); err != nil {
		return err
	}

	return nil
}

func flagsChecker(credentials *Credentials) error {
	if nameFlagVal != "tempfile" {
		credentials.FileName = nameFlagVal
	}

	if pathFlagVal != "./" {
		if pathFlagVal[0] == '~' {
			homePath, err := os.UserHomeDir()
			if err != nil {
				return err
			}

			credentials.Path = homePath + pathFlagVal[1:] + "/"

		} else {
			credentials.Path = pathFlagVal + "/"
		}

		err := os.Mkdir(credentials.Path, 0o700)
		if err != nil && !errors.Is(err, os.ErrExist) {
			return err
		}

	}

	if backgroundFlagVal {
		logFile, err := os.Create(utils.DefaultLog)
		if err != nil {
			return err
		}
		fmt.Println(`Output will be written to "wget-log"`)
		credentials.OutPut = logFile

		parentProc := syscall.Getppid()
		syscall.Kill(parentProc, syscall.SIGTSTP)
		defer syscall.Kill(parentProc, syscall.SIGCONT)
	}

	if rateLimit != "max" {
		speed, err := setSpeed(rateLimit)
		if err != nil {
			return err
		}
		credentials.RateLimit = speed

	}

	return nil
}

func setSpeed(rateLimit string) (int64, error) {
	var speed int64
	var isValid bool
	err := errors.New("wrong type of speed")

	for _, ch := range rateLimit {
		if isValid {
			return 0, err
		}

		if ch >= '0' && ch <= '9' {
			speed *= 10
			speed += int64(ch - 48)
			continue
		}

		switch ch {
		case 'k':
			speed *= 1000
			isValid = true
		case 'M':
			speed *= 1000000
			isValid = true
		default:
			return 0, err
		}
	}

	return speed, nil
}
