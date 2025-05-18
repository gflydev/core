package core

// Map generic map type for almost usage purpose. It is shorthand
// for a map with string keys and values of any type.
type Map map[string]any

// Data generic map type for almost usage purpose.
// It is used as an alias for the Map type.
type Data Map

// Get returns the value associated with the given key.
// If the key does not exist, it returns nil.
func (d Data) Get(key string) any {
	return d[key]
}

// Set sets the value for the given key and returns the modified Data.
// If the key already exists, its value will be overwritten.
func (d Data) Set(key string, value any) Data {
	d[key] = value
	return d
}

// Has checks if the given key exists in the Data.
// Returns true if the key exists, false otherwise.
func (d Data) Has(key string) bool {
	_, ok := d[key]
	return ok
}

// Delete removes the given key and its associated value from the Data.
// Returns the modified Data. If the key does not exist, no action is taken.
func (d Data) Delete(key string) Data {
	delete(d, key)
	return d
}

func (d Data) Keys() []string {
	keys := make([]string, 0, len(d))
	for k := range d {
		keys = append(keys, k)
	}
	return keys
}

func (d Data) Vals() []any {
	values := make([]any, 0, len(d))
	for _, v := range d {
		values = append(values, v)
	}
	return values
}

// ------------------------ Extra Functions ------------------------ //

// GetInt returns the integer value associated with the given key.
// It attempts to convert various numeric types to int.
// Returns 0 if the key does not exist or value cannot be converted to int.
func (d Data) GetInt(key string) int {
	if v, ok := d.Get(key).(int); ok {
		return v
	}
	if v, ok := d.Get(key).(int64); ok {
		return int(v)
	}
	if v, ok := d.Get(key).(int32); ok {
		return int(v)
	}
	if v, ok := d.Get(key).(int16); ok {
		return int(v)
	}
	if v, ok := d.Get(key).(int8); ok {
		return int(v)
	}
	if v, ok := d.Get(key).(uint32); ok {
		return int(v)
	}
	if v, ok := d.Get(key).(uint16); ok {
		return int(v)
	}
	if v, ok := d.Get(key).(uint8); ok {
		return int(v)
	}
	return 0
}

// GetFloat returns the float64 value associated with the given key.
// It attempts to convert various numeric types to float64.
// Returns 0 if the key does not exist or value cannot be converted to float64.
func (d Data) GetFloat(key string) float64 {
	if v, ok := d.Get(key).(float64); ok {
		return v
	}
	if v, ok := d.Get(key).(float32); ok {
		return float64(v)
	}
	if v, ok := d.Get(key).(int); ok {
		return float64(v)
	}
	if v, ok := d.Get(key).(int64); ok {
		return float64(v)
	}
	if v, ok := d.Get(key).(int32); ok {
		return float64(v)
	}
	if v, ok := d.Get(key).(int16); ok {
		return float64(v)
	}
	if v, ok := d.Get(key).(int8); ok {
		return float64(v)
	}
	if v, ok := d.Get(key).(uint); ok {
		return float64(v)
	}
	if v, ok := d.Get(key).(uint64); ok {
		return float64(v)
	}
	if v, ok := d.Get(key).(uint32); ok {
		return float64(v)
	}
	if v, ok := d.Get(key).(uint16); ok {
		return float64(v)
	}
	if v, ok := d.Get(key).(uint8); ok {
		return float64(v)
	}
	return 0
}

// GetString returns the string value associated with the given key.
// Returns empty string if the key does not exist or value is not a string.
func (d Data) GetString(key string) string {
	if v, ok := d.Get(key).(string); ok {
		return v
	}
	return ""
}

// GetBool returns the boolean value associated with the given key.
// Returns false if the key does not exist or value is not a boolean.
func (d Data) GetBool(key string) bool {
	if v, ok := d.Get(key).(bool); ok {
		return v
	}
	return false
}

// GetData returns the Data value associated with the given key.
// Returns nil if the key does not exist or value is not of type Data.
func (d Data) GetData(key string) Data {
	if v, ok := d.Get(key).(Data); ok {
		return v
	}
	return nil
}
