Using psql Manually

From your terminal:

# Connect to your database
psql "postgres://super_shop_dev_user:QmaDNHGpVtdD8sCv40MIvZFono48XZrW@localhost:5432/super_shop_dev_db?sslmode=disable"


Then run:

```\i E:\ongoing-projects\PS-Ecommerce\pse-api-v1\migrations\001_init.up.sql;```


```\i```→ executes SQL file in current session