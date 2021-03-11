FROM golang:1.15.8 as builder
WORKDIR '/app'
COPY . .
RUN go build -v -o app

FROM ubuntu as runner
COPY --from=builder /app/app /app
COPY --from=builder /app/config-prod.json /config.json
ENV PORT $PORT
CMD ["sh", "-c", "/app --port=$PORT"]