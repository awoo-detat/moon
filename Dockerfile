FROM golang:1.13 AS build

WORKDIR /app
COPY . .

#RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"'
RUN go build

#FROM scratch
#COPY --from=build /app/moon-werewolf .
#ENTRYPOINT ["./moon-werewolf"]
ENTRYPOINT ["/app/moon"]

