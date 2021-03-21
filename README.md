# Application simple endpoint

## How to run the application
Clone the application with `git@github.com:pmuaddib/sedb.git`

MySQL must be installed

1. Use the MySQL dump `sql-dump.sql` to create the table.

2. Once the application is cloned and database is created, change configurations in file config/config.yaml.

3. To run the application, please use the following command -

    ```shell
        go run main.go
    ```
> Note: By default the port number its being run on is **8082**.

## Config

Config file location exists in **config/config.yaml**

An example:
```shell
# Server
server:
  host: ""
  port: "8082"
# Database
db:
  host: "127.0.0.1"
  port: "3306"
  user: "user"
  pass: "pass"
  dbname: "sedb"
```

## Endpoint exposed

1. Send
   ```shell
   curl -X POST -H "Content-Type: application/json" -d '{"user_id": "134256", "currency": "EUR", "amount": 1000, "time_placed" : "24-JAN-20 10:27:44", "type": "deposit"}' 127.0.0.1:8082/send
   ```
2. Get
   ```shell
   curl 127.0.0.1:8082/get?userId=134256
   ```

## Tests

   ```shell
      go test
   ```