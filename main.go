package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/danieldiamont/go-yamcs-cli/yamcs/protobuf"
	"github.com/danieldiamont/go-yamcs-cli/yamcs/protobuf/events"
	pvalue "github.com/danieldiamont/go-yamcs-cli/yamcs/protobuf/pvalue"
)

type StreamParameterRequestDetails struct {
	Uri      string
	Instance string
	Start    string
	Stop     string
	Ids      []ID
}

type StreamParameterValuesRequest struct {
	Start   string   `json:"start"`
	Stop    string   `json:"stop"`
	Ids     []ID     `json:"ids"`
	TmLinks []string `json:"tmLinks,omitempty"`
}

type ID struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type StreamEventsRequestDetails struct {
	Uri      string
	Instance string
	Start    string
	Stop     string
	Source   []string
	Severity string
	Q        string
}

type StreamEventsRequest struct {
	Start    string   `json:"start"`
	Stop     string   `json:"stop"`
	Source   []string `json:"source,omitempty"`
	Severity string   `json:"severity,omitempty"`
	Q        string   `json:"q,omitempty"`
}

func splitDates(dates string, sep string) (string, string) {
	split := strings.Split(dates, sep)
	return split[0], split[1]
}

func getRequest(uri string) []byte {

	res, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody : %s\n", res.StatusCode, body)
	}

	return body
}

func fixJsonResp(jsonBytes []byte) ([][]byte, error) {

	var jsonBytesSlice [][]byte
	jsonChunks := strings.Split(string(jsonBytes), "}{")

	if len(jsonChunks) > 0 {
		for _, chunk := range jsonChunks {
			if len(chunk) > 0 {
				if chunk[0] == '{' && chunk[len(chunk)-1] != '}' {
					chunk += "}"
				} else if chunk[len(chunk)-1] == '}' && chunk[0] != '{' {
					chunk = "{" + chunk
				} else if chunk[0] != '{' && chunk[len(chunk)-1] != '}' {
					chunk = "{" + chunk + "}"
				} else {
					return nil, fmt.Errorf("Unexpected JSON formatting: %s", chunk)
				}
				jsonBytesSlice = append(jsonBytesSlice, []byte(chunk))
			}
		}
	}
	return jsonBytesSlice, nil
}

func configureStreamEventsRequest(requestDetails *StreamEventsRequestDetails) (*http.Request, error) {
	requestUri := fmt.Sprintf(
		"%s/stream-archive/%s:streamEvents",
		requestDetails.Uri,
		requestDetails.Instance,
	)

	streamPayload := StreamEventsRequest{
		Start:    requestDetails.Start,
		Stop:     requestDetails.Stop,
		Source:   requestDetails.Source,
		Severity: requestDetails.Severity,
		Q:        requestDetails.Q,
	}

	streamPayloadBytes, err := json.Marshal(streamPayload)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(streamPayloadBytes)

	req, err := http.NewRequest("POST", requestUri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func processEventStream(pdChan chan<- *events.Event, body io.ReadCloser) error {

	reader := bufio.NewReader(body)
	var buffer bytes.Buffer

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				close(pdChan)
				break
			}
			return err
		}

		buffer.Write(line)
		fixedJsonChunks, err := fixJsonResp(buffer.Bytes())

		for _, chunk := range fixedJsonChunks {
			if json.Valid(chunk) {
				var event events.Event
				if err := json.Unmarshal(chunk, &event); err != nil {
					log.Printf("Error unmarshaling event data: %v", err)
				} else {
					if event.String() != "" {
						pdChan <- &event
					}
				}
				buffer.Reset()
			}
		}
	}
	return nil
}

func streamEvents(pdChan chan<- *events.Event, errChan chan<- error, requestDetails *StreamEventsRequestDetails) {

	req, err := configureStreamEventsRequest(requestDetails)
	if err != nil {
		log.Printf("Error configuring Event stream request")
		errChan <- err
		close(pdChan)
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error executing Event stream request.")
		errChan <- err
		close(pdChan)
		return
	}
	defer res.Body.Close()

	err = processEventStream(pdChan, res.Body)
	if err != nil {
		log.Printf("Error processing Event stream")
		errChan <- err
		close(pdChan)
		return
	}

	close(errChan)
}

