package quaderno

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Transactions service

type TransactionCreateParams struct {
	// The transaction’s type.
	Type *TransactionType `json:"type,omitempty"`
	// Three-letter ISO currency code, in uppercase. Defaults to the account’s default currency.
	Currency *string `json:"currency,omitempty"`
	// Set of key-value pairs that you can attach to an object. This can be useful for storing additional information
	// about the object in a structured format. You can have up to 20 keys, with key names up to 40 characters long and
	// values up to 500 characters long.
	CustomMetadata map[string]any `json:"custom_metadata,omitempty"`
	// The data of the customer who paid the transaction. You can reference an existing contact by its id, or pass the
	// json object to create a new contact.
	Customer TransactionCreateCustomerParams `json:"customer,omitempty"`
	// The transaction’s date. Defaults to today.
	Date *string `json:"date,omitempty"`
	// The list of individual items that make up the transaction.
	Items []*TransactionCreateItemParams `json:"items,omitempty"`
	// The mailing address to where the order will be shipped. Use it if the order contains physical goods.
	ShippingAddress *TransactionCreateShippingAddressParams `json:"shipping_address,omitempty"`
	// Evidence of the customer’s location. **Highly recommended**.
	Evidence *TransactionCreateEvidenceParams `json:"evidence,omitempty"`
	// Detailed information about the transaction payment.
	Payment *TransactionCreatePaymentParams `json:"payment,omitempty"`
	// The name of the platform that processed the transaction. E.g. shopify, woocommerce or any user agent you may want
	// to use to identify yourself… **Recommended**
	Processor *string `json:"processor,omitempty"`
	// The ID of the transaction in the processor. Use the same ID to link sales and refunds for the same operation.
	// **Recommended**
	ProcessorId *string `json:"processor_id,omitempty"`
	// Processor total fee, in cents.
	ProcessorFeeCents *int64 `json:"processor_fee_cents,omitempty"`
	// The rate used for currency conversion when the transaction currency differs from the account base currency. If
	// omitted, the default ECB reference exchange rate for the transaction date is applied.
	ExchangeRate *float64 `json:"exchange_rate,omitempty"`
	// The number of the related order. **Recommended**.
	PoNumber *string `json:"po_number,omitempty"`
	// Optional notes attached to the transaction.
	Notes *string `json:"notes,omitempty"`
	// Tags attached to the transaction, formatted as a string of comma-separated values. Tags are additional short
	// descriptors, commonly used for filtering and searching. Each individual tag is limited to 40 characters in
	// length.
	Tags *string `json:"tags,omitempty"`
}

type TransactionCreateCustomerParams interface {
	isTransactionCreateCustomerParams()
}

type TransactionCreateCustomer struct {
	// City/District/Suburb/Town/Village.
	City *string `json:"city,omitempty"`
	// If the contact is a company, this is its contact person.
	ContactPerson *string `json:"contact_person,omitempty"`
	// 2-letter ISO country code.
	Country *string `json:"country,omitempty"`
	// If the contact is a company, this is the department.
	Department *string `json:"department,omitempty"`
	// Default discount for this contact.
	DefaultDiscount *float64 `json:"default_discount,omitempty"`
	// The contact's email address
	Email *string `json:"email,omitempty"`
	// The contact's first name. Only if the contact is a person.
	FirstName *string `json:"first_name,omitempty"`
	// The type of contact.
	Kind *CustomerKind `json:"kind,omitempty"`
	// The contact's preferred language. 2-letter ISO language code. Should be included in the account's translations
	// list.
	Language *string `json:"language,omitempty"`
	// The contact's last name. Only if the contact is a person.
	LastName *string `json:"last_name,omitempty"`
	// Internal notes about the contact.
	Notes *string `json:"notes,omitempty"`
	// The contact's phone number.
	Phone *string `json:"phone_1,omitempty"`
	// ZIP or postal code.
	PostalCode *string `json:"postal_code,omitempty"`
	// The external platform where the contact was imported from, if applicable.
	Processor *string `json:"processor,omitempty"`
	// The ID the payment_processor assigned to the contact.
	ProcessorId *string `json:"processor_id,omitempty"`
	// State/Province/Region.
	Region *string `json:"region,omitempty"`
	// Address line 1 (Street address/PO Box).
	StreetLine1 *string `json:"street_line_1,omitempty"`
	// Address line 2 (Apartment/Suite/Unit/Building).
	StreetLine2 *string `json:"street_line_2,omitempty"`
	// The contact's tax identification number. Quaderno can validate EU VAT numbers, ABN, and NZBN.
	TaxId *string `json:"tax_id,omitempty"`
	// Specifies the tax status of the contact.
	TaxStatus *TaxStatus `json:"tax_status,omitempty"`
	// The contact's website
	Web *string `json:"web,omitempty"`
}

