package core

import (
	"errors"
	"fmt"
	"github.com/gflydev/core/utils"
	"regexp"
	"sort"
	"strings"

	"github.com/valyala/bytebufferpool"
)

const (
	root nodeType = iota
	static
	param
	wildcard
)

const (
	// errSetHandler occurs when a handler is already registered for a specific path.
	errSetHandler = "a handler is already registered for path '%s'"

	// errSetWildcardHandler occurs when a wildcard handler is already registered for a specific path.
	errSetWildcardHandler = "a wildcard handler is already registered for path '%s'"

	// errWildPathConflict occurs when a new path conflicts with an existing wild path in the prefix.
	errWildPathConflict = "'%s' in new path '%s' conflicts with existing wild path '%s' in existing prefix '%s'"

	// errWildcardConflict occurs when a new path conflicts with an existing wildcard in the prefix.
	errWildcardConflict = "'%s' in new path '%s' conflicts with existing wildcard '%s' in existing prefix '%s'"

	// errWildcardSlash occurs when there is no '/' before a wildcard in the path.
	errWildcardSlash = "no / before wildcard in path '%s'"

	// errWildcardNotAtEnd occurs when wildcard routes are not at the end of the path as expected.
	errWildcardNotAtEnd = "wildcard routes are only allowed at the end of the path in path '%s'"
)

// radixError represents a custom error with a message and related parameters.
type radixError struct {
	msg    string // Error message template
	params []any  // Parameters for formatting the error message
}

// Error returns the formatted error message for a radixError.
func (err radixError) Error() string {
	return fmt.Sprintf(err.msg, err.params...)
}

// newRadixError creates a new radixError instance.
// Parameters:
//   - msg: The error message template.
//   - params: A variadic slice of parameters to format the message.
//
// Returns:
//   - radixError: The created custom error instance.
func newRadixError(msg string, params ...any) radixError {
	return radixError{msg, params}
}

// nodeType defines the type of a node in the tree.
type nodeType uint8

// nodeWildcard represents a wildcard node in the routing tree.
type nodeWildcard struct {
	path     string   // The actual wildcard path.
	paramKey string   // The name of the parameter associated with the wildcard.
	handler  IHandler // The request handler associated with the wildcard node.
}

// node represents a single node in the routing tree.
type node struct {
	nType        nodeType      // The type of the node (e.g., root, static, param, wildcard).
	path         string        // The value of the path segment for the node.
	tsr          bool          // Indicates if the node is a Trailing Slash Redirect (TSR) node.
	handler      IHandler      // The request handler associated with the node (if any).
	hasWildChild bool          // Indicates if the node has a wildcard child.
	children     []*node       // Direct child nodes of the current node.
	wildcard     *nodeWildcard // Wildcard configuration for the node, if applicable.

	paramKeys  []string       // Parameter keys for parameterized paths.
	paramRegex *regexp.Regexp // Regular expression for validating parameterized paths.
}

// wildPath represents details about a wildcard or parameterized path segment.
type wildPath struct {
	path    string         // The path of the wildcard.
	keys    []string       // Key names extracted from the wildcard path.
	start   int            // Start index of the wildcard in the path.
	end     int            // End index of the wildcard in the path.
	pType   nodeType       // Type of the wildcard (e.g., param, wildcard).
	pattern string         // Regex pattern for the wildcard.
	regex   *regexp.Regexp // Compiled regex pattern for wildcard validation.
}

// panicf formats a string with the provided arguments and triggers a panic with the resulting message.
// Parameters:
//   - s: The format string.
//   - args: Variadic arguments used to format the string.
func panicf(s string, args ...any) {
	panic(fmt.Sprintf(s, args...))
}

func bufferRemoveString(buf *bytebufferpool.ByteBuffer, s string) {
	buf.B = buf.B[:len(buf.B)-len(s)]
}

