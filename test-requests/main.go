package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Starting test requests")

	seconds := 0

	for {
		for i := 0; i < 100; i++ {
			r := rand.Intn(15)
			if r < 5 {
				r = 5
			}
			for j := 0; j < r; j++ {
				testQuery(i)
			}
		}
		time.Sleep(time.Second)
		seconds++
		fmt.Println("Testing ", seconds)

	}

}

func testQuery(bannerID int) {
	url := fmt.Sprintf("http://localhost:5004/api/counter/%d", bannerID)

	_, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка выполнения GET-запроса:", err)
		return
	}
}
