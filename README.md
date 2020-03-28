# qjson

Helper routines for JSON manipulation in Go

## Usage

```go
package main

import (
    "fmt"
    "encoding/json"
    "time"
    "github.com/Snawoot/qjson"
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

func PrintJson(j interface{}) error {
    json, err := json.MarshalIndent(j, "", "    ")
    if err != nil {
        return err
    }
    _, err = fmt.Println(string(json))
    return err
}

func main() {
    // Read and write parsed JSON
    var j interface{}
    err := json.Unmarshal(EXAMPLE, &j)
    if err != nil {
        panic(err)
    }
    fmt.Println("Example JSON:")
    err = PrintJson(j)

    fmt.Println("")
    fmt.Println("Queries:")
    fmt.Println(`Q(j, "glossary", "title") ->`)
    fmt.Println(qjson.Q(j, "glossary", "title"))
    fmt.Println(`Q(j, "glossary", "non-existent-key") ->`)
    fmt.Println(qjson.Q(j, "glossary", "non-existent-key"))
    fmt.Println(`Q(j, "glossary", "GlossDiv", "GlossList", "GlossEntry", "GlossDef", "GlossSeeAlso", 1) ->`)
    fmt.Println(qjson.Q(j, "glossary", "GlossDiv", "GlossList", "GlossEntry", "GlossDef", "GlossSeeAlso", 1))

    fmt.Println("")
    fmt.Println("Updates:")
    fmt.Println(`U(&j, "glossary", "GlossDiv", "GlossList", "GlossEntry", "GlossDef", "GlossSeeAlso", 1, "ABC") -> `)
    fmt.Println(qjson.U(&j, "glossary", "GlossDiv", "GlossList", "GlossEntry", "GlossDef", "GlossSeeAlso", 1, "ABC"))
    fmt.Println(`U(&j, "glossary", "GlossDiv", "GlossList", "GlossEntry", "meta", "updated", time.Now().Truncate(0).String()) -> `)
    fmt.Println(qjson.U(&j, "glossary", "GlossDiv", "GlossList", "GlossEntry", "meta", "updated", time.Now().Truncate(0).String()))
    fmt.Println(`U(&j, "glossary", "GlossDiv", "GlossList", "GlossEntry", "GlossDef", "GlossSeeAlso", 4, "DEF") -> `)
    fmt.Println(qjson.U(&j, "glossary", "GlossDiv", "GlossList", "GlossEntry", "GlossDef", "GlossSeeAlso", 4, "DEF"))

    fmt.Println("Edited JSON:")
    err = PrintJson(j)
    if err != nil {
        panic(err)
    }

    // Compose JSON from scratch
    var k interface{}
    qjson.U(&k, "menu", "id", "file")
    qjson.U(&k, "menu", "value", "File")
    qjson.U(&k, "menu", "popup", "menuitem", 0, "value", "New")
    qjson.U(&k, "menu", "popup", "menuitem", 0, "onclick", "CreateNewDoc()")
    qjson.U(&k, "menu", "popup", "menuitem", 1, "value", "Open")
    qjson.U(&k, "menu", "popup", "menuitem", 1, "onclick", "OpenDoc()")
    qjson.U(&k, "menu", "popup", "menuitem", 2, "value", "Close")
    qjson.U(&k, "menu", "popup", "menuitem", 2, "onclick", "CloseDoc()")

    fmt.Println("Composed JSON:")
    err = PrintJson(k)
    if err != nil {
        panic(err)
    }
}
```
