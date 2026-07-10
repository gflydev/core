package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"testing"
)

// Helper function to create a test context
func createTestCtx() *Ctx {
	// Create a new fasthttp request context
	reqCtx := &fasthttp.RequestCtx{}

	// Initialize the request with empty body, query, and form
	reqCtx.Request.SetBody([]byte("{}"))

	// Set content type to JSON for body tests
	reqCtx.Request.Header.SetContentType("application/json")

	// Create a new Ctx with the request context
	return &Ctx{
		root: reqCtx,
	}
}

// Helper function to get form parameter value
func getFormParam(ctx *Ctx, key string) string {
	return ctx.PostStr(key)
}

// Helper function to get query parameter value
func getQueryParam(ctx *Ctx, key string) string {
	return ctx.QueryStr(key)
}

// Helper function to get body parameter value
func getBodyParam(ctx *Ctx, key string) (string, bool) {
	var bodyMap map[string]any

	err := ctx.ParseBody(&bodyMap)
	if err != nil {
		return "", false
	}

	// Check if the key exists in the body map
	if _, exists := bodyMap[key]; !exists {
		return "", false
	}

	// Get the JSON data from the request body
	return bodyMap[key].(string), true
}

// TestAddParamToBody tests adding parameters to the body directly
func TestAddParamToBody(t *testing.T) {
	tests := map[string]struct {
		key        string
		value      string
		expectBody string
	}{
		"Add simple key-value": {
			key:        "testKey",
			value:      "testValue",
			expectBody: "testValue",
		},
		"Add with special characters": {
			key:        "specialKey",
			value:      "value with spaces & special chars!",
			expectBody: "value with spaces & special chars!",
		},
		"Add with empty value": {
			key:        "emptyKey",
			value:      "",
			expectBody: "",
		},
		"Add with empty key": {
			key:        "",
			value:      "emptyKeyValue",
			expectBody: "emptyKeyValue",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := createTestCtx()

			// Print debug information before adding parameter
			fmt.Printf("Test case: %s\n", name)
			fmt.Printf("Key: %s, Value: %s\n", tt.key, tt.value)
			fmt.Printf("Body content before adding parameter: %s\n", string(ctx.root.PostBody()))

			// Add parameter directly to the body
			ctx.AddParam(tt.key, tt.value, ParamTypeBody)

			// Print debug information after adding parameter
			fmt.Printf("Body content after adding parameter: %s\n", string(ctx.root.PostBody()))

			// Check body parameter
			bodyValue, exists := getBodyParam(ctx, tt.key)
			fmt.Printf("Body parameter check - Key: %s, Exists: %v, Value: %s\n", tt.key, exists, bodyValue)
			assert.True(t, exists)
			assert.Equal(t, tt.expectBody, bodyValue)
		})
	}
}

