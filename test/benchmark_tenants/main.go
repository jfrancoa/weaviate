package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/weaviate/weaviate/entities/models"
)

const (
	schemaName   = "Benchmark"
	weaviatePort = "8383"
)

// creates a dummy class with multi-tenancy enabled
func createTenantRequest(url, className, tenantName string) *http.Request {
	tenantObj := models.Tenant{
		Name: tenantName,
	}
	request := createRequest(url+"schema/"+className+"/tenants", "POST", []models.Tenant{tenantObj})
	return request
}

// deletes the created dummy class
func deleteMultiTenantClassRequest(url, className string) *http.Request {
	request := createRequest(url+"schema/"+className, "DELETE", nil)
	return request
}

// creates a dummy class with multi-tenancy enabled
func createMultiTenantClassRequest(url, className string) *http.Request {
	classObj := &models.Class{
		Class:       className,
		Description: "Dummy class for benchmarking purposes",
		Properties: []*models.Property{
			{
				DataType:    []string{"int"},
				Description: "The value of the counter in the dataset",
				Name:        "counter",
			},
		},
		MultiTenancyConfig: &models.MultiTenancyConfig{
			Enabled: true,
		},
		Vectorizer: "none",
	}
	request := createRequest(url+"schema", "POST", classObj)
	return request
}

// createRequest creates requests
func createRequest(url string, method string, payload interface{}) *http.Request {
	var body io.Reader = nil
	if payload != nil {
		jsonBody, err := json.Marshal(payload)
		if err != nil {
			panic(errors.Wrap(err, "Could not marshal request"))
		}
		body = bytes.NewBuffer(jsonBody)
	}
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(errors.Wrap(err, "Could not create request"))
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	return request
}

// performRequest runs requests
func performRequest(c *http.Client, request *http.Request) (int, []byte, int64, error) {
	timeStart := time.Now()
	response, err := c.Do(request)
	requestTime := time.Since(timeStart).Milliseconds()

	if err != nil {
		return 0, nil, requestTime, err
	}

	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return 0, nil, requestTime, err
	}

	return response.StatusCode, body, requestTime, nil
}

// check if weaviate is running before starting the benchmark
func isWeaviateRunning(c *http.Client, url string) bool {
	requestReady := createRequest(url+".well-known/ready", "GET", nil)

	responseStartedCode, _, _, err := performRequest(c, requestReady)
	return err == nil && responseStartedCode == 200

}

// calculates the percentile 99 for all the response times stored in the slice
func calculate99thPercentile(responseTimes []int64) int64 {
	// Sort the response times in ascending order
	sort.Slice(responseTimes, func(i, j int) bool {
		return responseTimes[i] < responseTimes[j]
	})

	// Calculate the index for the 99th percentile
	index := int(float64(len(responseTimes)) * 0.99)

	// Retrieve the value at the calculated index
	percentile99 := responseTimes[index]

	return percentile99
}

// stores the response times to a json file
func saveTimestampsToFile(filename string, timestamps []int64) error {
	// Serialize the timestamps to JSON
	data, err := json.Marshal(timestamps)
	if err != nil {
		return err
	}

	// Write the serialized data to the file
	if err = os.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	return nil
}

// loads the response times from a json file
func readTimestampsFromFile(filename string) ([]int64, error) {
	// Read the serialized data from the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Deserialize the data into a slice of int64
	var timestamps []int64
	err = json.Unmarshal(data, &timestamps)
	if err != nil {
		return nil, err
	}

	return timestamps, nil
}

func main() {

	var numberTenants int
	var totalElapsedTime int64
	var tenantRequestsTimeStore []int64

	flag.IntVar(&numberTenants, "number", 1000, "Number of tenants to be created on the schema.")
	flag.Parse()

	fmt.Printf("Benchmarking the creation of %v tenants\n", numberTenants)

	t := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 120 * time.Second,
		}).DialContext,
		MaxIdleConnsPerHost:   100,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	c := &http.Client{Transport: t}
	url := "http://localhost:" + weaviatePort + "/v1/"

	if !isWeaviateRunning(c, url) {
		panic(fmt.Sprintf("Weaviate isn't running on endpoint %v", url))
	}
	fmt.Print("Weaviate is fresh!: up and running.\n")

	tenantRequestsTimeStore = []int64{}
	totalElapsedTime = 0

	defer performRequest(c, deleteMultiTenantClassRequest(url, schemaName))
	requestSchema := createMultiTenantClassRequest(url, schemaName)

	// Add schema
	responseSchemaCode, body, _, err := performRequest(c, requestSchema)

	if err != nil {
		panic(errors.Wrap(err, "Could not add schema, error: "))
	} else if responseSchemaCode != 200 {
		panic(fmt.Sprintf("Could not add schema, http error code: %v, body: %v", responseSchemaCode, string(body)))
	}

	for i := 0; i < numberTenants; i++ {
		requestTenant := createTenantRequest(url, schemaName, "tenant"+strconv.Itoa(i))
		responseCreateTenant, body, timeTenantAdd, err := performRequest(c, requestTenant)
		if err != nil {
			panic(errors.Wrap(err, "Could not add tenant to schema, error: "))
		} else if responseCreateTenant != 200 {
			panic(fmt.Sprintf("Could not add tenant, http error code: %v, body: %v", responseCreateTenant, string(body)))
		}
		totalElapsedTime += timeTenantAdd
		tenantRequestsTimeStore = append(tenantRequestsTimeStore, timeTenantAdd)
	}

	p99ElapsedTime := calculate99thPercentile(tenantRequestsTimeStore)
	fmt.Printf("Percentile 99 response time per tenant creation: %v ms\n", p99ElapsedTime)
	fmt.Printf("Total elapsed time to create %v tenants: %v ms\n", numberTenants, totalElapsedTime)

}
