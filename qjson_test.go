package qjson

import (
    "testing"
    "encoding/json"
    "strings"
    "time"
)

var EXAMPLE = `
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
`

var EXAMPLE2 = `
{
    "menu": {
        "id": "file",
        "popup": {
            "menuitem": [
                {
                    "onclick": "CreateNewDoc()",
                    "value": "New"
                },
                {
                    "onclick": "OpenDoc()",
                    "value": "Open"
                },
                {
                    "onclick": "CloseDoc()",
                    "value": "Close"
                }
            ]
        },
        "value": "File"
    }
}
`



func loadJSON(data string, t *testing.T) interface{} {
    var j interface{}
    err := json.Unmarshal([]byte(data), &j)
    if err != nil {
        t.FailNow()
    }
    return j
}

func dumpJSON(data interface{}, t *testing.T) string {
    d, err := json.Marshal(data)
    if err != nil {
        t.FailNow()
    }
    return string(d)
}

func TestQueryEmptyMap(t *testing.T) {
    j := loadJSON("{}", t)
    q, err := Q(j)
    if err != nil {
        t.Fail()
    }
    if m, ok := q.(map[string]interface{}); !ok || len(m) != 0 {
        t.Fail()
    }

    _, err = Q(j, "somekey")
    switch e := err.(type) {
    case KeyError:
        if e.Key() != "somekey" {
            t.Fail()
        }
        if !strings.Contains(e.Error(), "somekey") {
            t.Fail()
        }
    default:
        t.Fail()
    }

    _, err = Q(j, 1)
    switch err.(type) {
    case TypeError:
    default:
        t.Fail()
    }

    _, err = Q(j, -1)
    switch err.(type) {
    case TypeError:
    default:
        t.Fail()
    }

    _, err = Q(j, -1, "aaa")
    switch err.(type) {
    case TypeError:
    default:
        t.Fail()
    }
}

func TestQueryEmptyList(t *testing.T) {
    j := loadJSON("[]", t)
    q, err := Q(j)
    if err != nil {
        t.Fail()
    }
    if m, ok := q.([]interface{}) ; !ok || len(m) != 0 {
        t.Fail()
    }

    _, err = Q(j, "somekey")
    switch err.(type) {
    case TypeError:
    default:
        t.Fail()
    }

    _, err = Q(j, 0)
    switch err.(type) {
    case IndexError:
    default:
        t.Fail()
    }

    _, err = Q(j, 1)
    switch err.(type) {
    case IndexError:
    default:
        t.Fail()
    }

    _, err = Q(j, -1)
    switch e := err.(type) {
    case IndexError:
        if e.Index() != -1 {
            t.Fail()
        }
    default:
        t.Fail()
    }

    _, err = Q(j, -1, "aaa")
    switch err.(type) {
    case IndexError:
    default:
        t.Fail()
    }
}

func TestQueryNull(t *testing.T) {
    j := loadJSON("null", t)

    q, err := Q(j)
    if err != nil {
        t.Fail()
    }
    if q != nil {
        t.Fail()
    }

    q, err = Q(j, "aaa")
    switch err.(type) {
    case TypeError:
    default:
        t.Fail()
    }
}

func TestQueryObj(t *testing.T) {
    j := loadJSON(`{"a":{"b":{"c": "d"}}}`, t)

    q, err := Q(j)
    if err != nil {
        t.Fail()
    }
    _, ok := q.(map[string]interface{})
    if !ok {
        t.Fail()
    }

    q, err = Q(j, "a")
    if err != nil {
        t.Fail()
    }
    _, ok = q.(map[string]interface{})
    if !ok {
        t.Fail()
    }

    q, err = Q(j, "a", "b")
    if err != nil {
        t.Fail()
    }
    _, ok = q.(map[string]interface{})
    if !ok {
        t.Fail()
    }

    q, err = Q(j, "a", "a")
    _, ok = err.(KeyError)
    if !ok {
        t.Fail()
    }

    q, err = Q(j, "a", 0)
    _, ok = err.(TypeError)
    if !ok {
        t.Fail()
    }

    q, err = Q(j, "a", "b", "c")
    if err != nil {
        t.Fail()
    }
    s, ok := q.(string)
    if !ok {
        t.Fail()
    }
    if s != "d" {
        t.Fail()
    }
}