// longestCommonPrefix finds the longest common prefix between two strings.
// This function also ensures that the common prefix does not contain ':' or '*'
// characters, since those are invalid in existing keys.
//
// Parameters:
//   - a: The first string to compare.
//   - b: The second string to compare.
//
// Returns:
//   - int: The length of the longest common prefix between the two input strings.
func longestCommonPrefix(a, b string) int {
	// Compare byte-by-byte. The previous rune-based implementation compared a
	// byte offset (i += rune size) against a rune count, so it terminated early
	// for any multi-byte UTF-8 path and reported a too-short prefix, corrupting
	// node splitting. Byte comparison is correct here (identical byte prefixes
	// imply identical rune prefixes) and matches httprouter's approach.
	maxVal := min(len(a), len(b))

	i := 0
	for i < maxVal && a[i] == b[i] {
		i++
	}

	return i
}

// segmentEndIndex returns the index where the segment ends from the given path.
//
// Parameters:
//   - path: The input string path to search for the segment end.
//   - includeTSR: Boolean flag indicating whether to consider a trailing slash redirect (TSR) as part of the segment.
//
// Returns:
//   - int: The index where the current segment ends.
func segmentEndIndex(path string, includeTSR bool) int {
	end := 0
	for end < len(path) && path[end] != '/' {
		end++
	}

	if includeTSR && path[end:] == "/" {
		end++
	}

	return end
}

// findWildPath searches for a wildcard segment in a given path and validates the segment name
// for invalid characters. It analyses the path to extract wildcard details, including its name,
// type, and optional regex pattern.
//
// Parameters:
//   - path: The current path segment where the search for a wildcard begins.
//   - fullPath: The original full path from which `path` is derived. Used for error reporting.
//
// Returns:
//   - *wildPath: A pointer to a `wildPath` struct containing wildcard details such as its
//     name, position, type, and optional regex pattern. Returns nil if no wildcard is found.
func findWildPath(path, fullPath string) *wildPath {
	// Find start
	for start, c := range []byte(path) {
		// A wildcard starts with ':' (param) or '*' (wildcard)
		if c != '{' {
			continue
		}

		withRegex := false
		keys := 0

		// Find end and check for invalid characters
		for end, c := range []byte(path[start+1:]) {
			switch c {
			case '}':
				if keys > 0 {
					keys--
					continue
				}

				end := start + end + 2
				wp := &wildPath{
					path:  path[start:end],
					keys:  []string{path[start+1 : end-1]},
					start: start,
					end:   end,
					pType: param,
				}

				if len(path) > end && path[end] == '{' {
					panic("the wildcards must be separated by at least 1 char")
				}

				sn := strings.SplitN(wp.keys[0], ":", 2)
				if len(sn) > 1 {
					wp.keys = []string{sn[0]}
					pattern := sn[1]

					if pattern == "*" {
						wp.pattern = pattern
						wp.pType = wildcard
					} else {
						wp.pattern = "(" + pattern + ")"
						wp.regex = regexp.MustCompile(wp.pattern)
					}
				} else if path[len(path)-1] != '/' {
					wp.pattern = "(.*)"
				}

				if wp.keys[0] == "" {
					panicf("wildcards must be named with a non-empty name in path '%s'", fullPath)
				}

				segEnd := end + segmentEndIndex(path[end:], true)
				path = path[end:segEnd]

				if path == "/" {
					// Last segment, so include the TSR
					path = ""
					wp.end++
				}

				if path != "" {
					// Rebuild the wildpath with the prefix
					wp2 := findWildPath(path, fullPath)
					if wp2 != nil {
						prefix := path[:wp2.start]

						wp.end += wp2.end
						wp.path += prefix + wp2.path
						wp.pattern += prefix + wp2.pattern
						wp.keys = append(wp.keys, wp2.keys...)
					} else {
						wp.path += path
						wp.pattern += path
						wp.end += len(path)
					}

					wp.regex = regexp.MustCompile(wp.pattern)
				}

				return wp

			case ':':
				withRegex = true

			case '{':
				if !withRegex && keys == 0 {
					panic("the char '{' is not allowed in the param name")
				}

				keys++
			}
		}
	}

	return nil
}

// Tree is a routes storage, which organizes routing paths in a hierarchical structure.
type Tree struct {
	root *node

	// If enabled, the node handler could be updated.
	Mutable bool
}

// NewTree creates and returns a new instance of a Tree with an initialized root node.
//
// Returns:
//   - *Tree: An empty tree structure with a root node of type 'root'.
func NewTree() *Tree {
	return &Tree{
		root: &node{
			nType: root,
		},
	}
}

