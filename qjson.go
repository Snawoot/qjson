// Helper routines for JSON manipulation in Go.
package qjson

import (
    "fmt"
)

type sliceResizeNeeded uint64

func newSliceResizeNeeded(newsize uint64) sliceResizeNeeded {
    return sliceResizeNeeded(newsize)
}

func (e sliceResizeNeeded) Error() string {
    return fmt.Sprintf("Slice needs to be at least %v elements long", int64(e))
}

// This error is returned when string key is not found in map.
type KeyError string

func newKeyError(key string) KeyError {
    return KeyError(key)
}

func (e KeyError) Error() string {
    return fmt.Sprintf("Key \"%s\" not found", string(e))
}

// Returns absent key name
func (e KeyError) Key() string {
    return string(e)
}

// This error is returned when array index is out of range.
type IndexError int

func newIndexError(index int) IndexError {
    return IndexError(index)
}

func (e IndexError) Error() string {
    return fmt.Sprintf("Index \"%d\" is out of range", int(e))
}

// Returns absent index value.
func (e IndexError) Index() int {
    return int(e)
}

// This error is returned in case when function parameters are incorrect.
type ArgError string

func newArgError(msg string) ArgError {
    return ArgError(msg)
}

func (e ArgError) Error() string {
    return string(e)
}

// This error is returned on mismatch of data types.
type TypeError string

func newTypeError(msg string) TypeError {
    return TypeError(msg)
}

func (e TypeError) Error() string {
    return string(e)
}

func s(keys ...interface{}) (interface{}, error) {
    if len(keys) == 0 {
        return nil, newArgError("No values passed")
    } else if len(keys) == 1 {
        return keys[0], nil
    }
    key := keys[0]
    switch k := key.(type) {
    case string:
        m := make(map[string]interface{})
        elem, err := s(keys[1:]...)
        if err != nil {
            return nil, err
        }
        m[k] = elem
        return m, nil
    case int:
        if k < 0 {
            return nil, newArgError("Negative index is not allowed")
        }
        a := make([]interface{}, k+1)
        elem, err := s(keys[1:]...)
        if err != nil {
            return nil, err
        }
        a[k] = elem
        return a, nil
    default:
        return nil, newTypeError("Unknown key type")
    }
}

// Query some JSON paths.
// Invocation: Q(object {}interface, path... interface{}, newvalue interface{}).
// Returns value and error.
func Q(V interface{}, keys ...interface{}) (interface{}, error) {
    if len(keys) == 0 {
        return V, nil
    }
    key := keys[0]

    var next interface{}
    switch k := key.(type) {
    case string:
        v, ok := V.(map[string]interface{})
        if !ok {
            return nil, newTypeError("Bad container type: not a map")
        }
        next, ok = v[k]
        if !ok {
            return nil, newKeyError(k)
        }
    case int:
        v, ok := V.([]interface{})
        if !ok {
            return nil, newTypeError("Bad container type: not an array")
        }
        if len(v) <= k || k < 0 {
            return nil, newIndexError(k)
        }
        next = v[k]
    default:
        return nil, newTypeError("Unknown key type")
    }
    return Q(next, keys[1:]...)
}

func u(V interface{}, keys ...interface{}) (interface{}, error) {
    if V == nil {
        // Should never happen if this function is called only by U()
        return nil, newTypeError("Can't update nil value")
    }
    l := len(keys)
    if l < 2 {
        return nil, newArgError("Incorrect arg length")
    }
    key := keys[0]
    switch k := key.(type) {
    case string:
        m, ok := V.(map[string]interface{})
        if !ok {
            return nil, newTypeError("Container type mismatch")
        }
        if l == 2 {
            // Reached path destination
            old := m[k]
            m[k] = keys[1]
            return old, nil
        } else {
            // Follow next container
            if m[k] == nil {
                // Recreate subtree
                tree, err := s(keys[1:]...)
                if err != nil {
                    return nil, err
                }
                m[k] = tree
                return nil, nil
            } else {
                res, err := u(m[k], keys[1:]...)
                if size, ok := err.(sliceResizeNeeded) ; ok {
                    // Handle slice resize
                    newslice := make([]interface{}, size)
                    copy(newslice, m[k].([]interface{}))
                    m[k] = newslice
                    // Retry with resized array
                    return u(m[k], keys[1:]...)
                }
                return res, err
            }
        }
    case int:
        a, ok := V.([]interface{})
        if !ok {
            return nil, newTypeError("Container type mismatch")
        }
        if k < 0 {
            return nil, newIndexError(k)
        }
        if k >= len(a) {
            return nil, newSliceResizeNeeded(uint64(k + 1))
        }
        if l == 2 {
            // Reached path destination
            old := a[k]
            a[k] = keys[1]
            return old, nil
        } else {
            // Follow next container
            if a[k] == nil {
                // Recreate subtree
                tree, err := s(keys[1:]...)
                if err != nil {
                    return nil, err
                }
                a[k] = tree
                return nil, nil
            } else {
                res, err := u(a[k], keys[1:]...)
                if size, ok := err.(sliceResizeNeeded) ; ok {
                    // Handle slice resize
                    newslice := make([]interface{}, size)
                    copy(newslice, a[k].([]interface{}))
                    a[k] = newslice
                    // Retry with resized array
                    return u(a[k], keys[1:]...)
                }
                return res, err
            }
        }
    default:
        return nil, newTypeError("Unknown key type")
    }
}

// Apply some changes to JSON.
// Invocation: U(object {}interface, path... interface{}, newvalue interface{}).
// Returns old value and error.
func U(V *interface{}, keys ...interface{}) (interface{}, error) {
    if V == nil {
        return nil, newArgError("nil pointer dereference")
    }
    l := len(keys)
    if l < 1 {
        return nil, newArgError("Incorrect arg length")
    } else if l == 1 {
        oldval := *V
        *V = keys[0]
        return oldval, nil
    } else {
        if *V == nil {
            tree, err := s(keys...)
            *V = tree
            return nil, err
        }
        res, err := u(*V, keys...)
        if size, ok := err.(sliceResizeNeeded) ; ok {
            // Handle slice resize
            newslice := make([]interface{}, size)
            copy(newslice, (*V).([]interface{}))
            *V = newslice
            // Retry with resized array
            return u(*V, keys...)
        }
        return res, err
    }
}
