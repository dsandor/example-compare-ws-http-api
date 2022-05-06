package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/k0kubun/pp/v3"
	"gonum.org/v1/gonum/stat"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"
)

type BenchmarkResults struct {
	ConnectAndAuth  float64
	InitialCall     float64
	SubsequentCalls []int // each call time is in milliseconds
}

type BenchmarkMetrics struct {
	MeanCallDuration   float64
	MedianCallDuration float64
	MinCallDuration    float64
	MaxCallDuration    float64
	Variance           float64
	StandardDeviation  float64
}

const numberOfCalls = 1000

func main() {
	log.Printf("Running websocket tests, %d iterations.\n", numberOfCalls+1)
	wsResults := performWebsocketBenchmark()
	_, _ = pp.Println(wsResults.GetMetrics())
	_ = wsResults.WriteJSON("websocket_performance_100ms.json")
	_ = wsResults.WriteCSV("websocket_performance_call_times_100ms_auth_delay.csv", "call_time")

	log.Printf("Running http api tests, %d iterations.\n", numberOfCalls+1)
	httpResults := performHttpApiBenchmark()
	_, _ = pp.Println(httpResults.GetMetrics())
	_ = httpResults.WriteJSON("httpapi_performance_100ms.json")
	_ = httpResults.WriteCSV("httpapi_performance_call_times_100ms_auth_delay.csv", "call_time")
}

func performWebsocketBenchmark() BenchmarkResults {
	benchmarkResults := BenchmarkResults{}

	connectionAndAuthTime := time.Now()

	conn, _, _, err := ws.DefaultDialer.Dial(context.TODO(), "wss://3heeykqcg9.execute-api.us-east-1.amazonaws.com/v1?token=1234")

	if err != nil {
		log.Fatalf("Error dialing websocket: %s", err.Error())
	}

	elapsedConnectionAndAuthTime := time.Since(connectionAndAuthTime)
	benchmarkResults.ConnectAndAuth = float64(elapsedConnectionAndAuthTime.Milliseconds())

	sendInitialMessageTime := time.Now()

	err = wsutil.WriteClientMessage(conn, ws.OpText, []byte("Test"))

	if err != nil {
		log.Fatalf("Error writing client message: %s", err.Error())
	}

	_, err = wsutil.ReadServerText(conn)

	if err != nil {
		log.Fatalf("Error reading client message: %s", err.Error())
	}

	sendInitialMessageRoundTripTime := time.Since(sendInitialMessageTime)
	benchmarkResults.InitialCall = float64(sendInitialMessageRoundTripTime.Milliseconds())

	var callResultTiming [numberOfCalls]int

	for i := 0; i < numberOfCalls; i++ {
		sendMessageTime := time.Now()

		err = wsutil.WriteClientMessage(conn, ws.OpText, []byte("Test"))

		if err != nil {
			log.Fatalf("Error writing client message: %s", err.Error())
		}

		_, err := wsutil.ReadServerText(conn)

		if err != nil {
			log.Fatalf("Error reading client message: %s", err.Error())
		}

		sendMessageRoundTripTime := time.Since(sendMessageTime)
		callResultTiming[i] = int(sendMessageRoundTripTime.Milliseconds())
	}

	benchmarkResults.SubsequentCalls = callResultTiming[:]

	return benchmarkResults
}

func performHttpApiBenchmark() BenchmarkResults {
	benchmarkResults := BenchmarkResults{}

	benchmarkResults.ConnectAndAuth = 0 // this is not applicable for http

	sendInitialMessageTime := time.Now()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://g6p6yapybd.execute-api.us-east-1.amazonaws.com/hello", nil)
	req.Header.Set("Authorization", "1234")

	_, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	sendInitialMessageRoundTripTime := time.Since(sendInitialMessageTime)
	benchmarkResults.InitialCall = float64(sendInitialMessageRoundTripTime.Milliseconds())

	var callResultTiming [numberOfCalls]int

	for i := 0; i < numberOfCalls; i++ {
		sendMessageTime := time.Now()

		req, _ := http.NewRequest("GET", "https://g6p6yapybd.execute-api.us-east-1.amazonaws.com/hello", nil)
		req.Header.Set("Authorization", "1234")

		_, err := client.Do(req)

		if err != nil {
			log.Fatalf("Error reading client message: %s", err.Error())
		}

		sendMessageRoundTripTime := time.Since(sendMessageTime)
		callResultTiming[i] = int(sendMessageRoundTripTime.Milliseconds())
	}

	benchmarkResults.SubsequentCalls = callResultTiming[:]

	return benchmarkResults
}

func (bm *BenchmarkResults) WriteJSON(filename string) error {
	jsonString, _ := json.MarshalIndent(bm, "", " ")

	err := ioutil.WriteFile(filename, jsonString, 0644)

	return err
}

func (bm *BenchmarkResults) WriteCSV(filename string, columnName string) error {
	csvString := fmt.Sprintf("%s\n", columnName)

	for _, timing := range bm.SubsequentCalls {
		csvString += fmt.Sprintf("%d\n", timing)
	}

	err := ioutil.WriteFile(filename, []byte(csvString), 0644)

	return err
}

func (bm *BenchmarkResults) GetMetrics() BenchmarkMetrics {
	sortedSubsequentCalls := make([]float64, len(bm.SubsequentCalls))
	var subsequentCallsFloat []float64
	for _, val := range bm.SubsequentCalls {
		subsequentCallsFloat = append(subsequentCallsFloat, float64(val))
	}

	copy(sortedSubsequentCalls, subsequentCallsFloat)

	sort.Float64s(sortedSubsequentCalls)

	metrics := BenchmarkMetrics{
		MeanCallDuration:   stat.Mean(subsequentCallsFloat, nil),
		MedianCallDuration: stat.Quantile(0.5, stat.Empirical, sortedSubsequentCalls, nil),
		MinCallDuration:    sortedSubsequentCalls[0],
		MaxCallDuration:    sortedSubsequentCalls[len(sortedSubsequentCalls)-1],
		Variance:           stat.Variance(subsequentCallsFloat, nil),
		StandardDeviation:  stat.StdDev(subsequentCallsFloat, nil),
	}

	return metrics
}
