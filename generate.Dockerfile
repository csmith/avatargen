FROM reg.c5h.io/golang AS build

WORKDIR /go/src/app

COPY . .
RUN GO111MODULE=on go install ./cmd/generate
RUN mkdir /images /uniques

FROM reg.c5h.io/base-glibc
COPY --from=build /go/bin/generate /generate
COPY --from=build --chown=65532 /images /images
COPY --from=build --chown=65532 /uniques /uniques

VOLUME /images
VOLUME /uniques
ENTRYPOINT ["/generate", "--dirs", "/images:1000-2500,/uniques:1000-2000"]
EXPOSE 8080