func (TransactionCreateCustomer) isTransactionCreateCustomerParams() {}

type TransactionCreateCustomerId string

func (TransactionCreateCustomerId) isTransactionCreateCustomerParams() {}

func (id TransactionCreateCustomerId) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"id": string(id)})
}

type TransactionCreateItemParams struct {
	// Code of an existing product. If present, you don't need to set the description.
	ProductCode *string `json:"product_code,omitempty"`
	// Concept being sold or refunded. Meant to be displayable to the customer.
	Description *string `json:"description,omitempty"`
	// Quantity of units for the transaction item. Defaults to 1.
	Quantity *int64 `json:"quantity,omitempty"`
	// This represents the discount percent out of 100 included in the amount, if applicable.
	DiscountRate *float64 `json:"discount_rate,omitempty"`
	// Total amount of the transaction item. This should always be equal to the amount charged after discounts and
	// taxes.
	Amount *float64 `json:"amount,omitempty"`
	// The tax details applied to the transaction item. Responses from the tax calculation endpoint are valid input.
	Tax *TransactionCreateTaxParams `json:"tax,omitempty"`
}

type TransactionCreateTaxParams struct {
	// Name of the tax.
	Name *string `json:"name,omitempty"`
	// Tax rate applied.
	Rate *float64 `json:"rate,omitempty"`
	// Percentage of the subtotal used for calculating the tax amount. It's usually 100% but there are a few exceptions.
	// For example, in Texas the taxable base is 80% for SaaS products.
	TaxablePart *float64 `json:"taxable_part,omitempty"`
	// Country used for the tax calculation.
	Country *string `json:"country,omitempty"`
	// Region used for the tax calculation.
	Region *string `json:"region,omitempty"`
	// The transaction's tax code used for calculating the tax rate. When not specified, it uses the account's default
	// tax code.
	TaxCode *TaxCode `json:"tax_code,omitempty"`
	// Only for jurisdictions that need to apply two tax rates (e.g. some Canadian provinces).
	AdditionalName *string `json:"additional_name,omitempty"`
	// Only for jurisdictions that need to apply two tax rates (e.g. some Canadian provinces).
	AdditionalRate *float64 `json:"additional_rate,omitempty"`
	// Only for jurisdictions that need to apply two tax rates (e.g. some Canadian provinces).
	AdditionalTaxablePart *float64 `json:"additional_taxable_part,omitempty"`
	// Whether the product is considered an import sale or not.
	Import *bool `json:"import,omitempty"`
}

type TransactionCreateShippingAddressParams struct {
	// Address line 1 (Street address/PO Box).
	StreetLine1 *string `json:"street_line_1,omitempty"`
	// Address line 2 (Apartment/Suite/Unit/Building).
	StreetLine2 *string `json:"street_line_2,omitempty"`
	// City/District/Suburb/Town/Village.
	City *string `json:"city,omitempty"`
	// 2-letter ISO country code.
	Country *string `json:"country,omitempty"`
	// State/Province/Region.
	Region *string `json:"region,omitempty"`
	// ZIP or postal code.
	PostalCode *string `json:"postal_code,omitempty"`
}

