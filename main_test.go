package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/danieldiamont/go-yamcs-cli/yamcs/protobuf"
	"github.com/danieldiamont/go-yamcs-cli/yamcs/protobuf/events"
	pvalue "github.com/danieldiamont/go-yamcs-cli/yamcs/protobuf/pvalue"
)

func TestInitEventCsvFile(t *testing.T) {
    filename := "test-events-header.csv.csv"
    expectedFile := filename

	expectedHeader := fmt.Sprintf(
		"%s,%s,%s,%s,%s,%s,%s",
        "Source",
		"GenerationTime",
		"ReceptionTime",
		"SeqNum",
		"Message",
		"Severity",
		"CreatedBy",
    )

	initEventsCsvFile(filename)

    f, err := os.Open(expectedFile)
    if os.IsNotExist(err) {
        t.Errorf("Expected file %s to exist\n%q", expectedFile, err)
    } else if err != nil {
        t.Errorf("Error opening file: %q\n", err)
    } else {
        defer f.Close()

        scanner := bufio.NewScanner(f)
        if scanner.Scan() {
            line := scanner.Text()
            if line != expectedHeader {
                t.Errorf("got %s want %s\n", line, expectedHeader)
            }
        }

        if err = scanner.Err(); err != nil {
            t.Errorf("Error reading the file: %q\n", err)
		}

	}
}

func TestInitParameterDataCsvFiles(t *testing.T) {

	f, err := os.Create("tmp.csv")
	if err != nil {
		t.Fatalf("Error creating tmp file\n")
	}
	f.WriteString("abcdefg")
	f.Close()

	paramListJson := `
		[
			{
				"namespace": "/yamcs/spacestation",
				"name":      "jvmMemoryUsed"
			},
			{
				"namespace": "/yamcs/spacestation/df/dev",
				"name":      "sdc"
			}
		]
	`

	var paramList []ID
	err = json.Unmarshal([]byte(paramListJson), &paramList)
	if err != nil {
		t.Fatalf("Error unmarshalling JSON\n")
	}

	var expectedFiles []string
	for _, param := range paramList {
		expectedFiles = append(expectedFiles, fmt.Sprintf("%s.csv", param.Name))
	}

	expectedHeader := fmt.Sprintf(
		"%s,%s,%s,%s,%s,%s",
		"GenerationTime",
		"AcquisitionTime",
		"Namespace",
		"Name",
		"EngValue",
		"RawValue",
	)

	initParameterDataCsvFiles(paramList)

	for _, expectedFile := range expectedFiles {
		f, err = os.Open(expectedFile)
		if os.IsNotExist(err) {
			t.Errorf("Expected file %s to exist\n%q", expectedFile, err)
		} else if err != nil {
			t.Errorf("Error opening file: %q\n", err)
		} else {
			defer f.Close()

			scanner := bufio.NewScanner(f)
			if scanner.Scan() {
				line := scanner.Text()
				if line != expectedHeader {
					t.Errorf("got %s want %s\n", line, expectedHeader)
				}
			}

			if err = scanner.Err(); err != nil {
				t.Errorf("Error reading the file: %q\n", err)
			}

		}
	}

}

func TestParamToStringSint64ValueEngValue(t *testing.T) {

	namespace := "/yamcs/spacestation"
	name := "jvmMemoryUsed"
	gen_time := "2023-12-03T23:47:34.942Z"
	ac_time := "2023-12-03T23:47:37.942Z"
	sint64value := "58783"

	want := fmt.Sprintf(
		"%s,%s,%s,%s,values=[(%s:%v)],%v\n",
		gen_time,
		ac_time,
		namespace,
		name,
		name,
		sint64value,
		nil,
	)

	objId := protobuf.NamedObjectId{}
	objId.Name = &name
	objId.Namespace = &namespace

	engValProto := protobuf.Value{}
	engValProto.Type = protobuf.Value_SINT64.Enum()
	engValProto.Sint64Value = &sint64value

	param := pvalue.ParameterValue{
		Id:              &objId,
		EngValue:        &engValProto,
		RawValue:        nil,
		AcquisitionTime: &ac_time,
		GenerationTime:  &gen_time,
	}

	got := paramToString(&param)
	if got != want {
		t.Errorf("Got %s, want %s", got, want)
	}
}

func TestParamToStringAggregateValue(t *testing.T) {

	namespace := "/yamcs/spacestation/df/dev"
	name := "sdc"
	gen_time := "2023-12-03T23:47:34.942Z"
	ac_time := "2023-12-03T23:47:37.942Z"

	totalVal := "10557621868"
	availableVal := "1034958734"
	var percentageUseVal float32
	percentageUseVal = 5.646844

	totalValProto := protobuf.Value{}
	totalValProto.Type = protobuf.Value_SINT64.Enum()
	totalValProto.Sint64Value = &totalVal
	availableValProto := protobuf.Value{}
	availableValProto.Type = protobuf.Value_SINT64.Enum()
	availableValProto.Sint64Value = &availableVal

	percentageUseValProto := protobuf.Value{}
	percentageUseValProto.Type = protobuf.Value_FLOAT.Enum()
	percentageUseValProto.FloatValue = &percentageUseVal

	lvl2aggValProto := protobuf.AggregateValue{}
	lvl2aggValProto.Name = []string{"percentageUse"}
	lvl2aggValProto.Value = []*protobuf.Value{&percentageUseValProto}

	lvl2aggVal := protobuf.Value{}
	lvl2aggVal.Type = protobuf.Value_AGGREGATE.Enum()
	lvl2aggVal.AggregateValue = &lvl2aggValProto

	aggValProto := protobuf.AggregateValue{}
	aggValProto.Name = []string{"total", "available", "lvl2agg"}
	aggValProto.Value = []*protobuf.Value{
		&totalValProto,
		&availableValProto,
		&lvl2aggVal,
	}

	aggValueStr := fmt.Sprintf(
		"values=[(%s:%v)|(%s:%v)|(%s:%v)]",
		"sdc.total",
		totalVal,
		"sdc.available",
		availableVal,
		"sdc.lvl2agg.percentageUse",
		percentageUseVal,
	)

	want := fmt.Sprintf(
		"%s,%s,%s,%s,%s,%v\n",
		gen_time,
		ac_time,
		namespace,
		name,
		aggValueStr,
		nil,
	)

	objId := protobuf.NamedObjectId{}
	objId.Name = &name
	objId.Namespace = &namespace

	engValProto := protobuf.Value{}
	engValProto.Type = protobuf.Value_AGGREGATE.Enum()
	engValProto.AggregateValue = &aggValProto

	param := pvalue.ParameterValue{
		Id:              &objId,
		EngValue:        &engValProto,
		RawValue:        nil,
		AcquisitionTime: &ac_time,
		GenerationTime:  &gen_time,
	}

	got := paramToString(&param)
	if got != want {
		t.Errorf("Got %s, want %s", got, want)
	}

}

