package main

import (
    "bytes"
    "fmt"
    "flag"
    "strings"
    "net/http"
    "encoding/json"
    "log"
    "io"
)

type RequestDetails struct {
    Uri         string
    Instance    string
    Start       string
    Stop        string
    Ids         []ID
}

type StreamParameterValuesRequest struct {
    Start   string  `json:"start"`
    Stop    string  `json:"stop"`
    Ids     []ID    `json:"ids"`
    TmLinks []string`json:"tmLinks,omitempty"`
}

type ID struct {
    Namespace   string `json:"namespace"`
    Name        string `json:"name"`
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
            
            if chunk[0] == '{' && chunk[len(chunk)-1] != '}' {
                chunk += "}"
            } else if chunk[len(chunk)-1] == '}' && chunk[0] != '{' {
                chunk = "{" + chunk
            } else if chunk[0] != '{' && chunk[len(chunk)-1] != '}' {
                chunk = "{" + chunk + "}"
            } else {
                return nil, fmt.Errorf("Unexpected JSON formatting: %s", chunk)
            }
            fmt.Println(chunk)

            jsonBytesSlice = append(jsonBytesSlice, []byte(chunk))
        }
    }

    return jsonBytesSlice, nil
}

func streamParameterValues( requestDetails *RequestDetails) {

        requestUri := fmt.Sprintf(
            "%s/stream-archive/%s:streamParameterValues",
            requestDetails.Uri,
            requestDetails.Instance,
        )

        streamPayload := StreamParameterValuesRequest{
            Start: requestDetails.Start,
            Stop: requestDetails.Stop,
            Ids: requestDetails.Ids,
        }

        streamPayloadBytes, err := json.Marshal(streamPayload)
        if err != nil {
            log.Fatal(err)
        }
        log.Printf("Marshalled payload:\n\n%s\n\n", streamPayloadBytes)
        body := bytes.NewReader(streamPayloadBytes)

        req, err := http.NewRequest("POST", requestUri, body)
        if err != nil {
            log.Fatal(err)
        }

        req.Header.Set("Content-Type", "application/json")

        client := &http.Client{}        
        res, err := client.Do(req)
        if err != nil {
            log.Fatal(err)
        }

        defer res.Body.Close()

        resBody, err := io.ReadAll(res.Body)
        if err != nil {
            log.Fatal(err)
        }
        if res.StatusCode > 299 {
            log.Fatalf("Response failed with status code: %d and\nbody : %s\n", res.StatusCode, resBody)
        }
        
        fmt.Printf("%q\n", resBody)
        
        cleanedJsonChunks, err := fixJsonResp(resBody)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(cleanedJsonChunks)

        var jsonResp map[string]interface{}

        for _, cleanedJson := range cleanedJsonChunks {

            err = json.Unmarshal(cleanedJson, &jsonResp)
            if err != nil {
                log.Println("Could not unmarshal json: %s\n", err)
            }
            for k := range jsonResp {
                fmt.Println(k)
            }

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

    flag.Parse()

    start, stop := splitDates(*dates, "..")
    log.Printf("%q, %q\n", start, stop)

    res := getRequest(*yamcsURI)
    log.Printf("res = %s\n", res)

    streamRequestDetails := RequestDetails {
            Uri: *yamcsURI,
            Instance: *instance,
            Start: start,
            Stop: stop,
            Ids: []ID{
                {
                    Name: "jvmMemoryUsed",
                    Namespace: "/yamcs/spacestation",
                },
            },
    }

    streamParameterValues(&streamRequestDetails)

}
