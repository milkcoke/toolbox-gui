package layout

import (
	"fmt"
	"github.com/imroc/req/v3"
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
		log.Fatalln("아니 이것도 왜 안열리냐 ㅋㅋ")
	}
	// TODO
	//  여기서 SetHeader 를 갱신하고 재시도하는 로직이 필요함.

	log.Println("재시도중 현재 파일 크기 : ", retryFileInfo.Size())

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
	log.Println("써진 바이트 : ", written)

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

func Test_Retry_Download(t *testing.T) {
	testFile := "Docker.exe"

	// 테스트 파일 없음
	if _, err := os.Stat(testFile); err != nil {
		t.Fatal("Test readFD notFound.")
	}

	readFD, err := os.Open(testFile)
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
	log.Println("전체 파일 길이: ", contentLength)
	log.Println("현재 파일 크기: ", fileInfo.Size())

	// 이미 완성본이 된 경우
	if fullFileLength == fileInfo.Size() {
		// 이미 완성된 파일이라면, 해당 파일을 열어주고 끝내야함.
		// 단, 이벤트 등록은 여전히!
		if err != nil {
			log.Println("잘못된 디렉토리 경로")
		}
		log.Println(testFile + " download complete")
	} else {
		retryRequest(readFD)
		log.Println("After Retry : " + testFile + " download complete.")
	}

	os.Chdir(testFile)
}
