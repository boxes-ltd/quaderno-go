package quaderno

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTransactionCreateCustomerId_MarshalJSON(t *testing.T) {
	id := TransactionCreateCustomerId("cust_123")
	b, err := json.Marshal(id)
	if err != nil {
		t.Fatalf("MarshalJSON error: %v", err)
	}
	want := `{"id":"cust_123"}`
	if string(b) != want {
		t.Errorf("MarshalJSON = %q, want %q", string(b), want)
	}
}

func TestTransactionCreate_nilParams(t *testing.T) {
	c := NewClient("key", "https://api.example.com")
	_, err := c.Transactions.Create(context.Background(), nil)
	if err == nil {
		t.Error("expected error for nil params")
	}
}

func TestTransactionCreate_success(t *testing.T) {
	id := int64(42)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/transactions" {
			t.Errorf("path = %q, want /transactions", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(TransactionCreateResponse{Id: &id})
	}))
	defer srv.Close()

	txType := TransactionTypeSale
	c := NewClient("key", srv.URL)
	resp, err := c.Transactions.Create(context.Background(), &TransactionCreateParams{
		Type: &txType,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Id == nil || *resp.Id != 42 {
		t.Errorf("Id = %v, want 42", resp.Id)
	}
}

func TestTransactionCreate_encodesCustomerIdAsObject(t *testing.T) {
	var body map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&body)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(TransactionCreateResponse{})
	}))
	defer srv.Close()

	c := NewClient("key", srv.URL)
	_, err := c.Transactions.Create(context.Background(), &TransactionCreateParams{
		Customer: TransactionCreateCustomerId("cust_abc"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	customer, ok := body["customer"].(map[string]any)
	if !ok {
		t.Fatalf("customer field = %T, want map[string]any", body["customer"])
	}
	if customer["id"] != "cust_abc" {
		t.Errorf("customer.id = %v, want %q", customer["id"], "cust_abc")
	}
}

func TestTransactionCreate_propagatesApiError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"error":"invalid transaction"}`))
	}))
	defer srv.Close()

	c := NewClient("key", srv.URL)
	_, err := c.Transactions.Create(context.Background(), &TransactionCreateParams{})
	if _, ok := err.(*ApiError); !ok {
		t.Errorf("expected *ApiError, got %T: %v", err, err)
	}
}