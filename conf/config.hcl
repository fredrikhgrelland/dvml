source {
  json "test.json" {
    fields {
      varchar A {
        path = "$.a"
      }
      varchar B {
        path = "$.b"
      }
    }
  }
}
target {
  hub X {
    key = "source.json.fields.varchar[a].path"
  }
}