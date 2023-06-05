#FROM golang:1.19.4 as build

#WORKDIR /go/src/app
#COPY . .

#RUN go mod download
#RUN CGO_ENABLED=0 go build -ldflags "-w -s" -a -o /go/bin/app

# Now copy it into our base image.
#FROM gcr.io/distroless/static-debian11
#COPY --from=build /go/bin/app /
#ENV APP_SERVE_ADDR=":5040"
#EXPOSE 5040
#CMD ["/app"]

FROM gcr.io/distroless/static-debian11

#ADD main /app/main
#RUN ["chmod", "+x", "/app/main"]
COPY --chmod=755 main /app/main

ENV APP_SERVE_ADDR=":5040"
EXPOSE 5040
CMD ["/app/main"]