func TestUnrollAggregate(t *testing.T) {

	totalVal := "10557621868"
	availableVal := "1034958734"
	var percentageUseVal float32
	percentageUseVal = 5.646844

	totalValProto := protobuf.Value{}
	totalValProto.Type = protobuf.Value_SINT64.Enum()
	totalValProto.Sint64Value = &totalVal
	availableValProto := protobuf.Value{}
	availableValProto.Type = protobuf.Value_SINT64.Enum()
	availableValProto.Sint64Value = &availableVal

	percentageUseValProto := protobuf.Value{}
	percentageUseValProto.Type = protobuf.Value_FLOAT.Enum()
	percentageUseValProto.FloatValue = &percentageUseVal

	lvl2aggValProto := protobuf.AggregateValue{}
	lvl2aggValProto.Name = []string{"percentageUse"}
	lvl2aggValProto.Value = []*protobuf.Value{&percentageUseValProto}

	lvl2aggVal := protobuf.Value{}
	lvl2aggVal.Type = protobuf.Value_AGGREGATE.Enum()
	lvl2aggVal.AggregateValue = &lvl2aggValProto

	aggValProto := protobuf.AggregateValue{}
	aggValProto.Name = []string{"total", "available", "lvl2agg"}
	aggValProto.Value = []*protobuf.Value{
		&totalValProto,
		&availableValProto,
		&lvl2aggVal,
	}

	engValProto := protobuf.Value{}
	engValProto.Type = protobuf.Value_AGGREGATE.Enum()
	engValProto.AggregateValue = &aggValProto

	expectedPaths := []string{"sdc.total", "sdc.available", "sdc.lvl2agg.percentageUse"}
	expectedVals := []*protobuf.Value{
		&totalValProto,
		&availableValProto,
		&percentageUseValProto,
	}

	var actualPaths []string
	var actualVals []*protobuf.Value

	unrollAggregate(&engValProto, "sdc", &actualPaths, &actualVals)

	for i, val := range expectedVals {
		if actualPaths[i] != expectedPaths[i] {
			t.Errorf("Want %s, got %s", expectedPaths[i], actualPaths[i])
		}
		if actualVals[i] != val {
			t.Errorf("Want %s, got %s", val, actualVals[i])
		}
	}

}

func TestEventToCsv(t *testing.T) {

    user := "User"
    genTime0 := "2023-12-03T23:47:09Z"
    recTime0 := "2023-12-03T23:47:10Z"
    genTime1 := "2023-12-03T23:48:09Z"
    recTime1 := "2023-12-03T23:48:10Z"
    message0 := "test message 0"
    message1 := "test message 1"
    var seqNum0 int32
    var seqNum1 int32
    seqNum0 = 0
    seqNum1 = 0
    severity := events.Event_WARNING
    createdBy := "guest"

    eventFilename := "events.csv"
    initEventsCsvFile(eventFilename)

    event0 := events.Event {
        Source: &user,
        GenerationTime: &genTime0,
        ReceptionTime: &recTime0,
        SeqNumber: &seqNum0,
        Message: &message0,
        Severity: &severity,
        CreatedBy: &createdBy,    
    }
    event1 := events.Event {
        Source: &user,
        GenerationTime: &genTime1,
        ReceptionTime: &recTime1,
        SeqNumber: &seqNum1,
        Message: &message1,
        Severity: &severity,
        CreatedBy: &createdBy,
    }

    eventList := []*events.Event{ &event0, &event1 }

    for _, e := range eventList {
        eventToCsv(e, &eventFilename)
    }


    expectedFile := "events.csv"
    f, err := os.Open(expectedFile)
    if os.IsNotExist(err) {
        t.Errorf("Expected file %s to exist\n%q", expectedFile, err)
    } else if err != nil {
        t.Errorf("Error opening file: %q\n", err)
    } else {
        defer f.Close()

        scanner := bufio.NewScanner(f)
        
        i := 0 // ignore the header
        for scanner.Scan() {
            line := scanner.Text()
            if i > 0 {
                expected := strings.TrimSpace(eventToString(eventList[i-1]))
                if line != expected {   
                    t.Errorf("got:\n%s\nwant:\n%s\n", line, eventToString(eventList[i-1]))
                }
            }
            i++
        }

        if err = scanner.Err(); err != nil {
            t.Errorf("Error reading the file: %q\n", err)
        }

	}

}
