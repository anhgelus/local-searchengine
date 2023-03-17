FROM node:alpine as build

WORKDIR /app

COPY . .

RUN npm i && npm run build

FROM golang:1.18-alpine

WORKDIR /app

COPY --from=build /app .

RUN go mod tidy

RUN go build -o local-searchengine .

CMD ./local-searchengine

EXPOSE 8042
