package pool

import (
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestPool_NewTick(t *testing.T) {
	startTime := time.Now()
	rand.Seed(time.Now().UnixNano())
	fmt.Println(strconv.FormatInt(rand.Int63(), 10))
	elapsed := time.Since(startTime)
	fmt.Println(elapsed)
}

func TestParse(t *testing.T) {
	fmt.Println(url.Parse("114.114.114.114:8080"))
}