func configureStreamParameterValuesRequest(requestDetails *StreamParameterRequestDetails) (*http.Request, error) {

	requestUri := fmt.Sprintf(
		"%s/stream-archive/%s:streamParameterValues",
		requestDetails.Uri,
		requestDetails.Instance,
	)

	streamPayload := StreamParameterValuesRequest{
		Start: requestDetails.Start,
		Stop:  requestDetails.Stop,
		Ids:   requestDetails.Ids,
	}

	streamPayloadBytes, err := json.Marshal(streamPayload)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(streamPayloadBytes)

	req, err := http.NewRequest("POST", requestUri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func processParameterValueStream(pdChan chan<- *pvalue.ParameterData, body io.ReadCloser) error {

	reader := bufio.NewReader(body)
	var buffer bytes.Buffer

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				close(pdChan)
				break
			}
			return err
		}

		buffer.Write(line)
		fixedJsonChunks, err := fixJsonResp(buffer.Bytes())

		for _, chunk := range fixedJsonChunks {
			if json.Valid(chunk) {
				var paramData pvalue.ParameterData
				if err := json.Unmarshal(chunk, &paramData); err != nil {
					log.Printf("Error unmarshaling ParameterData: %v", err)
				} else {
					if paramData.String() != "" {
						pdChan <- &paramData
					}
				}
				buffer.Reset()
			}
		}
	}
	return nil
}

func streamParameterValues(pdChan chan<- *pvalue.ParameterData, errChan chan<- error, requestDetails *StreamParameterRequestDetails) {

	req, err := configureStreamParameterValuesRequest(requestDetails)
	if err != nil {
		log.Printf("Error configuring ParameterValue stream request")
		errChan <- err
		close(pdChan)
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error executing ParameterValue stream request.")
		errChan <- err
		close(pdChan)
		return
	}
	defer res.Body.Close()

	err = processParameterValueStream(pdChan, res.Body)
	if err != nil {
		log.Printf("Error processing ParameterValue stream")
		errChan <- err
		close(pdChan)
		return
	}

	close(errChan)
}

func paramToCsv(param *pvalue.ParameterValue) {
	filename := *param.Id.Name + ".csv"
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening file %s", filename)
	}
	defer f.Close()

	outStr := paramToString(param)

	_, err = f.WriteString(outStr)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func paramToString(param *pvalue.ParameterValue) string {

	engValue := param.GetEngValue()

	var paths []string
	var vals []*protobuf.Value

	unrollAggregate(engValue, param.GetId().GetName(), &paths, &vals)

	engValueStr := "values=["

	for i, val := range vals {

		switch val.GetType() {
		case protobuf.Value_SINT64:
			engValueStr = fmt.Sprintf("%s(%s:%v)|", engValueStr, paths[i], val.GetSint64Value())
		case protobuf.Value_SINT32:
			engValueStr = fmt.Sprintf("%s(%s:%v)|", engValueStr, paths[i], val.GetSint32Value())
		case protobuf.Value_FLOAT:
			engValueStr = fmt.Sprintf("%s(%s:%v)|", engValueStr, paths[i], val.GetFloatValue())
		default:
			engValueStr = fmt.Sprintf("%s(%s:%v)|", engValueStr, paths[i], val)
		}

	}

	if engValueStr[len(engValueStr)-1] == '[' {
		engValueStr += "]"
	} else {
		engValueStr = engValueStr[:len(engValueStr)-1] + "]"
	}

	outStr := fmt.Sprintf(
		"%s,%s,%s,%s,%s,%s\n",
		param.GetGenerationTime(),
		param.GetAcquisitionTime(),
		param.GetId().GetNamespace(),
		param.GetId().GetName(),
		engValueStr,
		param.GetRawValue(),
	)
	return outStr
}

func unrollAggregate(val *protobuf.Value, path string, paths *[]string, vals *[]*protobuf.Value) {

	if val.GetType() != *protobuf.Value_AGGREGATE.Enum() {
		*paths = append(*paths, path)
		*vals = append(*vals, val)
		return
	}

	if val.GetAggregateValue() != nil {
		for i, value := range val.AggregateValue.Value {
			name := val.AggregateValue.Name[i]
			unrollAggregate(value, fmt.Sprintf("%s.%s", path, name), paths, vals)
		}
	}

}

func initParameterDataCsvFiles(ids []ID) {

	header := fmt.Sprintf(
		"%s,%s,%s,%s,%s,%s\n",
		"GenerationTime",
		"AcquisitionTime",
		"Namespace",
		"Name",
		"EngValue",
		"RawValue",
	)

	for _, id := range ids {
		fp := id.Name + ".csv"

		err := os.Remove(fp)
		if !os.IsNotExist(err) && err != nil {
			log.Fatalf("Error deleting file %s with error: %q\n", fp, err)
		}

		f, err := os.Create(fp)
		if err != nil {
			log.Fatalf("Error creating file %s with error: %q\n", fp, err)
		}

		f.WriteString(header)
		f.Close()
	}
}

