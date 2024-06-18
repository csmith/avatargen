FROM reg.c5h.io/golang AS build

WORKDIR /go/src/app

COPY . .
RUN GO111MODULE=on go install ./cmd/serve
RUN mkdir /images /uniques

FROM reg.c5h.io/base-glibc
COPY --from=build /go/bin/serve /serve
COPY --from=build --chown=65532 /images /images
COPY --from=build --chown=65532 /uniques /uniques

VOLUME /images
VOLUME /uniques
ENTRYPOINT ["/serve", "--fallback-dir", "/images", "--unique-dir", "/uniques"]
EXPOSE 8080
