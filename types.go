package core

// Map generic map type for almost usage purpose. It is a shorthand
// for a map with string keys and values of any type.
type Map map[string]any

// Data generic map type for almost usage purpose.
// It is used as an alias for the Map type.
type Data Map

// JsonData generic JSON data type for almost usage purpose.
// It represents any data structure that can be marshaled as JSON.
type JsonData interface{}
