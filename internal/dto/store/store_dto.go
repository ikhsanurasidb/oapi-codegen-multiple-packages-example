// Package store provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package store

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

const (
	Api_keyScopes = "api_key.Scopes"
)

// Defines values for OrderStatus.
const (
	Approved  OrderStatus = "approved"
	Delivered OrderStatus = "delivered"
	Placed    OrderStatus = "placed"
)

// Order defines model for Order.
type Order struct {
	Complete *bool      `json:"complete,omitempty"`
	Id       *int64     `json:"id,omitempty"`
	PetId    *int64     `json:"petId,omitempty"`
	Quantity *int32     `json:"quantity,omitempty"`
	ShipDate *time.Time `json:"shipDate,omitempty"`

	// Status Order Status
	Status *OrderStatus `json:"status,omitempty"`
}

// OrderStatus Order Status
type OrderStatus string

// PlaceOrderJSONRequestBody defines body for PlaceOrder for application/json ContentType.
type PlaceOrderJSONRequestBody = Order

// PlaceOrderFormdataRequestBody defines body for PlaceOrder for application/x-www-form-urlencoded ContentType.
type PlaceOrderFormdataRequestBody = Order

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