// Add adds a node with the specified path and associated handler to the tree.
// The path must begin with a '/' character, and the handler should not be nil.
// This function is not safe to call concurrently.
//
// WARNING: Not concurrency-safe!
//
// Parameters:
//   - path (string): The path of the node to add. It must start with '/'.
//   - handler (IHandler): The handler to associate with the specified path.
//
// Panics:
//   - If the path does not begin with a '/' character.
//   - If the handler is nil.
//
// Behavior:
//   - If the given path already exists in the tree structure and the node can
//     be updated (i.e., t.Mutable is true), the handler will be updated.
//   - If conflicts occur or the node can't be updated, the function may panic.
func (t *Tree) Add(path string, handler IHandler) {
	if !strings.HasPrefix(path, "/") {
		panicf("path must begin with '/' in path '%s'", path)
	} else if handler == nil {
		panic("nil handler")
	}

	fullPath := path

	i := longestCommonPrefix(path, t.root.path)
	if i > 0 {
		if len(t.root.path) > i {
			t.root.split(i)
		}

		path = path[i:]
	}

	n, err := t.root.add(path, fullPath, handler)
	if err != nil {
		var radixErr radixError

		if errors.As(err, &radixErr) && t.Mutable && !n.tsr {
			switch radixErr.msg {
			case errSetHandler:
				n.handler = handler
				return
			case errSetWildcardHandler:
				n.wildcard.handler = handler
				return
			}
		}

		panic(err)
	}

	if t.root.path == "" {
		t.root = t.root.children[0]
		t.root.nType = root
	}

	// Reorder the nodes
	t.root.sort()
}

// Get searches for a handler associated with the given path in the tree structure.
// It processes the given path to find the appropriate handler and updates the context
// with the values of parameters or wildcards if they exist. Additionally, it can suggest
// a trailing slash redirection if applicable.
//
// Parameters:
//   - path (string): The path to search for in the tree.
//   - ctx (*Ctx): The context to be updated with user values (e.g., parameters or wildcards).
//
// Returns:
//   - (IHandler): The handler associated with the given path, or nil if none is found.
//   - (bool): A boolean indicating if a trailing slash redirect (TSR) is recommended.
//
// Behavior:
//   - If the path matches the root or its child nodes, the relevant handler is returned.
//   - In case of trailing slash differences, it may recommend a TSR (if possible).
//   - Wildcard handlers are matched depending on the path and updated in the context.
//
// Notes:
//   - If no handler exists for the path, nil is returned along with false for TSR.
func (t *Tree) Get(path string, ctx *Ctx) (IHandler, bool) {
	if len(path) > len(t.root.path) {
		if path[:len(t.root.path)] != t.root.path {
			return nil, false
		}

		path = path[len(t.root.path):]

		return t.root.getFromChild(path, ctx)

	} else if path == t.root.path {
		switch {
		case t.root.tsr:
			return nil, true
		case t.root.handler != nil:
			return t.root.handler, false
		case t.root.wildcard != nil:
			if ctx != nil {
				ctx.Root().SetUserValue(t.root.wildcard.paramKey, "")
			}

			return t.root.wildcard.handler, false
		}
	}

	return nil, false
}

// FindCaseInsensitivePath performs a case-insensitive search for the specified path
// within the tree structure, taking into account optional trailing slash handling.
//
// Parameters:
//   - path (string): The path to search for. Case-insensitiveness ensures matches
//     regardless of the path's letter casing.
//   - fixTrailingSlash (bool): Determines if trailing slash differences should be corrected
//     during the lookup. If true, finds matches with discrepancies
//     in trailing slashes and adjusts accordingly.
//   - buf (*bytebufferpool.ByteBuffer): A buffer used for temporary path modifications
//     during the search process.
//
// Returns:
//   - (bool): A boolean indicating whether the path was successfully found:
//   - true: The path was found, optionally corrected for case or trailing slashes.
//   - false: The path was not found or could not be resolved due to differences.
func (t *Tree) FindCaseInsensitivePath(path string, fixTrailingSlash bool, buf *bytebufferpool.ByteBuffer) bool {
	found, tsr := t.root.find(path, buf)

	if !found || (tsr && !fixTrailingSlash) {
		buf.Reset()

		return false
	}

	return true
}

