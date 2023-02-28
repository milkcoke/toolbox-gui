package file

import (
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/inhies/go-bytesize"
	"github.com/milkcoke/auto-setup-gui/src/app"
	"io"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

func retryRequest(readFileFD *os.File) {
	retryFileInfo, err := readFileFD.Stat()
	if err != nil {
		log.Fatalln("Failed to open test file description")
	}

	log.Println("Retrying file size : ", bytesize.New(float64(retryFileInfo.Size())))

	streamFile, err := os.OpenFile(readFileFD.Name(), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Failed to open readFileFD : ", err)
	}
	defer streamFile.Close()

	res, err := req.R().
		SetHeader("Range", fmt.Sprintf("bytes=%d-", retryFileInfo.Size())).
		Get(app.DockerInstaller.Url)

	defer res.Body.Close()

	written, copyErr := io.Copy(streamFile, res.Body)
	log.Println("Written : ", bytesize.New(float64(written)))

	if err != nil {
		log.Println("Http request fail : ", err)
		time.Sleep(2 * time.Second)
		retryRequest(readFileFD)
	}

	if copyErr != nil {
		log.Println("Error occurs appending stream from response .. ", err)
		time.Sleep(2 * time.Second)
		retryRequest(readFileFD)
	}

}

func Test_Retry_Download_Complete(t *testing.T) {

	testFileFullPath, err := getCompleteFileFullPath()

	if err != nil {
		t.Fatal("Failed to current directory path")
	}

	// 테스트 파일 없음
	if _, err := os.Stat(testFileFullPath); err != nil {
		t.Fatal("Test completeTestFile notFound.")
	}

	readFD, err := os.Open(testFileFullPath)
	if err != nil {
		t.Fatal("Failed to open test readFD ", err)
	}

	defer readFD.Close()

	fileInfo, err := readFD.Stat()
	if err != nil {
		t.Fatal("Failed to get test readFD info ", err)
	}

	headerRes, err := req.R().Head(app.DockerInstaller.Url)
	if err != nil {
		log.Println("헤더 응답 오류 ", err)
	}

	// Convert string to int64
	contentLength := headerRes.GetHeader("Content-Length")
	fullFileLength, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		log.Println("파일 크기 변환 에러 ", err)
	}
	log.Println("Total file size: ", bytesize.New(float64(fullFileLength)))
	log.Println("Current file size: ", bytesize.New(float64(fileInfo.Size())))

	// 이미 완성본이 된 경우
	if fullFileLength == fileInfo.Size() {
		// 이미 완성된 파일이라면, 해당 파일을 열어주고 끝내야함.
		// 단, 이벤트 등록은 여전히!
		if err != nil {
			log.Println("Invalid directory path")
		}
		log.Println(completeTestFile + " already exists!")
	} else {
		retryRequest(readFD)
		log.Println("After Retry : " + completeTestFile + " download complete.")
	}

	NavigateToDir(testFileFullPath)
}

func Test_Retry_Download_Partial(t *testing.T) {

	testFileFullPath, err := getPartialFileFullPath()

	if err != nil {
		t.Fatal("Failed to current directory path")
	}

	// 테스트 파일 없음
	if _, err := os.Stat(testFileFullPath); err != nil {
		t.Fatal("Test completeTestFile notFound.")
	}

	readFD, err := os.Open(testFileFullPath)
	if err != nil {
		t.Fatal("Failed to open test readFD ", err)
	}

	defer readFD.Close()

	fileInfo, err := readFD.Stat()
	if err != nil {
		t.Fatal("Failed to get test readFD info ", err)
	}

	headerRes, err := req.R().Head(app.DockerInstaller.Url)
	if err != nil {
		log.Println("헤더 응답 오류 ", err)
	}

	// Convert string to int64
	contentLength := headerRes.GetHeader("Content-Length")
	fullFileLength, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		log.Println("파일 크기 변환 에러 ", err)
	}
	log.Println("Total file size: ", bytesize.New(float64(fullFileLength)))
	log.Println("Current file size: ", bytesize.New(float64(fileInfo.Size())))

	// 이미 완성본이 된 경우
	if fullFileLength == fileInfo.Size() {
		// 이미 완성된 파일이라면, 해당 파일을 열어주고 끝내야함.
		// 단, 이벤트 등록은 여전히!
		if err != nil {
			log.Println("Invalid directory path")
		}
		log.Println(partialTestFile + " already exists!")
	} else {
		retryRequest(readFD)
		log.Println("After Retry : " + partialTestFile + " download complete.")
	}

	NavigateToDir(testFileFullPath)
}