func initEventsCsvFile(filename string) {
    
	header := fmt.Sprintf(
		"%s,%s,%s,%s,%s,%s,%s\n",
        "Source",
		"GenerationTime",
		"ReceptionTime",
		"SeqNum",
		"Message",
		"Severity",
		"CreatedBy",
    )

    err := os.Remove(filename)
    if !os.IsNotExist(err) && err != nil {
        log.Fatalf("Error deleting file %s with error: %q\n", filename, err)
    }

    f, err := os.Create(filename)
    if err != nil {
        log.Fatalf("Error creating file %s with error: %q\n", filename, err)
    }

    f.WriteString(header)
}

func eventToString(e *events.Event) string {
    eventStr := fmt.Sprintf(
        "%s, %s, %s, %d, %s, %s, %s\n",
        e.GetSource(),
        e.GetGenerationTime(),
        e.GetReceptionTime(),
        e.GetSeqNumber(),
        e.GetMessage(),
        e.GetSeverity(),
        e.GetCreatedBy(),
    )

    return eventStr
}

func eventToCsv(e *events.Event, filename *string) {
    
	f, err := os.OpenFile(*filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening file %s", *filename)
	}
	defer f.Close()

	_, err = f.WriteString(eventToString(e))
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func main() {

	dates := flag.String(
		"dates",
		"2023-12-03T23:46:00Z..2023-12-03T23:50:59Z",
		"date range",
	)
	yamcsURI := flag.String("uri", "http://localhost:8090/api", "yamcs server URI")
	instance := flag.String("instance", "myproject", "yamcs server instance name")
	paramFile := flag.String("paramFile", "params.json", "list of parameters in .json format")
    eventsOut := flag.String("eventsOut", "events.csv", "output path (relative) for events")
	verbose := flag.Bool("verbose", false, "print log statements to console")

	flag.Parse()

	start, stop := splitDates(*dates, "..")
	if *verbose {
		log.Printf("%q, %q\n", start, stop)
	}

	res := getRequest(*yamcsURI)
	if *verbose {
		log.Printf("res = %s\n", res)
	}

	dat, err := os.ReadFile(*paramFile)
	if err != nil {
		log.Fatalf("Error reading %s\n%q", *paramFile, err)
	}
	if !json.Valid(dat) {
		log.Fatalf("%s is not valid JSON\n", *paramFile)
	}

	var paramList []ID
	err = json.Unmarshal(dat, &paramList)
	if err != nil {
		log.Fatalf("Error de-serializing contents of %s\n", *paramFile)
	}

	initParameterDataCsvFiles(paramList)

	// stream parameterData
	paramDataRequestDetails := StreamParameterRequestDetails{
		Uri:      *yamcsURI,
		Instance: *instance,
		Start:    start,
		Stop:     stop,
		Ids:      paramList,
	}

	parameterDataChannel := make(chan *pvalue.ParameterData)
	parameterDataErrorChannel := make(chan error, 1)

	go streamParameterValues(
		parameterDataChannel,
		parameterDataErrorChannel,
		&paramDataRequestDetails,
	)

	for paramData := range parameterDataChannel {
		if *verbose {
			fmt.Printf("Received ParameterData:\n%q\n", paramData)
		}
		for _, param := range paramData.Parameter {
			paramToCsv(param)
		}
	}

	err = <-parameterDataErrorChannel
	if err != nil {
		log.Printf("Error occurred during ParameterData Stream: %q", err)
		close(parameterDataErrorChannel)
	}

	// stream events
    initEventsCsvFile(*eventsOut)

	eventsRequestDetails := StreamEventsRequestDetails{
		Uri:      *yamcsURI,
		Instance: *instance,
		Start:    start,
		Stop:     stop,
	}

	eventsChannel := make(chan *events.Event)
	eventsErrorChannel := make(chan error, 1)

	go streamEvents(
		eventsChannel,
		eventsErrorChannel,
		&eventsRequestDetails,
	)

	for event := range eventsChannel {
		if *verbose {
			fmt.Printf("Received Event:\n%q\n", event)
		}
        eventToCsv(event, eventsOut)
	}

	err = <-eventsErrorChannel
	if err != nil {
		log.Printf("Error occurred during event Stream: %q", err)
		close(eventsErrorChannel)
	}
}
