target {
  hub X {
    key = "source.json.fields.varchar[a].path"
    #key = source.json.fields.varchar[a].path
  }
  result = "test-${upper(var.A)}"
  #result = var.A
}