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
