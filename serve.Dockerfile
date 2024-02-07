FROM golang:1.22 AS build

WORKDIR /go/src/app

COPY . .
RUN GO111MODULE=on go install ./cmd/serve
RUN mkdir /images /uniques

FROM gcr.io/distroless/base:nonroot
COPY --from=build /go/bin/serve /serve
COPY --from=build --chown=nonroot /images /images
COPY --from=build --chown=nonroot /uniques /uniques

VOLUME /images
VOLUME /uniques
ENTRYPOINT ["/serve", "--fallback-dir", "/images", "--unique-dir", "/uniques"]
EXPOSE 8080
