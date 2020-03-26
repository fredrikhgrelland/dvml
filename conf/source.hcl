source {
  json "sirius.json" {
    fields {
      varchar A {
        #type = "smepath"
        path = "$.a"
      }
      varchar B {
        path = "$.b"
      }
      number X {
        path = "$.x"
      }
    }
  }
}