// newNode creates a new instance of a node with the specified path.
//
// Parameters:
//   - path (string): The value of the path segment for the new node.
//
// Returns:
//   - *node: A pointer to the newly created node with the given path.
func newNode(path string) *node {
	return &node{
		nType: static,
		path:  path,
	}
}

// conflict checks for a wildcard path conflict and raises a panic with details if a conflict is found.
//
// Parameters:
//   - path (string): The conflicting path segment causing the conflict.
//   - fullPath (string): The full path that includes the conflicting segment.
//
// Returns:
//   - error: An error detailing the wildcard path conflict.
func (n *nodeWildcard) conflict(path, fullPath string) error {
	prefix := fullPath[:strings.LastIndex(fullPath, path)] + n.path

	return newRadixError(errWildcardConflict, path, fullPath, n.path, prefix)
}

// wildPathConflict checks for a wild path conflict in the node and raises a panic with details if a conflict is found.
//
// Parameters:
//   - path (string): The conflicting path segment being checked.
//   - fullPath (string): The full path where the conflict occurs.
//
// Returns:
//   - error: An error detailing the wild path conflict.
func (n *node) wildPathConflict(path, fullPath string) error {
	pathSeg := strings.SplitN(path, "/", 2)[0]
	prefix := fullPath[:strings.LastIndex(fullPath, path)] + n.path

	return newRadixError(errWildPathConflict, pathSeg, fullPath, n.path, prefix)
}

// clone creates a deep copy of the current node and returns it.
//
// Returns:
//   - *node: A pointer to the newly created cloned node with copied data.
func (n *node) clone() *node {
	cloneNode := new(node)
	cloneNode.nType = n.nType     // Copy the node type.
	cloneNode.path = n.path       // Copy the path of the node.
	cloneNode.tsr = n.tsr         // Copy the trailing slash redirect flag.
	cloneNode.handler = n.handler // Copy the handler associated with the node.

	// Clone all child nodes recursively.
	if len(n.children) > 0 {
		cloneNode.children = make([]*node, len(n.children))
		for i, child := range n.children {
			cloneNode.children[i] = child.clone() // Recursively clone each child node.
		}
	}

	// Clone wildcard configuration if present.
	if n.wildcard != nil {
		cloneNode.wildcard = &nodeWildcard{
			path:     n.wildcard.path,     // Copy the wildcard path.
			paramKey: n.wildcard.paramKey, // Copy the wildcard's parameter key.
			handler:  n.wildcard.handler,  // Copy the handler for the wildcard.
		}
	}

	// Clone parameter keys if present.
	if len(n.paramKeys) > 0 {
		cloneNode.paramKeys = make([]string, len(n.paramKeys))
		copy(cloneNode.paramKeys, n.paramKeys) // Copy parameter keys.
	}

	cloneNode.paramRegex = n.paramRegex // Copy the parameter regex if present.

	return cloneNode
}

// split splits the current node at the specified index, creating a new child node
// with the remaining part of the path and reassigning properties accordingly.
//
// Parameters:
//   - i (int): The index at which the current node's path should be split.
//
// Behavior:
//   - Creates a cloned child node with the part of the path after the split index.
//   - Updates the current node's path to retain the part before the split index.
//   - Clears handler, TSR flag, wildcard configuration, and parameter-related properties
//     of the current node.
//   - Reassigns the cloned child node as the current node's only child.
func (n *node) split(i int) {
	cloneChild := n.clone()
	cloneChild.nType = static
	cloneChild.path = cloneChild.path[i:]
	cloneChild.paramKeys = nil
	cloneChild.paramRegex = nil

	n.path = n.path[:i]
	n.handler = nil
	n.tsr = false
	n.wildcard = nil
	n.children = append(n.children[:0], cloneChild)
}

