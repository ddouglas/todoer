env "local" {
  src = data.hcl_schema.app.url
  url = "postgres://todoer:password@localhost:5432/todoer?sslmode=disable"
}

data "hcl_schema" "app" {
  paths = fileset("*.pg.hcl")
}