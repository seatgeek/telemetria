package telemetria

import (
	"reflect"
	"testing"

	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	influxdb "github.com/influxdata/influxdb1-client/v2"
)

var client = createClient()

// The purpose of this function is to prevent compilation
// if we forgot to correctly implement the Recorder interface
// for any of the exposed structs in the library
func failCompilationIfNotRecorder(r Recorder) {
	// Nothing to do here, the compiler should do
	// the work by itself
}

func createClient() influxdb.Client {
	client, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr: "http://localhost:8086",
	})

	if err != nil {
		panic(err)
	}

	return client
}

func queryDB(client influxdb.Client, cmd string) (res []influxdb.Result, err error) {
	q := influxdb.Query{
		Command:  cmd,
		Database: "test",
	}

	if response, err := client.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}

	return res, nil
}

func TestIsRecorder(t *testing.T) {
	failCompilationIfNotRecorder(SimpleRecorder{})
	failCompilationIfNotRecorder(NoRecorder{})
}

func TestNewUDPRecorder(t *testing.T) {
	recorder, err := NewRecorder("udp://localhost:8089")

	if err != nil {
		t.Errorf("Error in TestNewUDPRecorder:\n %s", err.Error())
		return
	}

	failCompilationIfNotRecorder(recorder)
}

func TestNewUDPRecorderFail(t *testing.T) {
	_, err := NewRecorder("localhost")

	if err == nil {
		t.Errorf("This should fail because of missing port")
	}
}

func TestNewHttpRecorder(t *testing.T) {
	recorder, err := NewRecorder("http://localhost")

	if err != nil {
		t.Errorf("Error in TestNewHttpRecorder:\n %s", err.Error())
		return
	}

	failCompilationIfNotRecorder(recorder)
}

func TestErrorInNewRecorder(t *testing.T) {
	_, err := NewRecorder("")

	if err == nil {
		t.Errorf("I was expecting an empty string to not be valid")
	}
}

func TestWriteOne(t *testing.T) {
	queryDB(client, "DROP DATABASE test")
	queryDB(client, "CREATE DATABASE test")
	recorder, err := NewRecorder("http://localhost:8086/test")

	if err != nil {
		t.Errorf("Error in TestWriteOne:\n %s", err.Error())
		return
	}

	err = recorder.WriteOne(Metric{
		Name:   "things",
		Tags:   map[string]string{"one": "two", "three": "four"},
		Fields: map[string]interface{}{"field_one": "1", "field_2": "2"},
	})

	if err != nil {
		t.Errorf("Error inserting metric:\n %s", err.Error())
		return
	}

	results, err := queryDB(client, "SELECT * from things")

	if err != nil {
		t.Fatalf("Could not query the databse:\n %s", err.Error())
	}

	if len(results) == 0 {
		t.Errorf("No metrics where inserted")
		return
	}

	result := results[0]

	if result.Series[0].Name != "things" {
		t.Errorf("Invalid name for metric found '%s' but expected '%s'", result.Series[0].Name, "things")
		return
	}

	columns := result.Series[0].Columns
	expected := []string{"time", "field_2", "field_one", "one", "three"}

	if reflect.DeepEqual(columns, expected) != true {
		t.Errorf("Wrong columns: found '%s' but expected '%s'", columns, expected)
		return
	}

	values := result.Series[0].Values[0][1:]
	expectedV := []interface{}{"2", "1", "two", "four"}

	if reflect.DeepEqual(values, expectedV) != true {
		t.Errorf("Wrong values: found '%s' but expected '%s'", values, expectedV)
		return
	}
}
