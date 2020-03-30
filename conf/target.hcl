target {
  hub X {
    #key = "source.json.fields.varchar[a].path"
    key = key(var.A, var.B)
    key = "test-${upper(var.A)}"
    date = nows
    //computed = now()
  }
}