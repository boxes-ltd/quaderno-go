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
	// Optional: The seller's country. 2-letter ISO country code. Defaults to the account's country.
	FromCountry *string
	// Optional: The seller's ZIP or postal code. Defaults to the account's postal code.
	FromPostalCode *string
	// Required: The customer's country. 2-letter ISO country code.
	ToCountry *string
	// Optional: The customer's ZIP or postal code.
	ToPostalCode *string
	// Optional: The customer's city. Recommended for US Sales Tax calculations.
	ToCity *string
	// Optional: The customer's street address. Recommended for US Sales Tax calculations.
	ToStreet *string
	// Optional: The customer's tax ID. Quaderno can validate VAT/GST numbers from the EU, United Kingdom, Switzerland,
	// Québec (Canada), Australia, and New Zealand.
	TaxID *string
	// Optional: The transaction's tax code. Tax codes can be obtained via GET /tax_codes. Defaults to the account's
	// default tax code.
	TaxCode *TaxCode
	// Optional: Specifies whether the price is considered inclusive of taxes or exclusive of taxes.
	TaxBehavior *TaxBehavior
	// Optional: Specifies whether the product is a good or a service. Defaults to the account's default.
	ProductType *ProductType
	// Optional: The transaction's date. Defaults to today.
	Date *string
	// Optional: The transaction's amount.
	Amount *float64
	// Optional: The transaction's currency. Three-letter ISO currency code, in uppercase. Defaults to the account's
	// currency.
	Currency *string
}

type TaxCalculateResponse struct {
	// Name of the tax.
	Name *string `json:"name"`
	// Tax rate applied.
	Rate *float64 `json:"rate"`
	// Percentage of the subtotal used for calculating the tax amount. It's usually 100% but there are a few exceptions.
	// For example, in Texas the taxable base is 80% for SaaS products.
	TaxablePart *float64 `json:"taxable_part"`
	// Country used for the tax calculation.
	Country *string `json:"country"`
	// Currency used for the tax calculation.
	Currency *string `json:"currency"`
	// Region used for the tax calculation.
	Region *string `json:"region"`
	// Tax county. Only for US sales tax.
	County *string `json:"county"`
	// City used for the tax calculation.
	City *string `json:"city"`
	// The transaction's tax code used for calculating the tax rate.
	TaxCode *TaxCode `json:"tax_code"`
	// Whether the price was considered inclusive of taxes or exclusive of taxes.
	TaxBehavior *TaxBehavior `json:"tax_behavior"`
	// Whether the product is a good or a service.
	ProductType *ProductType `json:"product_type"`
	// The tax amount.
	TaxAmount *float64 `json:"tax_amount"`
	// Price before taxes.
	Subtotal *float64 `json:"subtotal"`
	// Total price, including taxes.
	TotalAmount *float64 `json:"total_amount"`
	// Tax calculation status.
	Status *TaxStatus `json:"status"`
	// Help message complementing the tax calculation status.
	Notice *string `json:"notice"`
	// Only for jurisdictions that need to apply two tax rates (e.g. some Canadian provinces).
	AdditionalName *string `json:"additional_name"`
	// Only for jurisdictions that need to apply two tax rates (e.g. some Canadian provinces).
	AdditionalRate *float64 `json:"additional_rate"`
	// Only for jurisdictions that need to apply two tax rates (e.g. some Canadian provinces).
	AdditionalTaxablePart *float64 `json:"additional_taxable_part"`
	// Only for jurisdictions that need to apply two tax rates (e.g. some Canadian provinces).
	AdditionalTaxAmount *float64 `json:"additional_tax_amount"`
}

func (s *Taxes) Calculate(ctx context.Context, params *TaxCalculateParams) (*TaxCalculateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params is required")
	}
	if params.ToCountry == nil || *params.ToCountry == "" {
		return nil, fmt.Errorf("to_country is required")
	}

	q := url.Values{}
	q.Set("to_country", *params.ToCountry)

	if params.FromCountry != nil {
		q.Set("from_country", *params.FromCountry)
	}
	if params.FromPostalCode != nil {
		q.Set("from_postal_code", *params.FromPostalCode)
	}
	if params.ToPostalCode != nil {
		q.Set("to_postal_code", *params.ToPostalCode)
	}
	if params.ToCity != nil {
		q.Set("to_city", *params.ToCity)
	}
	if params.ToStreet != nil {
		q.Set("to_street", *params.ToStreet)
	}
	if params.TaxID != nil {
		q.Set("tax_id", *params.TaxID)
	}
	if params.TaxCode != nil {
		q.Set("tax_code", string(*params.TaxCode))
	}
	if params.TaxBehavior != nil {
		q.Set("tax_behavior", string(*params.TaxBehavior))
	}
	if params.ProductType != nil {
		q.Set("product_type", string(*params.ProductType))
	}
	if params.Date != nil {
		q.Set("date", *params.Date)
	}
	if params.Amount != nil {
		q.Set("amount", strconv.FormatFloat(*params.Amount, 'f', -1, 64))
	}
	if params.Currency != nil {
		q.Set("currency", *params.Currency)
	}

	var resp TaxCalculateResponse

	err := s.client.doRequest(ctx, http.MethodGet, "/tax_rates/calculate", q, nil, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
