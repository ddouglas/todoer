table "todo_nag_settings" {
  schema = schema.public

  column "id" {
    type = text
  }

  column "todo_id" {
    type = text
    null = false
  }

  column "enabled" {
    type    = boolean
    null    = false
    default = true
  }

  column "interval_minutes" {
    type    = integer
    null    = false
    default = 10
  }

  column "duration_minutes" {
    type    = integer
    null    = false
    default = 60
  }

  column "last_nagged_at" {
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

  foreign_key "fk_todo_nag_todo" {
    columns     = [column.todo_id]
    ref_columns = [table.todos.column.id]
    on_delete   = CASCADE
  }

  index "uq_todo_nag_todo_id" {
    columns = [column.todo_id]
    unique  = true
  }

  index "idx_nag_settings_enabled" {
    columns = [column.enabled, column.last_nagged_at]
    where   = "enabled = true"
  }

  index "idx_nag_settings_todo_id" {
    columns = [column.todo_id]
  }
}