package service

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"food_delivery/model"
	"io"
	"net/http"
	"time"
)

func fetchSuppliers(page, limit int) ([]model.Supplier, error) {
	// url := "https://foodapi.golang.nixdev.co/suppliers?limit=10&page=1"
	url := fmt.Sprintf("http://foodapi.golang.nixdev.co/suppliers?limit=%d&page=%d", limit, page)

	// Create a custom HTTP client that skips SSL verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching suppliers: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var suppliers []model.Supplier
	var response model.Response

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}
	suppliers = response.Suppliers
	// fmt.Printf("Suppliers: %+v\n", suppliers)

	return suppliers, nil
}

func FetchAllSuppliers() ([]model.Supplier, error) {
	var allSuppliers []model.Supplier
	limit := 10
	page := 1

	for {
		// Fetch suppliers for the current page
		suppliers, err := fetchSuppliers(page, limit)
		if err != nil {
			return nil, err
		}

		// If no more suppliers are returned, stop the loop
		if len(suppliers) == 0 {
			break
		}

		// Append the fetched suppliers to our list
		allSuppliers = append(allSuppliers, suppliers...)

		// Move to the next page
		page++
	}

	return allSuppliers, nil
}