type TransactionCreateEvidenceParams struct {
	// Additional evidence used to proof the customer's location. Required if there's a `additional_evidence_country'.
	AdditionalEvidence *string `json:"additional_evidence,omitempty"`
	// 2-letter ISO 3166-1 code.
	AdditionalEvidenceCountry *string `json:"additional_evidence_country,omitempty"`
	// 2-letter ISO 3166-1 code.
	BankCountry *string `json:"bank_country,omitempty"`
	// 2-letter ISO 3166-1 code.
	BillingCountry *string `json:"billing_country,omitempty"`
	// Unique identifier for the related invoice.
	DocumentId *string `json:"document_id,omitempty"`
	// The customer's IP address
	IpAddress *string `json:"ip_address,omitempty"`
	// Internal notes about the evidence.
	Notes *string `json:"notes,omitempty"`
	// State of the evidence
	State *EvidenceState `json:"state,omitempty"`
}

type TransactionCreatePaymentParams struct {
	// The payment method used to pay the transaction.
	Method *PaymentMethod `json:"method,omitempty"`
	// The name of the payment processor used to take the payment. E.g. stripe, paypal…
	Processor *string `json:"processor,omitempty"`
	// The ID of the transaction in the payment processor.
	ProcessorId *string `json:"processor_id,omitempty"`
}

type TransactionCreateResponse struct {
	// Unique identifier for the transaction.
	Id *int64 `json:"id"`
	// The transaction’s type.
	Type *DocumentType `json:"type"`
	// Timestamp when the transaction was created. Measured in seconds since the Unix epoch.
	CreatedAt *int64 `json:"created_at"`
	// Unique, sequential number identifying this transaction. Transaction numbers must be sequential without gaps or
	// duplicates to comply with legal requirements. Automatically generated if not provided.
	Number *string `json:"number"`
	// Date of the transaction.
	IssueDate *string `json:"issue_date"`
	// Reference to the original invoice or receipt that this transaction is refunding or correcting.
	RelatedDocument *TransactionCreateResponseRelatedDocument `json:"related_document"`
	// Original purchase order number from the related invoice, maintained for reference.
	PoNumber *string `json:"po_number"`
	// Due date for credit application or refund processing. May differ from issue date for complex refund workflows.
	DueDate *string `json:"due_date"`
	// Three-letter ISO 4217 currency code in uppercase.
	Currency *string `json:"currency"`
	// Array of tags for categorizing and filtering credit notes. Useful for tracking refund reasons or processing
	//status.
	Tags *string `json:"tags"`
	// Optional notes attached to the transaction.
	Notes *string `json:"notes"`
	// Conteact related to this transaction. Can reference an existing contact by ID or include full contact details.
	Contact *TransactionCreateResponseContact `json:"contact"`
	// Brief summary or reason for this credit note. Displayed prominently on the document (e.g., "Refund for Order
	// #12345").
	Subject *string `json:"subject"`
	// Delivery history showing when and how this credit note was sent to the customer (email, etc.).
	Deliveries []*TransactionCreateResponseDelivery `json:"deliveries"`
	// Line items included in this transaction, typically matching items from the original invoice or receipt. Maximum
	// 200 items per request; use update requests to add more. Total limit: 1000 items per transaction.
	Items []*TransactionCreateResponseItem `json:"items"`
	// Calculated tax breakdown showing tax being applied or credited.
	Taxes []*TransactionCreateResponseTax `json:"taxes"`
	// Payment records associated with the transaction.
	Payments []*TransactionCreateResponsePayment `json:"payments"`
	// The name of the platform that processed the transaction.
	Processor *string `json:"processor"`
	// The ID of the transaction in the processor.
	ProcessorId *string `json:"processor_id"`
	// Processor total fee, in cents.
	ProcessorFeeCents *int64 `json:"processor_fee_cents"`
	// Custom key-value pairs for storing additional structured data.
	CustomMetadata map[string]any `json:"custom_metadata"`
	// Exchange rate applied when converting from credit currency to account base currency.
	ExchangeRate *float64 `json:"exchange_rate"`
	// Sum of all line items before discounts and taxes are applied, in cents.
	SubtotalCents *int64 `json:"subtotal_cents"`
	// Total discount amount from the original transaction.
	DiscountCents *int64 `json:"discount_cents"`
	// Final amount after applying discounts and taxes, in cents.
	TotalCents *int64 `json:"total_cents"`
	// Current processing status. Can be: "outstanding" (credit issued, refund pending), "late" (refund overdue), "paid"
	// (refund completed), "void" (credit cancelled), or "archived" (removed from active credits).
	State *TransactionState `json:"state"`
	// Public URL where customers can view this transaction online.
	Permalink *string `json:"permalink"`
	// Direct download URL for the PDF version of this transaction.
	Pdf *string `json:"pdf"`
	// API endpoint URL for this transaction.
	Url *string `json:"url"`
}

