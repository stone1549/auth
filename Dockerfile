FROM golang:1.20

ENV AUTH_SERVICE_ENVIRONMENT=DEV
ENV AUTH_SERVICE_REPO_TYPE=POSTGRESQL
ENV AUTH_SERVICE_TIMEOUT=60
ENV AUTH_SERVICE_PORT=3333
ENV AUTH_SERVICE_PG_URL=postgres://postgres:postgres@yapyapyap-db:5432/postgres?sslmode=disable
ENV AUTH_SERVICE_TOKEN_SECRET=SECRET!

WORKDIR /go/src/github.com/stone1549/yapyapyap/auth/
COPY . .

RUN go mod tidy

CMD ["go", "run", "main.go"]

