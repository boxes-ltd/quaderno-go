package quaderno

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Taxes service

type TaxCalculateParams struct {
	FromCountry    string      // Optional: The seller's country. 2-letter ISO country code. Defaults to the account's country.
	FromPostalCode string      // Optional: The seller's ZIP or postal code. Defaults to the account's postal code.
	ToCountry      string      // Required: The customer's country. 2-letter ISO country code.
	ToPostalCode   string      // Optional: The customer's ZIP or postal code.
	ToCity         string      // Optional: The customer's city. Recommended for US Sales Tax calculations.
	ToStreet       string      // Optional: The customer's street address. Recommended for US Sales Tax calculations.
	TaxID          string      // Optional: The customer's tax ID. Quaderno can validate VAT/GST numbers from the EU, United Kingdom, Switzerland, Québec (Canada), Australia, and New Zealand.
	TaxCode        TaxCode     // Optional: The transaction's tax code. Tax codes can be obtained via GET /tax_codes. Defaults to the account's default tax code.
	TaxBehavior    TaxBehavior // Optional: Specifies whether the price is considered inclusive of taxes or exclusive of taxes.
	ProductType    ProductType // Optional: Specifies whether the product is a good or a service. Defaults to the account's default.
	Date           string      // Optional: The transaction's date. Defaults to today.
	Amount         *float64    // Optional: The transaction's amount.
	Currency       string      // Optional: The transaction's currency. Three-letter ISO currency code, in uppercase. Defaults to the account's currency.
}

type TaxCalculateResponse struct {
	Name                  string      `json:"name"`                    // Name of the tax.
	Rate                  float64     `json:"rate"`                    // Tax rate applied.
	TaxablePart           float64     `json:"taxable_part"`            // Percentage of the subtotal used for calculating the tax amount. It's usually 100% but there are a few exceptions. For example, in Texas the taxable base is 80% for SaaS products.
	Country               string      `json:"country"`                 // Country used for the tax calculation.
	Region                string      `json:"region"`                  // Region used for the tax calculation.
	County                string      `json:"county"`                  // Tax county. Only for US sales tax.
	City                  string      `json:"city"`                    // City used for the tax calculation.
	TaxCode               string      `json:"tax_code"`                // The transaction's tax code used for calculating the tax rate.
	TaxBehavior           TaxBehavior `json:"tax_behavior"`            // Whether the price was considered inclusive of taxes or exclusive of taxes.
	TaxAmount             float64     `json:"tax_amount"`              // The tax amount.
	Subtotal              float64     `json:"subtotal"`                // Price before taxes.
	TotalAmount           float64     `json:"total_amount"`            // Total price, including taxes.
	Status                TaxStatus   `json:"status"`                  // Tax calculation status.
	Notice                string      `json:"notice"`                  // Help message complementing the tax calculation status.
	AdditionalName        string      `json:"additional_name"`         // Only for jurisdictions that need to apply two tax rates (e.g. some Canadian provinces).
	AdditionalRate        float64     `json:"additional_rate"`         // Only for jurisdictions that need to apply two tax rates (e.g. some Canadian provinces).
	AdditionalTaxablePart float64     `json:"additional_taxable_part"` // Only for jurisdictions that need to apply two tax rates (e.g. some Canadian provinces).
	AdditionalTaxAmount   float64     `json:"additional_tax_amount"`   // Only for jurisdictions that need to apply two tax rates (e.g. some Canadian provinces).
}

func (s *Taxes) Calculate(ctx context.Context, params *TaxCalculateParams) (*TaxCalculateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params is required")
	}
	if params.ToCountry == "" {
		return nil, fmt.Errorf("to_country is required")
	}

	q := url.Values{}
	q.Set("to_country", params.ToCountry)

	if params.FromCountry != "" {
		q.Set("from_country", params.FromCountry)
	}
	if params.FromPostalCode != "" {
		q.Set("from_postal_code", params.FromPostalCode)
	}
	if params.ToPostalCode != "" {
		q.Set("to_postal_code", params.ToPostalCode)
	}
	if params.ToCity != "" {
		q.Set("to_city", params.ToCity)
	}
	if params.ToStreet != "" {
		q.Set("to_street", params.ToStreet)
	}
	if params.TaxID != "" {
		q.Set("tax_id", params.TaxID)
	}
	if params.TaxCode != "" {
		q.Set("tax_code", string(params.TaxCode))
	}
	if params.TaxBehavior != "" {
		q.Set("tax_behavior", string(params.TaxBehavior))
	}
	if params.ProductType != "" {
		q.Set("product_type", string(params.ProductType))
	}
	if params.Date != "" {
		q.Set("date", params.Date)
	}
	if params.Amount != nil {
		q.Set("amount", strconv.FormatFloat(*params.Amount, 'f', -1, 64))
	}
	if params.Currency != "" {
		q.Set("currency", params.Currency)
	}

	var resp TaxCalculateResponse

	err := s.client.doRequest(ctx, http.MethodGet, "/tax_rates/calculate", q, nil, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