func TestAddParam(t *testing.T) {
	tests := map[string]struct {
		key         string
		value       string
		paramType   []string
		checkForm   bool
		checkQuery  bool
		checkBody   bool
		expectForm  string
		expectQuery string
		expectBody  string
	}{
		"Add to form and query (default)": {
			key:         "testKey",
			value:       "testValue",
			paramType:   []string{},
			checkForm:   true,
			checkQuery:  true,
			checkBody:   false,
			expectForm:  "testValue",
			expectQuery: "testValue",
			expectBody:  "",
		},
		"Add to form only": {
			key:         "formKey",
			value:       "formValue",
			paramType:   []string{ParamTypePost},
			checkForm:   true,
			checkQuery:  false,
			checkBody:   false,
			expectForm:  "formValue",
			expectQuery: "",
			expectBody:  "",
		},
		"Add to query only": {
			key:         "queryKey",
			value:       "queryValue",
			paramType:   []string{ParamTypeQuery},
			checkForm:   false,
			checkQuery:  true,
			checkBody:   false,
			expectForm:  "",
			expectQuery: "queryValue",
			expectBody:  "",
		},
		"Add to body only": {
			key:         "bodyKey",
			value:       "bodyValue",
			paramType:   []string{ParamTypeBody},
			checkForm:   false,
			checkQuery:  false,
			checkBody:   true,
			expectForm:  "",
			expectQuery: "",
			expectBody:  "bodyValue",
		},
		"Add with special characters": {
			key:         "specialKey",
			value:       "value with spaces & special chars!",
			paramType:   []string{},
			checkForm:   true,
			checkQuery:  true,
			checkBody:   false,
			expectForm:  "value with spaces & special chars!",
			expectQuery: "value with spaces & special chars!",
			expectBody:  "",
		},
		"Add with empty value": {
			key:         "emptyKey",
			value:       "",
			paramType:   []string{},
			checkForm:   true,
			checkQuery:  true,
			checkBody:   false,
			expectForm:  "",
			expectQuery: "",
			expectBody:  "",
		},
		"Add with empty key": {
			key:         "",
			value:       "emptyKeyValue",
			paramType:   []string{},
			checkForm:   true,
			checkQuery:  true,
			checkBody:   false,
			expectForm:  "emptyKeyValue",
			expectQuery: "emptyKeyValue",
			expectBody:  "",
		},
		"Add to form and body": {
			key:         "formBodyKey",
			value:       "formBodyValue",
			paramType:   []string{ParamTypePost, ParamTypeBody},
			checkForm:   true,
			checkQuery:  false,
			checkBody:   true,
			expectForm:  "formBodyValue",
			expectQuery: "",
			expectBody:  "formBodyValue",
		},
		"Add to all types explicitly": {
			key:         "allKey",
			value:       "allValue",
			paramType:   []string{ParamTypePost, ParamTypeQuery, ParamTypeBody},
			checkForm:   true,
			checkQuery:  true,
			checkBody:   true,
			expectForm:  "allValue",
			expectQuery: "allValue",
			expectBody:  "allValue",
		},
		"Add to query and form": {
			key:         "queryBodyKey",
			value:       "queryBodyValue",
			paramType:   []string{ParamTypeQuery, ParamTypePost},
			checkForm:   true,
			checkQuery:  true,
			checkBody:   false,
			expectForm:  "queryBodyValue",
			expectQuery: "queryBodyValue",
			expectBody:  "",
		},
		"Add to query and body": {
			key:         "queryBodyKey",
			value:       "queryBodyValue",
			paramType:   []string{ParamTypeQuery, ParamTypeBody},
			checkForm:   false,
			checkQuery:  true,
			checkBody:   true,
			expectForm:  "",
			expectQuery: "queryBodyValue",
			expectBody:  "queryBodyValue",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := createTestCtx()

			// Print debug information before AddParam
			fmt.Printf("Test case: %s\n", name)
			fmt.Printf("Key: %s, Value: %s, ParamType: %v\n", tt.key, tt.value, tt.paramType)
			fmt.Printf("Body content before AddParam: %s\n", string(ctx.root.PostBody()))

			// For body parameters, use our helper function
			ctx.AddParam(tt.key, tt.value, tt.paramType...)

			// Print debug information after AddParam
			fmt.Printf("Body content after AddParam: %s\n", string(ctx.root.PostBody()))

			// Check form parameter if needed
			if tt.checkForm {
				formValue := getFormParam(ctx, tt.key)
				assert.Equal(t, tt.expectForm, formValue)
			}

			// Check query parameter if needed
			if tt.checkQuery {
				queryValue := getQueryParam(ctx, tt.key)
				assert.Equal(t, tt.expectQuery, queryValue)
			}

			// Check body parameter if needed
			if tt.checkBody {
				bodyValue, exists := getBodyParam(ctx, tt.key)
				fmt.Printf("Body parameter check - Key: %s, Exists: %v, Value: %s\n", tt.key, exists, bodyValue)
				assert.True(t, exists)
				assert.Equal(t, tt.expectBody, bodyValue)
			}
		})
	}
}

