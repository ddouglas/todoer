schema "public" {}

table "categories" {
  schema = schema.public

  column "id" {
    type = text
  }

  column "name" {
    type = text
    null = false
  }

  column "color" {
    type    = text
    null    = false
    default = "#3b82f6"
  }

  column "icon" {
    type = text
    null = true
  }

  column "sort_order" {
    type    = integer
    null    = false
    default = 0
  }

  column "created_at" {
    type    = timestamp
    null    = false
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }
}

table "todos" {
  schema = schema.public

  column "id" {
    type = text
  }

  column "title" {
    type = text
    null = false
  }

  column "description" {
    type = text
    null = true
  }

  column "completed" {
    type    = boolean
    null    = false
    default = false
  }

  column "priority" {
    type    = text
    null    = false
    default = "medium"
  }

  column "category_id" {
    type = text
    null = true
  }

  column "due_date" {
    type = timestamp
    null = true
  }

  column "reminder_at" {
    type = timestamp
    null = true
  }

  column "created_at" {
    type    = timestamp
    null    = false
    default = sql("now()")
  }

  column "updated_at" {
    type    = timestamp
    null    = false
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_category" {
    columns     = [column.category_id]
    ref_columns = [table.categories.column.id]
    on_delete   = SET_NULL
    on_update   = CASCADE
  }

  index "idx_category_id" {
    columns = [column.category_id]
  }

  index "idx_due_date" {
    columns = [column.due_date]
  }

  index "idx_completed" {
    columns = [column.completed]
  }
}