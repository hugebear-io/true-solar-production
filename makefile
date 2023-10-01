apply_schema:
	atlas schema apply --url "sqlite://database.db" --to "file://script/schema.hcl"

api:
	go run ./cmd/api/main.go

dp:
	go run ./cmd/daily_production/main.go

mp:
	go run ./cmd/monthly_production/main.go

pa:
	go run ./cmd/performance_alarm/main.go

mock_invt:
	go run ./cmd/mock_solarman/main.go

daily_build:
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -ldflags "-linkmode external" -o daily ./cmd/daily_production/main.go

monthly_build:
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -ldflags "-linkmode external" -o monthly ./cmd/monthly_production/main.go

pa_build:
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -ldflags "-linkmode external" -o performance_alarm ./cmd/performance_alarm/main.go

mock_invt_build:
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -ldflags "-linkmode external" -o mock_solarman ./cmd/mock_solarman/main.go

invt_build:
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -ldflags "-linkmode external" -o solarman ./cmd/solarman/main.go