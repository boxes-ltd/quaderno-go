package quaderno

type TaxCode string

const (
	TaxCodeConsulting TaxCode = "consulting"
	TaxCodeEService   TaxCode = "eservice"
	TaxCodeEBook      TaxCode = "ebook"
	TaxCodeSaaS       TaxCode = "saas"
	TaxCodeStandard   TaxCode = "standard"
	TaxCodeReduced    TaxCode = "reduced"
	TaxCodeExempt     TaxCode = "exempt"
)

type TaxBehavior string

const (
	TaxBehaviorInclusive TaxBehavior = "inclusive"
	TaxBehaviorExclusive TaxBehavior = "exclusive"
)

type ProductType string

const (
	ProductTypeGood    ProductType = "good"
	ProductTypeService ProductType = "service"
)

type TaxStatus string

const (
	TaxStatusTaxable       TaxStatus = "taxable"
	TaxStatusNonTaxable    TaxStatus = "non_taxable"
	TaxStatusNotRegistered TaxStatus = "not_registered"
	TaxStatusReverseCharge TaxStatus = "reverse_charge"
)

type TransactionType string

const (
	TransactionTypeSale   TransactionType = "sale"
	TransactionTypeRefund TransactionType = "refund"
)

type PaymentMethod string

const (
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodCash         PaymentMethod = "cash"
	PaymentMethodWireTransfer PaymentMethod = "wire_transfer"
	PaymentMethodDirectDebit  PaymentMethod = "direct_debit"
	PaymentMethodCheck        PaymentMethod = "check"
	PaymentMethodIou          PaymentMethod = "iou"
	PaymentMethodPaypal       PaymentMethod = "paypal"
	PaymentMethodOther        PaymentMethod = "other"
)

type CustomerKind string

const (
	CustomerKindCompany CustomerKind = "company"
	CustomerKindPerson  CustomerKind = "person"
)

type EvidenceState string

const (
	EvidenceStateUnsettled   EvidenceState = "unsettled"
	EvidenceStateConflicting EvidenceState = "conflicting"
	EvidenceStateConfirmed   EvidenceState = "confirmed"
)

type DocumentType string

const (
	DocumentTypeInvoice DocumentType = "invoice"
	DocumentTypeReceipt DocumentType = "receipt"
)

type DeliveryRecipient string

const (
	DeliveryRecipientAraba    DeliveryRecipient = "ticketbai_araba"
	DeliveryRecipientGipuzkoa DeliveryRecipient = "ticketbai_gipuzkoa"
	DeliveryRecipientBizkaia  DeliveryRecipient = "ticketbai_bizkaia"
)

type EnvelopeType string

const (
	EnvelopeTypeRegistration EnvelopeType = "registration"
	EnvelopeTypeCancellation EnvelopeType = "cancellation"
)

type TransactionState string

const (
	TransactionStateOutstanding TransactionState = "outstanding"
	TransactionStateLate        TransactionState = "late"
	TransactionStatePaid        TransactionState = "paid"
	TransactionStateVoid        TransactionState = "void"
	TransactionStateArchived    TransactionState = "archived"
)
