package main

import (
	"bufio"
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

func main() {

	start := time.Now()

	excelFileName := "test.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	bodyPath := "call-bodies"
	if err != nil {
		fmt.Println(err)
	}

	var wg sync.WaitGroup

	fmt.Println("Making calls")

	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {

			wg.Add(1)

			name, _ := row.Cells[0].String()
			httpType, _ := row.Cells[1].String()
			url, _ := row.Cells[2].String()
			relativePath, _ := row.Cells[3].String()
			bodyPath := bodyPath + "/" + relativePath
			fmt.Print(".")
			go makeHttpCall(name, httpType, url, bodyPath, &wg)
		}
	}

	fmt.Println("Waiting..")
	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)

	go startWS()

	time.Sleep(60 * time.Second)
	fmt.Println("This is working")
}

func makeHttpCall(name string, httpType string, urll string, bodyPath string, wg *sync.WaitGroup) {

	baseUrl := "http://10.252.169.12:7788" + urll
	// Need to read values from file with path in bodyPath as JSON
	bodyValue, _ := os.Open(bodyPath)
	bodyVal := bufio.NewReader(bodyValue)
	var err error
	var resp *http.Response

	switch httpType {
	case "GET":
		resp, err := http.Get(baseUrl)
	case "POST":
		resp, err := http.Post(baseUrl, "application/json", bodyVal)
	}

	if resp != nil {

		if resp.StatusCode == 200 {

			fmt.Print("+")
		} else {

			fmt.Println(baseUrl, resp.StatusCode, err)
		}
	} else {

		fmt.Println("Didn't make call...", baseUrl)
	}

	wg.Done()
}

func startWS() {
	fs := http.FileServer(http.Dir("static-pages"))
	http.Handle("/", fs)

	log.Println("Listening...")

	go openBrowser("http://127.0.0.1:3000/index.html")

	http.ListenAndServe(":3000", nil)
}

// openBrowser tries to open the URL in a browser,
// and returns whether it succeed in doing so.
func openBrowser(url string) bool {

	time.Sleep(3 * time.Second)

	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}
