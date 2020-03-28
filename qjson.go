package main

import (
    "fmt"
    "encoding/json"
    "errors"
    "time"
)

var EXAMPLE = []byte(`
{
    "glossary": {
        "title": "example glossary",
        "GlossDiv": {
            "title": "S",
            "GlossList": {
                "GlossEntry": {
                    "ID": "SGML",
                    "SortAs": "SGML",
                    "GlossTerm": "Standard Generalized Markup Language",
                    "Acronym": "SGML",
                    "Abbrev": "ISO 8879:1986",
                    "GlossDef": {
                        "para": "A meta-markup language, used to create markup languages such as DocBook.",
                        "GlossSeeAlso": ["GML", "XML"]
                    },
                    "GlossSee": "markup"
                }
            }
        }
    }
}
`)

type SliceResizeNeeded uint64

func NewSliceResizeNeeded(newsize uint64) SliceResizeNeeded {
    return newsize
}

func (e SliceResizeNeeded) Error() string {
    return fmt.Sprintf("Slice needs to be at least %v elements long", e)
}

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
                return u(m[k], keys[1:]...)
            }
        }
    case int:
        a, ok := V.([]interface{})
        if !ok {
            return nil, errors.New("Container type mismatch")
        }
        if l == 2 {
            // Reached path destination
            if k >= len(a) {
                return nil, errors.New("Index out of range")
            }
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
                return u(a[k], keys[1:]...)
            }
        }
    default:
        return nil, errors.New("Unknown key type")
    }
}

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
        return u(*V, keys...)
    }
}

func PrintJson(j interface{}) error {
    json, err := json.MarshalIndent(j, "", "    ")
    if err != nil {
        return err
    }
    _, err = fmt.Println(string(json))
    return err
}

func main() {
    var j interface{}
    err := json.Unmarshal(EXAMPLE, &j)
    if err != nil {
        panic(err)
    }
    fmt.Println("Example JSON:")
    err = PrintJson(j)

    // Query some JSON paths
    fmt.Println(`Q(j, "glossary", "title") ->`)
    fmt.Println(Q(j, "glossary", "title"))
    fmt.Println(`Q(j, "glossary", "non-existent-key") ->`)
    fmt.Println(Q(j, "glossary", "non-existent-key"))
    fmt.Println(`Q(j, "glossary", "GlossDiv", "GlossList", "GlossEntry", "GlossDef", "GlossSeeAlso", 1) ->`)
    fmt.Println(Q(j, "glossary", "GlossDiv", "GlossList", "GlossEntry", "GlossDef", "GlossSeeAlso", 1))

    // Apply some changes to JSON
    fmt.Println(U(&j, "glossary", "GlossDiv", "GlossList", "GlossEntry", "GlossDef", "GlossSeeAlso", 1, "ABC"))
    fmt.Println(U(&j, "glossary", "GlossDiv", "GlossList", "GlossEntry", "meta", "updated", time.Now().String()))
    fmt.Println(U(&j, "glossary", "GlossDiv", "GlossList", "GlossEntry", "GlossDef", "GlossSeeAlso", 4, "DEF"))

    fmt.Println("Edited JSON:")
    err = PrintJson(j)
    if err != nil {
        panic(err)
    }
}