func TestQueryObjWithList(t *testing.T) {
    j := loadJSON(`{"a":{"b":[true, false, null]}}`, t)

    q, err := Q(j, "a", "b", 0)
    if err != nil {
        t.Fail()
    }
    b, ok := q.(bool)
    if !ok {
        t.Fail()
    }
    if !b {
        t.Fail()
    }

    q, err = Q(j, "a", "b", 1)
    if err != nil {
        t.Fail()
    }
    b, ok = q.(bool)
    if !ok {
        t.Fail()
    }
    if b {
        t.Fail()
    }

    q, err = Q(j, "a", "b", 2)
    if err != nil {
        t.Fail()
    }
    if q != nil {
        t.Fail()
    }
}

func TestQueryObjWithListWithObject(t *testing.T) {
    j := loadJSON(`{"a":{"b":[true, false, {"c": "d"}]}}`, t)

    q, err := Q(j, "a", "b", 0)
    if err != nil {
        t.Fail()
    }
    b, ok := q.(bool)
    if !ok {
        t.Fail()
    }
    if !b {
        t.Fail()
    }

    q, err = Q(j, "a", "b", 1)
    if err != nil {
        t.Fail()
    }
    b, ok = q.(bool)
    if !ok {
        t.Fail()
    }
    if b {
        t.Fail()
    }

    q, err = Q(j, "a", "b", 2, "c")
    if err != nil {
        t.Fail()
    }
    if s, ok := q.(string) ; ok {
        if s != "d" {
            t.Fail()
        }
    } else {
        t.Fail()
    }
}

func TestUpdateEmpty(t *testing.T) {
    var j interface{}

    old, err := U(&j)
    if old != nil {
        t.Fail()
    }
    switch err.(type) {
    case ArgError:
    default:
        t.Fail()
    }

    old, err = U(&j, "abc")
    if old != nil {
        t.Fail()
    }
    s, ok := j.(string)
    if !ok {
        t.Fail()
    }
    if s != "abc" {
        t.Fail()
    }

    old, err = U(&j, "def")
    s, ok = old.(string)
    if !ok {
        t.Fail()
    }
    if s != "abc" {
        t.Fail()
    }
    s, ok = j.(string)
    if !ok {
        t.Fail()
    }
    if s != "def" {
        t.Fail()
    }
}

func TestUpdateExamples(t *testing.T) {
    j := loadJSON(EXAMPLE, t)
    now := time.Now().Truncate(0).String()
    U(&j, "glossary", "GlossDiv", "GlossList", "GlossEntry", "GlossDef", "GlossSeeAlso", 1, "ABC")
    U(&j, "glossary", "GlossDiv", "GlossList", "GlossEntry", "meta", "updated", now)
    U(&j, "glossary", "GlossDiv", "GlossList", "GlossEntry", "GlossDef", "GlossSeeAlso", 4, "DEF")

    q, err := Q(j, "glossary", "GlossDiv", "GlossList", "GlossEntry", "GlossDef", "GlossSeeAlso", 1)
    s, ok := q.(string)
    if err != nil || !ok || s != "ABC" {
        t.Fail()
    }
    q, err = Q(j, "glossary", "GlossDiv", "GlossList", "GlossEntry", "meta", "updated")
    s, ok = q.(string)
    if err != nil || !ok || s != now {
        t.Fail()
    }
    q, err = Q(j, "glossary", "GlossDiv", "GlossList", "GlossEntry", "GlossDef", "GlossSeeAlso", 4)
    s, ok = q.(string)
    if err != nil || !ok || s != "DEF" {
        t.Fail()
    }
    q, err = Q(j, "glossary", "GlossDiv", "GlossList", "GlossEntry", "GlossDef", "GlossSeeAlso", 3)
    if err != nil || q != nil {
        t.Fail()
    }
}

func TestResizeOuter(t *testing.T) {
    var j interface{}
    U(&j, 0, "a")
    U(&j, 1, "b")
    U(&j, 2, "c")
    ref := []string{"a", "b", "c"}
    a, ok := j.([]interface{})
    if !ok {
        t.Fail()
    }
    if len(a) != len(ref) {
        t.Fail()
    }
    for i, v := range a {
        s, ok := v.(string)
        if !ok || s != ref[i] {
            t.Fail()
        }
    }
}

