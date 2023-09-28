schema "main" {}

table "tbl_solarman_credentials" {
  schema = schema.main

  column "id" {
    null    = false
    type    = int
  }

  column "username" {
    type = varchar(256)
  }

  column "password" {
    type = varchar(256)
  }

  column "app_id" {
    type = varchar(256)
  }

  column "app_secret" {
    type = varchar(256)
  }

  column "created_at" {
    null    = false
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null      = false
    type      = datetime
    default   = sql("CURRENT_TIMESTAMP")
    on_update = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }
}

table "tbl_installed_capacity" {
  schema = schema.main

  column "id" {
    null    = false
    type    = int
  }

  column "efficiency_factor" {
    type = float
  }

  column "focus_hour" {
    type = int
  }

  column "created_at" {
    null    = false
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null      = false
    type      = datetime
    default   = sql("CURRENT_TIMESTAMP")
    on_update = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }
}

table "tbl_site_region_mapping" {
  schema = schema.main

  column "id" {
    null    = false
    type    = int
  }

  column "code" {
    type = varchar(256)
  }

  column "name" {
    type = varchar(256)
  }

  column "area" {
    type = varchar(256)
    null = true
  }

  column "created_at" {
    null    = false
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null      = false
    type      = datetime
    default   = sql("CURRENT_TIMESTAMP")
    on_update = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }
}