// findEndIndexAndValues parses a path segment and extracts the ending index
// and associated parameter values based on the node's parameter regex.
//
// Parameters:
//   - path (string): The path segment to parse and extract parameter values from.
//
// Returns:
//   - (int): The index at which the matched parameter ends within the path.
//   - ([]string): A slice of extracted parameter values matched by the regex.
func (n *node) findEndIndexAndValues(path string) (int, []string) {
	index := n.paramRegex.FindStringSubmatchIndex(path)
	if len(index) == 0 || index[0] != 0 {
		return -1, nil
	}

	end := index[1]

	index = index[2:]
	values := make([]string, len(index)/2)

	i := 0
	for j := range index {
		if (j+1)%2 != 0 {
			continue
		}

		values[i] = utils.CopyStr(path[index[j-1]:index[j]])

		i++
	}

	return end, values
}

// setHandle assigns a request handler to the node, ensuring no conflicts in
// handler placement and handling trailing slash redirects.
//
// Parameters:
//   - handler (IHandler): The request handler to associate with the node.
//   - fullPath (string): The complete path for the handler being set.
//
// Returns:
//   - (*node): The current node with the handler set.
//   - (error): Error details in case of conflicts while setting the handler.
func (n *node) setHandle(handler IHandler, fullPath string) (*node, error) {
	if n.handler != nil || n.tsr {
		return n, newRadixError(errSetHandler, fullPath)
	}

	n.handler = handler
	foundTSR := false

	// Set TSR in method
	for i := range n.children {
		child := n.children[i]

		if child.path != "/" {
			continue
		}

		child.tsr = true
		foundTSR = true

		break
	}

	if n.path != "/" && !foundTSR {
		if strings.HasSuffix(n.path, "/") {
			n.split(len(n.path) - 1)
			n.tsr = true
		} else {
			childTSR := newNode("/")
			childTSR.tsr = true
			n.children = append(n.children, childTSR)
		}
	}

	return n, nil
}

// insert adds a new handler for the given path into the current node, handling
// static paths, parameters, and wildcards as needed.
//
// Parameters:
//   - path (string): The path segment to insert.
//   - fullPath (string): The full path of the handler being inserted, used for conflict detection and error reporting.
//   - handler (IHandler): The handler to associate with the specified path.
//
// Returns:
//   - (*node): A pointer to the node where the handler was successfully inserted.
//   - (error): An error if a conflict or invalid insertion occurs.
func (n *node) insert(path, fullPath string, handler IHandler) (*node, error) {
	end := segmentEndIndex(path, true)
	child := newNode(path)

	wp := findWildPath(path, fullPath)
	if wp != nil {
		j := end
		if wp.start > 0 {
			j = wp.start
		}

		child.path = path[:j]

		if wp.start > 0 {
			n.children = append(n.children, child)

			return child.insert(path[j:], fullPath, handler)
		}

		switch wp.pType {
		case param:
			n.hasWildChild = true

			child.nType = wp.pType
			child.paramKeys = wp.keys
			child.paramRegex = wp.regex
		case wildcard:
			if len(path) == end && n.path[len(n.path)-1] != '/' {
				return nil, newRadixError(errWildcardSlash, fullPath)
			} else if len(path) != end {
				return nil, newRadixError(errWildcardNotAtEnd, fullPath)
			}

			if n.path != "/" && n.path[len(n.path)-1] == '/' {
				n.split(len(n.path) - 1)
				n.tsr = true

				n = n.children[0]
			}

			if n.wildcard != nil {
				if n.wildcard.path == path {
					return n, newRadixError(errSetWildcardHandler, fullPath)
				}

				return nil, n.wildcard.conflict(path, fullPath)
			}

			n.wildcard = &nodeWildcard{
				path:     wp.path,
				paramKey: wp.keys[0],
				handler:  handler,
			}

			return n, nil
		}

		path = path[wp.end:]

		if path != "" {
			n.children = append(n.children, child)

			return child.insert(path, fullPath, handler)
		}
	}

	child.handler = handler
	n.children = append(n.children, child)

	switch {
	case child.path == "/":
		// Add TSR when split an edge and the remaining path to insert is "/"
		n.tsr = true
	case strings.HasSuffix(child.path, "/"):
		child.split(len(child.path) - 1)
		child.tsr = true
	default:
		childTSR := newNode("/")
		childTSR.tsr = true
		child.children = append(child.children, childTSR)
	}

	return child, nil
}