func TestResizeInner(t *testing.T) {
    j := loadJSON(`{"a":[[]]}`, t)
    refdump := dumpJSON(loadJSON(`{"a":[[true]]}`, t), t)
    U(&j, "a", 0, 0, true)
    if dumpJSON(j, t) != refdump {
        t.Fail()
    }
}

func TestUpdateNullPtr(t *testing.T) {
    _, err := U(nil, 0, "a")
    _, ok := err.(ArgError)
    if !ok {
        t.Fail()
    }
}

func TestUpdateBadKey(t *testing.T) {
    var j interface{}
    _, err := U(&j, 0, .0, "a")
    _, ok := err.(TypeError)
    if !ok {
        t.Fail()
    }
}

func TestUpdateBadKey2(t *testing.T) {
    j := loadJSON(`{"a":{"b": {"c": null}}}`, t)
    _, err := U(&j, "a", "b", .0, false)
    _, ok := err.(TypeError)
    if !ok {
        t.Fail()
    }
}

func TestQueryBadKey(t *testing.T) {
    j := loadJSON(`{"a":[]}`, t)
    _, err := Q(&j, .0)
    _, ok := err.(TypeError)
    if !ok {
        t.Fail()
    }
}

func TestSliceExcPrintable(t *testing.T) {
    e := newSliceResizeNeeded(100)
    s := e.Error()
    if !strings.Contains(s, "100") {
        t.Fail()
    }
}

func TestIndexErrorPrintable(t *testing.T) {
    e := newIndexError(100)
    s := e.Error()
    if !strings.Contains(s, "100") {
        t.Fail()
    }
}

func TestArgErrorPrintable(t *testing.T) {
    e := newArgError("abcdef")
    s := e.Error()
    if s != "abcdef" {
        t.Fail()
    }
}

func TestTypeErrorPrintable(t *testing.T) {
    e := newTypeError("abcdef")
    s := e.Error()
    if s != "abcdef" {
        t.Fail()
    }
}

func TestUpdateWithNegativeIndex(t *testing.T) {
    var j interface{}
    _, err := U(&j, -1, 1)
    if _, ok := err.(ArgError) ; !ok {
        t.Fail()
    }
}

func TestDeepUpdateWithNegativeIndex(t *testing.T) {
    var j interface{}
    _, err := U(&j, 0, "abc", "def", -1, 1)
    if _, ok := err.(ArgError) ; !ok {
        t.Fail()
    }
}

func TestBadSubtreeConstruct(t *testing.T) {
    _, err := s()
    if _, ok := err.(ArgError) ; !ok {
        t.Fail()
    }
}

func TestArtificalBadInnerUpdate(t *testing.T) {
    _, err := u(nil, 1, 2, 3)
    if _, ok := err.(TypeError) ; !ok {
        t.Fail()
    }
}

func TestArtificalInnerUpdateWithNoValue(t *testing.T) {
    j := loadJSON(`{}`, t)
    _, err := u(j, "test")
    if _, ok := err.(ArgError) ; !ok {
        t.Fail()
    }
}

func TestContainerTypeMismatch(t *testing.T) {
    j := loadJSON(`{}`, t)
    _, err := U(&j, 0, true)
    if _, ok := err.(TypeError) ; !ok {
        t.Fail()
    }
    j = loadJSON(`[]`, t)
    _, err = U(&j, "test", true)
    if _, ok := err.(TypeError) ; !ok {
        t.Fail()
    }
}

func TestJSONFromScratch(t *testing.T) {
	refj := loadJSON(EXAMPLE2, t)
	refstr := dumpJSON(refj, t)
    var k interface{}
    U(&k, "menu", "id", "file")
    U(&k, "menu", "value", "File")
    U(&k, "menu", "popup", "menuitem", 0, "value", "New")
    U(&k, "menu", "popup", "menuitem", 0, "onclick", "CreateNewDoc()")
    U(&k, "menu", "popup", "menuitem", 1, "value", "Open")
    U(&k, "menu", "popup", "menuitem", 1, "onclick", "OpenDoc()")
    U(&k, "menu", "popup", "menuitem", 2, "value", "Close")
    U(&k, "menu", "popup", "menuitem", 2, "onclick", "CloseDoc()")
    dump := dumpJSON(k, t)
    if dump != refstr {
        t.Fail()
    }
}

