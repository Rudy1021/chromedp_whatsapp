FROM golang:1.19-alpine as stage1
WORKDIR /app
copy . .

RUN go mod download
RUN go build

FROM alpine
WORKDIR /app
RUN mkdir data
RUN apk update && \
  apk add chromium
COPY --from=stage1 /app/chromedp_whatsapp /app

CMD ["./chromedp_whatsapp"]