// add adds a handler to the node for the specified path, processing static paths,
// parameters, and wildcards as appropriate.
//
// Parameters:
//   - path (string): The path segment where the handler is to be added.
//   - fullPath (string): The complete path being added, used for conflict detection and error reporting.
//   - handler (IHandler): The request handler to associate with the specified path.
//
// Returns:
//   - (*node): A pointer to the node where the handler was successfully added.
//   - (error): An error if a conflict or invalid addition occurs.
func (n *node) add(path, fullPath string, handler IHandler) (*node, error) {
	if path == "" {
		return n.setHandle(handler, fullPath)
	}

	for _, child := range n.children {
		i := longestCommonPrefix(path, child.path)
		if i == 0 {
			continue
		}

		switch child.nType {
		case static:
			if len(child.path) > i {
				child.split(i)
			}

			if len(path) > i {
				return child.add(path[i:], fullPath, handler)
			}
		case param:
			wp := findWildPath(path, fullPath)

			isParam := wp.start == 0 && wp.pType == param
			hasHandler := child.handler != nil || handler == nil

			if len(path) == wp.end && isParam && hasHandler {
				// The current segment is a param and it's duplicated
				if child.path == path {
					return child, newRadixError(errSetHandler, fullPath)
				}

				return nil, child.wildPathConflict(path, fullPath)
			}

			if len(path) > i {
				if child.path == wp.path {
					return child.add(path[i:], fullPath, handler)
				}

				return n.insert(path, fullPath, handler)
			}
		}

		if path == "/" {
			n.tsr = true
		}

		return child.setHandle(handler, fullPath)
	}

	return n.insert(path, fullPath, handler)
}

// getFromChild traverses the child nodes to find a matching handler for the given path.
//
// Parameters:
//   - path (string): The portion of the path to match against the child nodes.
//   - ctx (*Ctx): The context object to store user values if applicable.
//
// Returns:
//   - (IHandler): The handler associated with the matched path, or nil if no match is found.
//   - (bool): A boolean indicating if a trailing slash redirect (TSR) is required.
func (n *node) getFromChild(path string, ctx *Ctx) (IHandler, bool) {
	for _, child := range n.children {
		switch child.nType {
		case static:

			// Checks if the first byte is equal
			// It's faster than compare strings
			if path[0] != child.path[0] {
				continue
			}

			if len(path) > len(child.path) {
				if path[:len(child.path)] != child.path {
					continue
				}

				h, tsr := child.getFromChild(path[len(child.path):], ctx)
				if h != nil || tsr {
					return h, tsr
				}
			} else if path == child.path {
				switch {
				case child.tsr:
					return nil, true
				case child.handler != nil:
					return child.handler, false
				case child.wildcard != nil:
					if ctx != nil {
						ctx.Root().SetUserValue(child.wildcard.paramKey, "")
					}

					return child.wildcard.handler, false
				}

				return nil, false
			}

		case param:
			end := segmentEndIndex(path, false)
			values := []string{utils.CopyStr(path[:end])}

			if child.paramRegex != nil {
				end, values = child.findEndIndexAndValues(path[:end])
				if end == -1 {
					continue
				}
			}

			if len(path) > end {
				h, tsr := child.getFromChild(path[end:], ctx)
				if tsr {
					return nil, tsr
				} else if h != nil {
					if ctx != nil {
						for i, key := range child.paramKeys {
							ctx.Root().SetUserValue(key, values[i])
						}
					}

					return h, false
				}

			} else if len(path) == end {
				switch {
				case child.tsr:
					return nil, true
				case child.handler == nil:
					// try another child
					continue
				case ctx != nil:
					for i, key := range child.paramKeys {
						ctx.Root().SetUserValue(key, values[i])
					}
				}

				return child.handler, false
			}

		default:
			panic("invalid node type")
		}
	}

	if n.wildcard != nil {
		if ctx != nil {
			ctx.Root().SetUserValue(n.wildcard.paramKey, utils.CopyStr(path))
		}

		return n.wildcard.handler, false
	}

	return nil, false
}