func TestInnerSubtreeConstructFail(t *testing.T) {
    j := loadJSON(EXAMPLE2, t)
    _, err := U(&j, "menu", "aaa", -1, nil)
    if _, ok := err.(ArgError) ; !ok {
        t.Fail()
    }
}

func TestInnerIndexError(t *testing.T) {
    j := loadJSON(`[[[[]]]]`, t)
    _, err := U(&j, 0, -1, nil)
    if _, ok := err.(IndexError) ; !ok {
        t.Fail()
    }
}

func TestInnerSubtreeConstructFail2(t *testing.T) {
    j := loadJSON(EXAMPLE2, t)
    _, err := U(&j, "menu", "popup", "menuitem", 3, -1, nil)
    if _, ok := err.(ArgError) ; !ok {
        t.Fail()
    }
}

func TestQBool(t *testing.T) {
    j := loadJSON(`[true, false, null, false]`, t)
    b, err := QBool(j, 0)
    if err != nil || !b {
        t.Fail()
    }
    b, err = QBool(j, 1)
    if err != nil || b {
        t.Fail()
    }
    _, err = QBool(j, 4)
    if err == nil {
        t.Fail()
    }
    _, err = QBool(j, 2)
    if _, ok := err.(TypeError) ; !ok {
        t.Fail()
    }
}

func TestQNumber(t *testing.T) {
    j := loadJSON(`[0, 1, null, 3]`, t)
    b, err := QNumber(j, 0)
    if err != nil || b != 0 {
        t.Fail()
    }
    b, err = QNumber(j, 1)
    if err != nil || b != 1 {
        t.Fail()
    }
    _, err = QNumber(j, 4)
    if err == nil {
        t.Fail()
    }
    _, err = QNumber(j, 2)
    if _, ok := err.(TypeError) ; !ok {
        t.Fail()
    }
}

func TestQString(t *testing.T) {
    j := loadJSON(`["0", "1", 2, "3"]`, t)
    b, err := QString(j, 0)
    if err != nil || b != "0" {
        t.Fail()
    }
    b, err = QString(j, 1)
    if err != nil || b != "1" {
        t.Fail()
    }
    _, err = QString(j, 4)
    if err == nil {
        t.Fail()
    }
    _, err = QString(j, 2)
    if _, ok := err.(TypeError) ; !ok {
        t.Fail()
    }
}

func TestQList(t *testing.T) {
    j := loadJSON(`[[], [true], 2, [true,true,true]]`, t)
    b, err := QList(j, 0)
    if err != nil || len(b) != 0 {
        t.Fail()
    }
    b, err = QList(j, 1)
    if err != nil || len(b) != 1 {
        t.Fail()
    }
    _, err = QList(j, 4)
    if err == nil {
        t.Fail()
    }
    _, err = QList(j, 2)
    if _, ok := err.(TypeError) ; !ok {
        t.Fail()
    }
}

func TestQObject(t *testing.T) {
    j := loadJSON(`[{}, {"a": true}, 2, {"a": true, "b":true, "c": true}]`, t)
    b, err := QObject(j, 0)
    if err != nil || len(b) != 0 {
        t.Fail()
    }
    b, err = QObject(j, 1)
    if err != nil || len(b) != 1 {
        t.Fail()
    }
    _, err = QObject(j, 4)
    if err == nil {
        t.Fail()
    }
    _, err = QObject(j, 2)
    if _, ok := err.(TypeError) ; !ok {
        t.Fail()
    }
}

func TestQNull(t *testing.T) {
    j := loadJSON(`[null, null, 0, null]`, t)
    err := QNull(j, 0)
    if err != nil {
        t.Fail()
    }
    err = QNull(j, 1)
    if err != nil {
        t.Fail()
    }
    err = QNull(j, 4)
    if err == nil {
        t.Fail()
    }
    err = QNull(j, 2)
    if _, ok := err.(TypeError) ; !ok {
        t.Fail()
    }
}