type TransactionCreateResponseRelatedDocument struct {
	// Unique identifier for the related invoice or receipt.
	Id *int64 `json:"id"`
	// Type of the related document: "Invoice" or "Receipt".
	Type *DocumentType `json:"type"`
}

type TransactionCreateResponseContact struct {
	// Unique identifier for the object.
	Id *int64 `json:"id"`
	// City/District/Suburb/Town/Village.
	City *string `json:"city"`
	// If the contact is a company, this is its contact person.
	ContactPerson *string `json:"contact_person"`
	// 2-letter ISO country code.
	Country *string `json:"country"`
	// Time at which the object was created. Measured in seconds since the Unix epoch.
	CreatedAt *int64 `json:"created_at"`
	// If the contact is a company, this is the department.
	Department *string `json:"department"`
	// Default discount for this contact.
	Discount *float64 `json:"discount"`
	// The contact's email address
	Email *string `json:"email"`
	// The contact's first name. Only if the contact is a person.
	FirstName *string `json:"first_name"`
	// The type of contact.
	Kind *CustomerKind `json:"kind"`
	// The contact's preferred language. 2-letter ISO language code. Should be included in the account's translations
	// list.
	Language *string `json:"language"`
	// The contact's last name. Only if the contact is a person.
	LastName *string `json:"last_name"`
	// Internal notes about the contact.
	Notes *string `json:"notes"`
	// The contact's phone number.
	Phone *string `json:"phone_1"`
	// ZIP or postal code.
	PostalCode *string `json:"postal_code"`
	// The external platform where the contact was imported from, if applicable.
	Processor *string `json:"processor"`
	// The ID the payment_processor assigned to the contact.
	ProcessorId *string `json:"processor_id"`
	// The URL for the hosted billing area, where customers can download their invoices and update their billing
	// details.
	Permalink *string `json:"permalink"`
	// State/Province/Region.
	Region *string `json:"region"`
	// Address line 1 (Street address/PO Box).
	StreetLine1 *string `json:"street_line_1"`
	// Address line 2 (Apartment/Suite/Unit/Building).
	StreetLine2 *string `json:"street_line_2"`
	// The contact's tax identification number. Quaderno can validate EU VAT numbers, ABN, and NZBN.
	TaxId *string `json:"tax_id"`
	// Specifies the tax status of the contact.
	TaxStatus *TaxStatus `json:"tax_status"`
	// The contact's website
	Web *string `json:"web"`
	// URI of the object
	Url *string `json:"url"`
}

type TransactionCreateResponseDelivery struct {
	// Time at which the document was delivered. Measured in seconds since the Unix epoch.
	DeliveredAt *int64 `json:"delivered_at"`
	// Name of the destination system.
	Recipient *DeliveryRecipient `json:"recipient"`
	// Response from the system after delivering the document. May contain error messages.
	ServiceResponse *string `json:"service_response"`
	// Envelope type.
	Type *EnvelopeType `json:"type"`
}

