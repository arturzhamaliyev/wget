package transporter

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"wget/pkg/utils"
)

type Credentials struct {
	URL          string
	FileName     string
	Path         string
	RateLimit    int64
	IsBackground bool
	IsInDir      bool
	OutPut       *os.File
	Mutex        *sync.Mutex
}

func NewCredentialsConstructor(URL string) *Credentials {
	return &Credentials{
		URL:          URL,
		FileName:     getFileName(URL),
		Path:         "./",
		RateLimit:    0,
		IsBackground: false,
		IsInDir:      false,
		OutPut:       os.Stdout,
		Mutex:        &sync.Mutex{},
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
	rateLimitFlagVal  string
	directoryFlagVal  string
	logOutputFlagVal  string
	mirrorFlagVal     bool
)

func Switcher() error {
	flag.BoolVar(&backgroundFlagVal, "B", false, "background download")
	flag.StringVar(&nameFlagVal, "O", "tempfile", "give name to saved file")
	flag.StringVar(&pathFlagVal, "P", "./", "path to where you want to save the file")
	flag.StringVar(&rateLimitFlagVal, "rate-limit", "max", "handle limit limit")
	flag.StringVar(&directoryFlagVal, "i", "", "download from file that will contain all links")
	flag.StringVar(&logOutputFlagVal, "logoutput", "os.Stdout", "default log output")
	flag.BoolVar(&mirrorFlagVal, "mirror", false, "download the entire website")
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

// func bfs(credentials *Credentials) {
// 	q := []*Credentials{credentials}

// 	for len(q) > 0 {
// 		curCred := q[0]

// 		curCred.

// 		q=q[1:]
// 	}
// }

func flagsChecker(credentials *Credentials) error {
	// if mirrorFlagVal {
	// 	if err := os.Mkdir(getFileName(credentials.URL), 0o700); err != nil {
	// 		return err
	// 	}

	// 	// TODO: call download recursively
	// 	// bfs(credentials)
	// }

	if logOutputFlagVal != "os.Stdout" {
		file, err := os.Create(logOutputFlagVal)
		if err != nil {
			return err
		}

		credentials.OutPut = file
	}

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
		defer logFile.Close()

		fmt.Println(`Output will be written to "wget-log"`)
		credentials.OutPut = logFile

		var command string
		argLen := len(os.Args)
		command += os.Args[0] + " --logoutput=wget-log "

		for i := 1; i < argLen; i++ {
			if os.Args[i] != "-B" {
				command += os.Args[i] + " "
			}
		}
		command += "&"

		cmd := exec.Command("/bin/bash", "-c", command)
		if err := cmd.Run(); err != nil {
			return err
		}

		os.Exit(0)
	}

	if rateLimitFlagVal != "max" {
		limit, err := setLimit(rateLimitFlagVal)
		if err != nil {
			return err
		}
		credentials.RateLimit = limit

	}

	if directoryFlagVal != "" && credentials.URL == "" {
		file, err := os.Open(directoryFlagVal)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		var URLs []string

		for scanner.Scan() {
			URLs = append(URLs, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		var wg sync.WaitGroup
		var mutex sync.Mutex

		for _, URL := range URLs {
			go func(URL string) {
				Download(&Credentials{
					URL:          URL,
					FileName:     getFileName(URL),
					Path:         credentials.Path,
					RateLimit:    credentials.RateLimit,
					IsBackground: credentials.IsBackground,
					IsInDir:      true,
					OutPut:       credentials.OutPut,
					Mutex:        &mutex,
				})
				wg.Done()
			}(URL)

			wg.Add(1)
		}

		wg.Wait()

		os.Exit(0)
	}

	return nil
}

func setLimit(rateLimitFlagVal string) (int64, error) {
	var limit int64
	var isValid bool
	err := errors.New("wrong type of limit")

	for _, ch := range rateLimitFlagVal {
		if isValid {
			return 0, err
		}

		if ch >= '0' && ch <= '9' {
			limit *= 10
			limit += int64(ch - 48)
			continue
		}

		switch ch {
		case 'k':
			limit *= 1000
			isValid = true
		case 'M':
			limit *= 1000000
			isValid = true
		default:
			return 0, err
		}
	}

	return limit, nil
}