func TestDeleteParam(t *testing.T) {
	tests := map[string]struct {
		setupKey    string
		setupValue  string
		setupType   []string
		deleteKey   string
		deleteType  []string
		checkForm   bool
		checkQuery  bool
		checkBody   bool
		expectForm  bool
		expectQuery bool
		expectBody  bool
	}{
		"Full :: Remove body": {
			setupKey:    "testKey",
			setupValue:  "testValue",
			setupType:   []string{ParamTypeBody, ParamTypePost, ParamTypeQuery},
			deleteKey:   "testKey",
			deleteType:  []string{ParamTypeBody},
			checkForm:   true,
			checkQuery:  true,
			checkBody:   true,
			expectForm:  true,
			expectQuery: true,
			expectBody:  false,
		},
		"Full :: Remove query": {
			setupKey:    "testKey",
			setupValue:  "testValue",
			setupType:   []string{ParamTypeBody, ParamTypePost, ParamTypeQuery},
			deleteKey:   "testKey",
			deleteType:  []string{ParamTypeQuery},
			checkForm:   true,
			checkQuery:  true,
			checkBody:   true,
			expectForm:  true,
			expectQuery: false,
			expectBody:  true,
		},
		"Full :: Remove form": {
			setupKey:    "testKey",
			setupValue:  "testValue",
			setupType:   []string{ParamTypeBody, ParamTypePost, ParamTypeQuery},
			deleteKey:   "testKey",
			deleteType:  []string{ParamTypePost},
			checkForm:   true,
			checkQuery:  true,
			checkBody:   true,
			expectForm:  false,
			expectQuery: true,
			expectBody:  true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := createTestCtx()

			// Setup the parameter
			ctx.AddParam(tt.setupKey, tt.setupValue, tt.setupType...)

			// Call the DeleteParam method
			ctx.DeleteParam(tt.deleteKey, tt.deleteType...)

			// Check form parameter if needed
			if tt.checkForm {
				formValue := getFormParam(ctx, tt.setupKey)
				if tt.expectForm {
					assert.NotEmpty(t, formValue)
				} else {
					assert.Empty(t, formValue)
				}
			}

			// Check query parameter if needed
			if tt.checkQuery {
				queryValue := getQueryParam(ctx, tt.setupKey)
				if tt.expectQuery {
					assert.NotEmpty(t, queryValue)
				} else {
					assert.Empty(t, queryValue)
				}
			}

			// Check body parameter if needed
			if tt.checkBody {
				bodyValue, _ := getBodyParam(ctx, tt.setupKey)

				if tt.expectBody {
					assert.NotEmpty(t, bodyValue)
				} else {
					assert.Empty(t, bodyValue)
				}
			}
		})
	}
}

func TestParseQuery(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		expected Data
	}{
		{
			name: "single scalar key",
			uri:  "/?name=John",
			expected: Data{
				"name": "John",
			},
		},
		{
			name: "multiple scalar keys",
			uri:  "/?name=John&age=30&active=true",
			expected: Data{
				"name":   "John",
				"age":    "30",
				"active": "true",
			},
		},
		{
			name: "array key with bracket suffix",
			uri:  "/?tags[]=go&tags[]=web&tags[]=fast",
			expected: Data{
				"tags": []string{"go", "web", "fast"},
			},
		},
		{
			name: "single array key with bracket suffix",
			uri:  "/?tags[]=only",
			expected: Data{
				"tags": "only",
			},
		},
		{
			name: "mixed scalar and array keys",
			uri:  "/?name=John&ids[]=1&ids[]=2&active=true",
			expected: Data{
				"name":   "John",
				"ids":    []string{"1", "2"},
				"active": "true",
			},
		},
		{
			name:     "empty query string",
			uri:      "/",
			expected: Data{},
		},
		{
			name: "key with empty value",
			uri:  "/?empty=",
			expected: Data{
				"empty": "",
			},
		},
		{
			name: "repeated plain key accumulates into slice",
			uri:  "/?key=first&key=second",
			expected: Data{
				"key": []string{"first", "second"},
			},
		},
		{
			name: "repeated bracket key accumulates into slice",
			uri:  "/?key[]=first&key[]=second",
			expected: Data{
				"key": []string{"first", "second"},
			},
		},
		{
			name: "mixed plain and bracket forms for same key",
			uri:  "/?key=first&key[]=second",
			expected: Data{
				"key": []string{"first", "second"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createTestCtx()
			ctx.root.Request.SetRequestURI(tt.uri)

			data := Data{}
			err := ctx.ParseQuery(&data)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, data)
		})
	}
}

