FROM golang:1.22 AS build

WORKDIR /go/src/app

COPY . .
RUN GO111MODULE=on go install ./cmd/generate
RUN mkdir /images /uniques

FROM gcr.io/distroless/base:nonroot
COPY --from=build /go/bin/generate /generate
COPY --from=build --chown=nonroot /images /images
COPY --from=build --chown=nonroot /uniques /uniques

VOLUME /images
VOLUME /uniques
ENTRYPOINT ["/generate", "--dirs", "/images:1000-2500,/uniques:1000-2000"]
EXPOSE 8080
