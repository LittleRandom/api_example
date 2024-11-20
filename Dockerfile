# Builds the go modules as a binary
FROM golang AS build

WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o /bin/api_app --tags=docker -a -ldflags '-linkmode external -extldflags "-static"' ./main.go


# Final image to host application
FROM alpine

COPY --from=build /bin/api_app /app
# Create the base data and library directories.
RUN mkdir /data


EXPOSE 5050


# Run
CMD ["/app"]