package flightaware

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"sync"
	"testing"
)

const tabCount = 10

//go:embed flights_test_data.txt
var fPayload []byte

func TestGetFlightInfo(t *testing.T) {
	fChan := make(chan string, tabCount)
	errCount := 0
	count := 0

	go func() {
		s := bufio.NewScanner(bytes.NewReader(fPayload))
		for s.Scan() {
			if count > 50 {
				close(fChan)
				return
			}
			number := strings.Replace(s.Text(), "SU", "AFL", 1)
			fChan <- number
			count++
		}
		close(fChan)
	}()
	log := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)
	api, err := NewFlightInfoAPI(&ApiConfig{
		MaxTabCount: tabCount,
		Debug:       true,
	}, &log)
	if err != nil {
		t.Fatalf("create api error: %e", err)
	}

	wg := sync.WaitGroup{}
	ctx := context.Background()
	for i := 0; i < tabCount; i++ {
		wg.Add(1)
		go func() {
			for n := range fChan {
				_, err := api.GetFlightInfo(ctx, n)
				if err != nil {
					errCount++
					t.Errorf("get flight %s error: %e", n, err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	t.Logf("Error pecrent: %f", float64(errCount)*100/float64(count))
}
