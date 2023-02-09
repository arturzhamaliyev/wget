package transporter

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

func Download(URL, fileName string, output *os.File) error {
	fmt.Fprintf(output, "start at %v\n", time.Now().Format("2006-01-02 15:04:05"))

	fmt.Fprint(output, "sending request, awaiting response... status ")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MaxVersion: tls.VersionTLS13,
		},
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Chuhuahua.akl/pidor")

	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.Status != "200 OK" {
		return errors.New(response.Status)
	}
	fmt.Fprintln(output, response.Status)

	// utils.ContentSizeCheck(response.ContentLength)
	fmt.Fprintf(output, "content size: %d [~%.2fMB]\n", response.ContentLength, float64(response.ContentLength)/1000000)

	fmt.Fprintf(output, "saving file to: ./%s\n", fileName)

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	if output == os.Stdout {
		template := `{{ counters .}} {{ bar . "[" "=" (cycle . ">" ) "." "]"}} {{percent .}} {{speed .}} {{rtime .}}`

		bar := pb.ProgressBarTemplate(template).Start64(response.ContentLength)

		barReader := bar.NewProxyReader(response.Body)

		_, err = io.Copy(file, barReader)
		if err != nil {
			return err
		}

		bar.Finish()
	} else {
		_, err = io.Copy(file, response.Body)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(output, "Downloaded [%s]\n", URL)
	fmt.Fprintf(output, "finished at %v\n", time.Now().Format("2006-01-02 15:04:05"))

	// fmt.Println(response.Header)

	return nil
}
