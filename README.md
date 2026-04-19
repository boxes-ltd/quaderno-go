# Quaderno Go Client

A Go API client for [Quaderno](https://quaderno.io).

This is an unofficial library and is not affiliated with Quaderno.

If you encounter a bug, please [open an issue](../../issues). Feature requests and PRs are welcome.

**Note: API responses may differ from the public documentation. These client responses were adjusted based on sandbox
testing with API version 20260309.**

## Usage

### Create a client

```go
c := quaderno.NewClient("API_KEY", "API_URL", options...)
```

### Client options

`NewClient` accepts optional configuration options via variadic arguments

#### Use a specific API version

```go
c := quaderno.NewClient("API_KEY", "API_URL", quaderno.WithApiVersion("20260309"))
```

#### Use a custom HTTP client (rather than the default HTTP client)

```go
httpClient := &http.Client{...}
c := quaderno.NewClient("API_KEY", "API_URL", quaderno.WithHttpClient(httpClient))
```

#### Set the HTTP log level

Logging options:

- `quaderno.LogLevelNone` - no logging (default)
- `quaderno.LogLevelBasic` - just shows basic request/response info
- `quaderno.LogLevelHeaders` - basic logging plus headers
- `quaderno.LogLevelBody` - basic logging plus body
- `quaderno.LogLevelHeaders | quaderno.LogLevelBody` - full logging

```go
logLevel := quaderno.LogLevelHeaders | quaderno.LogLevelBody // log everything
c := quaderno.NewClient("API_KEY", "API_URL", quaderno.WithLogLevel(logLevel))
```

#### Set a custom user agent

```go
c := quaderno.NewClient("API_KEY", "API_URL", quaderno.WithUserAgent("myApp/1.0.0"))
```

### Error handling

If the API returns an error with a non-2xx status code the client error will be a `quaderno.ApiError` type containing
the status code and response body, which can be decoded if needed:

```go
if apiErr, ok := errors.AsType[*quaderno.ApiError](err); ok {
    // TODO decode apiErr.Body JSON and evaluate it
}
```

### Check API availability

Use `Ping` to test service availability and that credentials are correct

```go
c := quaderno.NewClient("API_KEY", "API_URL")
err := c.Ping(context.Background())
if err != nil {
    fmt.Printf("Error: %s\n", err)
    return
}
fmt.Println("Ping successful!")
```

### Calculate a tax rate

Example, with country and amount:

```go
c := quaderno.NewClient("API_KEY", "API_URL")
resp, err := c.Taxes.Calculate(
    context.Background(),
    &quaderno.TaxCalculateParams{
        ToCountry: new("GB"),
        Amount: new(10.00),
    },
)
if err != nil {
    fmt.Printf("Error calculating tax: %s\n", err)
    return
}
// TODO process response
```

### Record a transaction

Example, record a sale with new customer:

```go
c := quaderno.NewClient("API_KEY", "API_URL")
resp, err := c.Transactions.Create(
    context.Background(),
    &quaderno.TransactionCreateParams{
        Type:     new(quaderno.TransactionTypeSale),
        Currency: new("GBP"),
        Customer: &quaderno.TransactionCreateCustomer{
            City:        new("London"),
            Country:     new("GB"),
            Email:       new("test@example.com"),
            FirstName:   new("John"),
            Kind:        new(quaderno.CustomerKindPerson),
            LastName:    new("Smith"),
            PostalCode:  new("SW18 3JJ"),
            StreetLine1: new("123 Lyford Road"),
        },
        Items: []*quaderno.TransactionCreateItemParams{
            {
                ProductCode: new("TEST123"),
                Description: new("Test product"),
                Quantity:    new(int64(1)),
                Amount:      new(12.00),
                Tax: &quaderno.TransactionCreateTaxParams{
                    Rate:        new(20.0),
                    TaxablePart: new(100.0),
                    Country:     new("GB"),
                },
            },
        },
        Payment: &quaderno.TransactionCreatePaymentParams{
            Method:      new(quaderno.PaymentMethodCreditCard),
            Processor:   new("CardProcessor"),
            ProcessorId: new("ABC123"),
        },
        Processor:   new("CardProcessor"),
        ProcessorId: new("ABC123"),
    },
)
if err != nil {
    fmt.Printf("Error creating transaction: %s\n", err)
    return
}
// TODO process response
```

For existing customers the Quaderno id can be used instead:

```go
resp, err := c.Transactions.Create(
    context.Background(),
    &quaderno.TransactionCreateParams{
        ...
        Customer: new(quaderno.TransactionCreateCustomerId("12345678")),
        ...
    },
)
```

#### Refunds

The [official docs](https://developers.quaderno.io/guides/record-refunds/) state that only the `processor` and
`processor_id` are required to record a refund, but in our testing this did not work. Sending a similar request to a
sale but with type `quaderno.TransactionTypeRefund` seems to work reliably.
