## PostgreSQL based rest-api in Go

# Installations
go install github.com/joho/godotenv/cmd/godotenv@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
go get -u github.com/go-chi/chi/v5
go get -u github.com/golang-jwt/jwt/v5
go get github.com/google/uuid


# Creating Database
according to your db_url 
cd sql/schema
goose postgres postgres://postgres:Password@localhost:5432/dailytask up
cd ../..
sqlc generate

- Getting tasks of user

![dailytask-tasks2](https://github.com/aalperen0/daily-task/assets/88675716/e1f3cb26-d425-45f1-81e7-7cd0185fc37b)

![dailytask-tasks](https://github.com/aalperen0/daily-task/assets/88675716/31445ce8-eb80-4870-bf6f-9560bf2e99e7)
