package quaderno

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestTaxCalculate_nilParams(t *testing.T) {
	c := NewClient("key", "https://api.example.com")
	_, err := c.Taxes.Calculate(context.Background(), nil)
	if err == nil {
		t.Error("expected error for nil params")
	}
}

func TestTaxCalculate_missingToCountry(t *testing.T) {
	c := NewClient("key", "https://api.example.com")
	_, err := c.Taxes.Calculate(context.Background(), &TaxCalculateParams{})
	if err == nil {
		t.Error("expected error when ToCountry is nil")
	}
}

func TestTaxCalculate_emptyToCountry(t *testing.T) {
	c := NewClient("key", "https://api.example.com")
	_, err := c.Taxes.Calculate(context.Background(), &TaxCalculateParams{ToCountry: new("")})
	if err == nil {
		t.Error("expected error when ToCountry is empty")
	}
}

func TestTaxCalculate_sendsToCountry(t *testing.T) {
	var capturedQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(TaxCalculateResponse{})
	}))
	defer srv.Close()

	c := NewClient("key", srv.URL)
	_, err := c.Taxes.Calculate(context.Background(), &TaxCalculateParams{ToCountry: new("DE")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := capturedQuery.Get("to_country"); got != "DE" {
		t.Errorf("to_country = %q, want %q", got, "DE")
	}
}

func TestTaxCalculate_sendsAllOptionalParams(t *testing.T) {
	var capturedQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(TaxCalculateResponse{})
	}))
	defer srv.Close()

	params := &TaxCalculateParams{
		ToCountry:      new("DE"),
		FromCountry:    new("US"),
		FromPostalCode: new("10001"),
		ToPostalCode:   new("10115"),
		ToCity:         new("Berlin"),
		ToStreet:       new("Unter den Linden 1"),
		TaxID:          new("DE123456789"),
		TaxCode:        new(TaxCodeSaaS),
		TaxBehavior:    new(TaxBehaviorExclusive),
		ProductType:    new(ProductTypeService),
		Date:           new("2024-01-15"),
		Amount:         new(100.5),
		Currency:       new("EUR"),
	}

	c := NewClient("key", srv.URL)
	_, err := c.Taxes.Calculate(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cases := map[string]string{
		"to_country":       "DE",
		"from_country":     "US",
		"from_postal_code": "10001",
		"to_postal_code":   "10115",
		"to_city":          "Berlin",
		"to_street":        "Unter den Linden 1",
		"tax_id":           "DE123456789",
		"tax_code":         string(TaxCodeSaaS),
		"tax_behavior":     string(TaxBehaviorExclusive),
		"product_type":     string(ProductTypeService),
		"date":             "2024-01-15",
		"amount":           "100.5",
		"currency":         "EUR",
	}
	for param, want := range cases {
		if got := capturedQuery.Get(param); got != want {
			t.Errorf("%s = %q, want %q", param, got, want)
		}
	}
}

func TestTaxCalculate_parsesResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(TaxCalculateResponse{Rate: new(20.0), Name: new("VAT")})
	}))
	defer srv.Close()

	c := NewClient("key", srv.URL)
	resp, err := c.Taxes.Calculate(context.Background(), &TaxCalculateParams{ToCountry: new("DE")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Rate == nil || *resp.Rate != 20.0 {
		t.Errorf("Rate = %v, want 20.0", resp.Rate)
	}
	if resp.Name == nil || *resp.Name != "VAT" {
		t.Errorf("Name = %v, want VAT", resp.Name)
	}
}

func TestTaxCalculate_propagatesApiError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
	}))
	defer srv.Close()

	c := NewClient("key", srv.URL)
	_, err := c.Taxes.Calculate(context.Background(), &TaxCalculateParams{ToCountry: new("DE")})
	if _, ok := errors.AsType[*ApiError](err); !ok {
		t.Errorf("expected *ApiError, got %T: %v", err, err)
	}
}