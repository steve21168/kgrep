package main

import (
  "io"
  "fmt"
  "os"
  "log"
  "regexp"
  "encoding/json"
  "github.com/urfave/cli/v2"
)


func main() {
  app := &cli.App{
    Name:  "kgrep",
    Usage: "Used to grep keys in JSON or YAML",
    UsageText: "cat foo.json | kgrep [options] regexp",
    Flags: []cli.Flag{
      &cli.IntFlag{
        Name:    "match_count",
        Aliases: []string{"m"},
        Usage:   "The first n number of matches to return",
      },
    },
    Action: func(cCtx *cli.Context) error {
      match_count := cCtx.Int("match_count")
      user_regex := cCtx.Args().Get(0)

      kgrep(
        user_regex,
        match_count,
      )

      return nil
    },
  }

  if err := app.Run(os.Args); err != nil {
      log.Fatal(err)
  }
}

func kgrep(user_regex string, match_count int) {
  input := read_stdin()
  var key_values map[string]interface{}
  json.Unmarshal(input, &key_values)

  matches := find_matching_keys(user_regex, key_values)

  if match_count > 0 {
    matches = matches[0:match_count]
  }

  for _, match := range matches {
    fmt.Println(match)
  }
}

func read_stdin() []byte {
  stdin, err := io.ReadAll(os.Stdin)
  if err != nil {
    log.Fatal(err)
  }
  return stdin
}

func find_matching_keys(rx string, jsawn map[string]interface{}) []string {
  matches := make([]string, 0, 100)
  for k, v := range jsawn {
    switch v.(type) {
    case map[string]interface{}:
      matches = append(matches, find_matching_keys(rx, v.(map[string]interface{}))...)
    case []interface{}:
      for _, item := range v.([]interface{}) {
        switch item.(type) {
        case map[string]interface{}:
          matches = append(matches, find_matching_keys(rx, item.(map[string]interface{}))...)
        }
      }
    }

    rx, _ := regexp.Compile(rx)
    if rx.MatchString(k) {
      return_json_map := map[string]interface{}{
        k: v,
      }
      json_match, _ := json.MarshalIndent(return_json_map, "", "  ")
      matches = append(matches, string(json_match))
    }
  }

  return matches
}
