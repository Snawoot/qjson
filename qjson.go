// Helper routines for JSON manipulation in Go

package qjson

import (
    "fmt"
    "errors"
)

type SliceResizeNeeded uint64

func NewSliceResizeNeeded(newsize uint64) SliceResizeNeeded {
    return SliceResizeNeeded(newsize)
}

func (e SliceResizeNeeded) Error() string {
    return fmt.Sprintf("Slice needs to be at least %v elements long", e)
}

// Query some JSON paths
// Invocation: Q(object {}interface, path... interface{}, newvalue interface{})
// Returns value and error
func S(keys ...interface{}) (interface{}, error) {
    if len(keys) == 0 {
        return nil, errors.New("No values passed")
    } else if len(keys) == 1 {
        return keys[0], nil
    }
    key := keys[0]
    switch k := key.(type) {
    case string:
        m := make(map[string]interface{})
        elem, err := S(keys[1:]...)
        if err != nil {
            return nil, err
        }
        m[k] = elem
        return m, nil
    case int:
        if k < 0 {
            return nil, errors.New("Negative index is not allowed")
        }
        a := make([]interface{}, k+1)
        elem, err := S(keys[1:]...)
        if err != nil {
            return nil, err
        }
        a[k] = elem
        return a, nil
    default:
        return nil, errors.New("Unknown key type")
    }
}

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
            return nil, errors.New("Bad container type: not a map")
        }
        next, ok = v[k]
        if !ok {
            return nil, errors.New("Key not found")
        }
    case int:
        v, ok := V.([]interface{})
        if !ok {
            return nil, errors.New("Bad container type: not an array")
        }
        if len(v) <= k {
            return nil, errors.New("Index out of range")
        }
        next = v[k]
    default:
        return nil, errors.New("Unknown key type")
    }
    return Q(next, keys[1:]...)
}

func u(V interface{}, keys ...interface{}) (interface{}, error) {
    if V == nil {
        // Should never happen if this function is called only by U()
        return nil, errors.New("Can't update nil value")
    }
    l := len(keys)
    if l < 2 {
        return nil, errors.New("Incorrect arg length")
    }
    key := keys[0]
    switch k := key.(type) {
    case string:
        m, ok := V.(map[string]interface{})
        if !ok {
            return nil, errors.New("Container type mismatch")
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
                tree, err := S(keys[1:]...)
                if err != nil {
                    return nil, err
                }
                m[k] = tree
                return nil, nil
            } else {
                res, err := u(m[k], keys[1:]...)
                if size, ok := err.(SliceResizeNeeded) ; ok {
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
            return nil, errors.New("Container type mismatch")
        }
        if k >= len(a) {
            return nil, NewSliceResizeNeeded(uint64(k + 1))
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
                tree, err := S(keys[1:]...)
                if err != nil {
                    return nil, err
                }
                a[k] = tree
                return nil, nil
            } else {
                res, err := u(a[k], keys[1:]...)
                if size, ok := err.(SliceResizeNeeded) ; ok {
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
        return nil, errors.New("Unknown key type")
    }
}

// Apply some changes to JSON
// Invocation: U(object {}interface, path... interface{}, newvalue interface{})
// Returns old value and error
func U(V *interface{}, keys ...interface{}) (interface{}, error) {
    if V == nil {
        return nil, errors.New("nil pointer dereference")
    }
    l := len(keys)
    if l < 1 {
        return nil, errors.New("Incorrect arg length")
    } else if l == 1 {
        oldval := *V
        *V = keys[0]
        return oldval, nil
    } else {
        if *V == nil {
            tree, err := S(keys...)
            *V = tree
            return nil, err
        }
        res, err := u(*V, keys...)
        if size, ok := err.(SliceResizeNeeded) ; ok {
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