func TestQueryStrArr(t *testing.T) {
	ctx := createTestCtx()
	ctx.root.Request.SetRequestURI("/?tags=go&tags=web&tags=fast")

	result := ctx.QueryStrArr("tags")
	assert.Equal(t, []string{"go", "web", "fast"}, result)

	empty := ctx.QueryStrArr("missing")
	assert.Empty(t, empty)
}

func TestQueryIntArr(t *testing.T) {
	ctx := createTestCtx()
	ctx.root.Request.SetRequestURI("/?ids=1&ids=2&ids=3")

	result, err := ctx.QueryIntArr("ids")
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)

	empty, err := ctx.QueryIntArr("missing")
	assert.NoError(t, err)
	assert.Empty(t, empty)

	ctx2 := createTestCtx()
	ctx2.root.Request.SetRequestURI("/?ids=1&ids=bad")
	_, err = ctx2.QueryIntArr("ids")
	assert.Error(t, err)
}

func TestQueryBoolArr(t *testing.T) {
	ctx := createTestCtx()
	ctx.root.Request.SetRequestURI("/?flags=true&flags=false&flags=1")

	result, err := ctx.QueryBoolArr("flags")
	assert.NoError(t, err)
	assert.Equal(t, []bool{true, false, true}, result)

	empty, err := ctx.QueryBoolArr("missing")
	assert.NoError(t, err)
	assert.Empty(t, empty)

	ctx2 := createTestCtx()
	ctx2.root.Request.SetRequestURI("/?flags=true&flags=notbool")
	_, err = ctx2.QueryBoolArr("flags")
	assert.Error(t, err)
}

func TestQueryFloatArr(t *testing.T) {
	ctx := createTestCtx()
	ctx.root.Request.SetRequestURI("/?scores=1.1&scores=2.2&scores=3.3")

	result, err := ctx.QueryFloatArr("scores")
	assert.NoError(t, err)
	assert.Equal(t, []float64{1.1, 2.2, 3.3}, result)

	empty, err := ctx.QueryFloatArr("missing")
	assert.NoError(t, err)
	assert.Empty(t, empty)

	ctx2 := createTestCtx()
	ctx2.root.Request.SetRequestURI("/?scores=1.1&scores=notfloat")
	_, err = ctx2.QueryFloatArr("scores")
	assert.Error(t, err)
}

// bindPayload is a test type implementing Validator.
type bindPayload struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (p *bindPayload) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

func TestBind_JSONWithValidation(t *testing.T) {
	ctx := createTestCtx()
	ctx.root.Request.Header.SetContentType("application/json")
	ctx.root.Request.SetBody([]byte(`{"name":"vinh","age":30}`))

	var p bindPayload
	err := ctx.Bind(&p)

	assert.NoError(t, err)
	assert.Equal(t, "vinh", p.Name)
	assert.Equal(t, 30, p.Age)
}

func TestBind_ValidationError(t *testing.T) {
	ctx := createTestCtx()
	ctx.root.Request.Header.SetContentType("application/json")
	ctx.root.Request.SetBody([]byte(`{"age":30}`))

	var p bindPayload
	err := ctx.Bind(&p)

	assert.Error(t, err, "missing required field should fail validation")
}

// formPayload uses the encoding/json ",string" tag so a numeric field binds
// from a string form value.
type formPayload struct {
	Name string `json:"name"`
	Age  int    `json:"age,string"`
}

func TestBind_FormURLEncoded(t *testing.T) {
	ctx := createTestCtx()
	ctx.root.Request.Header.SetContentType("application/x-www-form-urlencoded")
	ctx.root.Request.SetBody([]byte("name=vinh&age=42"))

	var p formPayload
	err := ctx.Bind(&p)

	assert.NoError(t, err)
	assert.Equal(t, "vinh", p.Name)
	assert.Equal(t, 42, p.Age)
}