type TransactionCreateResponseItem struct {
	// Unique identifier for this line item.
	Id *int64 `json:"id"`
	// Timestamp when the item was created. Measured in seconds since the Unix epoch.
	CreatedAt *int64 `json:"created_at"`
	// Description of the product or service being billed, displayed to the customer on the invoice.
	Description *string `json:"description"`
	// Total price for this line item (unit_price × quantity) before discounts or taxes, in cents.
	SubtotalCents *int64 `json:"subtotal_cents"`
	// Discount amount applied to this line item, in cents.
	DiscountCents *int64 `json:"discount_cents"`
	// Percentage discount applied to this item (0-100). Alternative to specifying discount_cents directly.
	DiscountRate *float64 `json:"discount_rate"`
	// Number of units of this product or service being billed.
	Quantity *float64 `json:"quantity"`
	// SKU or product code from your Quaderno product catalog. Used for sales tracking and reporting by product.
	ProductCode *string `json:"product_code"`
	// Primary tax classification for this item. Determines applicable tax rules and rates (e.g., "eservice" for digital
	// services, "saas" for software subscriptions).
	TaxCode *TaxCode `json:"tax_code"`
	// Country where the primary tax applies, as a 2-letter ISO code.
	TaxCountry *string `json:"tax_country"`
	// Name of the primary tax (e.g., "VAT", "GST", "Sales Tax").
	TaxName *string `json:"tax_name"`
	// Primary tax rate as a percentage.
	TaxRate *float64 `json:"tax_rate"`
	// Percentage of the subtotal used for calculating the tax amount. It's usually 100% but there are a few exceptions.
	// For example, in Texas the taxable base is 80% for SaaS products.
	TaxablePart *float64 `json:"taxable_part,string"`
	// State/province where the primary tax applies. Used for US sales tax and Canadian GST/PST.
	TaxRegion *string `json:"tax_region"`
	// Only for jurisdictions that need to apply two tax rates (e.g. some Canadian provinces).
	AdditionalTaxName *string `json:"additional_tax_name"`
	// Only for jurisdictions that need to apply two tax rates (e.g. some Canadian provinces).
	AdditionalTaxRate *float64 `json:"additional_tax_rate"`
	// Only for jurisdictions that need to apply two tax rates (e.g. some Canadian provinces).
	AdditionalTaxablePart *float64 `json:"additional_taxable_part,string"`
	// Final total for this line item after discounts and taxes. Required if unit_price is not provided. Quaderno will
	// back-calculate the base price.
	TotalAmountCents *int64 `json:"total_amount_cents"`
	// Price per unit before discounts or taxes. Required if total_amount is not provided. Quaderno will calculate the
	// final total.
	UnitPrice *float64 `json:"unit_price"`
}

type TransactionCreateResponseTax struct {
	// Display name of the tax.
	Label *string `json:"label"`
	// Tax rate applied as a percentage.
	Rate *float64 `json:"rate"`
	// Country where this tax applies, as a 2-letter ISO code.
	Country *string `json:"country"`
	// State or province where this tax applies. Relevant for US sales tax and Canadian provincial taxes (GST/PST).
	Region *string `json:"region"`
	// County where this tax applies. Only applicable for US sales tax in states with county-level taxes.
	County *string `json:"county"`
	// Tax classification that determined this tax rate (e.g., "eservice" for digital services, "saas" for software).
	TaxCode *TaxCode `json:"tax_code"`
	// Total tax amount for this tax line, in cents.
	AmountCents *int64 `json:"amount_cents"`
}

type TransactionCreateResponsePayment struct {
	// Unique identifier for the object.
	Id *int64 `json:"id"`
	// Total payment amount, in cents.
	AmountCents *int64 `json:"amount_cents"`
	// Time at which the object was created. Measured in seconds since the Unix epoch.
	CreatedAt *int64 `json:"created_at"`
	// Date of the payment.
	Date *string `json:"date"`
	// The payment method used to pay the transaction.
	PaymentMethod *PaymentMethod `json:"payment_method"`
	// The payment processor used to process the payment.
	Processor *string `json:"processor"`
	// The ID the payment_processor assigned to the payment.
	ProcessorId *string `json:"processor_id"`
	// URI of the object
	Url *string `json:"url"`
}

func (s *Transactions) Create(ctx context.Context, params *TransactionCreateParams) (
	*TransactionCreateResponse,
	error,
) {
	if params == nil {
		return nil, fmt.Errorf("create transaction params cannot be nil")
	}

	var resp TransactionCreateResponse

	err := s.client.doRequest(ctx, http.MethodPost, "/transactions", nil, params, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
