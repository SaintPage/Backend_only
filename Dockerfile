# Etapa de construcción
FROM golang:1.20 as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o series-tracker-backend .

# Etapa final
FROM alpine:latest
# Instala bash
RUN apk add --no-cache bash
WORKDIR /root/
COPY --from=builder /app/series-tracker-backend .
COPY wait-for-it.sh .
RUN chmod +x wait-for-it.sh
EXPOSE 8080

CMD ["./wait-for-it.sh", "db:5432", "--", "./series-tracker-backend"]