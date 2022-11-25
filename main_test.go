package main

import (
    "reflect"
    "log"
    "encoding/json"
    "fmt"
    "testing"
    "os"
)

func TestKgrep(t *testing.T) {
    var tests = []struct {
        regex string
        match_count int
        expected []string
    }{
            { "^test$",
                0,
                []string{ "{\n \"test\": \"ing\" \n}" },
            },
            { "arr",
                0,
                []string{ "{\n \"arr\": [\"ayy\"] \n}" },
            },
            { "foo",
                0,
                []string{ "{\n \"foo\": \"bar\" \n}", "{\n \"foo\": \"bar\" \n}" },
            },
            { "foo",
                1,
                []string{ "{\n \"foo\": \"bar\" \n}"},
            },
            { "map",
                1,
                []string{ "{ \"map\": { \"romeo\": { \"juliet\": \"foo\" }}}" },
            },
            { "nested_a",
                1,
                []string{ "{ \"nested_a\": [ { \"str\": \"ing\" }, { \"foo\": \"bar\" }, { \"deep\": { \"lee\": \"woo\" } } ] }" },
            },
        }

    for _, tt := range tests {
        set_stdin("./fixtures/test.json")

        testname := fmt.Sprintf("%v,%v", tt.regex, tt.match_count)
        t.Run(testname, func(t *testing.T) {
            ans :=  Kgrep(tt.regex, tt.match_count)

            if len(tt.expected) != len(ans) {
                t.Errorf("got %v, want %v", ans, tt.expected)
            }
            for i := range tt.expected {
                var decoded_ans map[string]interface{}
                json.Unmarshal([]byte(ans[i]), &decoded_ans)
                var decoded_expected map[string]interface{}
                json.Unmarshal([]byte(tt.expected[i]), &decoded_expected)

                if reflect.DeepEqual(decoded_ans, decoded_expected) != true {
                    t.Errorf("got %v, want %v", decoded_ans, decoded_expected)
                }
            }
        })
    }
}

func set_stdin(fixture string) {
    file, err := os.Open(fixture)
    if err != nil {
        log.Fatal(err)
    }
    os.Stdin = file
}