// find traverses the current node and its children to find a matching path,
// while appending the successfully matched segments to the provided buffer.
//
// Parameters:
//   - path (string): The path to search for within the current node and its descendants.
//   - buf (*bytebufferpool.ByteBuffer): The buffer to store matched path segments.
//
// Returns:
//   - (bool): Indicates whether a match was found.
//   - (bool): Indicates whether a trailing slash redirect (TSR) is required.
func (n *node) find(path string, buf *bytebufferpool.ByteBuffer) (bool, bool) {
	if len(path) > len(n.path) {
		if !strings.EqualFold(path[:len(n.path)], n.path) {
			return false, false
		}

		path = path[len(n.path):]
		_, err := buf.WriteString(n.path)
		if err != nil {
			return false, false
		}

		found, tsr := n.findFromChild(path, buf)
		if found {
			return found, tsr
		}

		bufferRemoveString(buf, n.path)

	} else if strings.EqualFold(path, n.path) {
		_, err := buf.WriteString(n.path)
		if err != nil {
			return false, false
		}

		if n.tsr {
			if n.path == "/" {
				bufferRemoveString(buf, n.path)
			} else {
				err := buf.WriteByte('/')
				if err != nil {
					return false, false
				}
			}

			return true, true
		}

		if n.handler != nil {
			return true, false
		} else {
			bufferRemoveString(buf, n.path)
		}
	}

	return false, false
}

// findFromChild traverses the child nodes to find a matching path segment
// and appends the successfully matched segment to the provided buffer.
//
// Parameters:
//   - path (string): The path segment to match against the child nodes.
//   - buf (*bytebufferpool.ByteBuffer): The buffer to append matched path segments.
//
// Returns:
//   - (bool): Indicates whether a match was found for the path segment.
//   - (bool): Indicates whether a trailing slash redirect (TSR) is required.
func (n *node) findFromChild(path string, buf *bytebufferpool.ByteBuffer) (bool, bool) {
	for _, child := range n.children {
		switch child.nType {
		case static:
			found, tsr := child.find(path, buf)
			if found {
				return found, tsr
			}

		case param:
			end := segmentEndIndex(path, false)

			if child.paramRegex != nil {
				end, _ = child.findEndIndexAndValues(path[:end])
				if end == -1 {
					continue
				}
			}

			_, err := buf.WriteString(path[:end])
			if err != nil {
				return false, false
			}

			if len(path) > end {
				found, tsr := child.findFromChild(path[end:], buf)
				if found {
					return found, tsr
				}

			} else if len(path) == end {
				if child.tsr {
					err := buf.WriteByte('/')
					if err != nil {
						return false, false
					}

					return true, true
				}

				if child.handler != nil {
					return true, false
				}
			}

			bufferRemoveString(buf, path[:end])

		default:
			panic("invalid node type")
		}
	}

	if n.wildcard != nil {
		_, err := buf.WriteString(path)
		if err != nil {
			return false, false
		}

		return true, false
	}

	return false, false
}

// sort recursively sorts the current node and its children nodes
// in order based on type and their children count (using `Less` function).
func (n *node) sort() {
	for _, child := range n.children {
		child.sort()
	}

	sort.Sort(n)
}

// Len returns the total number of children nodes.
//
// Returns:
//   - (int): The total count of child nodes.
func (n *node) Len() int {
	return len(n.children)
}

// Swap switches the order of two child nodes at the specified indices.
//
// Parameters:
//   - i (int): The index of the first child node.
//   - j (int): The index of the second child node.
func (n *node) Swap(i, j int) {
	n.children[i], n.children[j] = n.children[j], n.children[i]
}

// Less determines if child node at index `i` has less priority than the child node at index `j`.
// The priority is determined first by node type and, if equal, by the number of children.
//
// Parameters:
//   - i (int): The index of the first child node to compare.
//   - j (int): The index of the second child node to compare.
//
// Returns:
//   - (bool): `true` if the child node at index `i` has less priority than the one at index `j`, `false` otherwise.
func (n *node) Less(i, j int) bool {
	if n.children[i].nType < n.children[j].nType {
		return true
	} else if n.children[i].nType > n.children[j].nType {
		return false
	}

	return len(n.children[i].children) > len(n.children[j].children)
}
