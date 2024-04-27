FROM cgr.dev/chainguard/wolfi-base as build

USER root

RUN apk update \
    && apk add --no-cache make go

WORKDIR /deepwalk

COPY . .

RUN CGO_ENABLED=0 go build -o deepwalk . && chmod +x deepwalk

FROM cgr.dev/chainguard/static

COPY --from=build /deepwalk/deepwalk /deepwalk

ENTRYPOINT ["/deepwalk"]
CMD ["--help"]
