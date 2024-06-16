FROM cgr.dev/chainguard/wolfi-base@sha256:3eff851ab805966c768d2a8107545a96218426cee1e5cc805865505edbe6ce92 as build

USER root

RUN apk update \
    && apk add --no-cache go~1.22

WORKDIR /deepwalk

COPY . .

RUN CGO_ENABLED=0 go build -o deepwalk . && chmod +x deepwalk

FROM cgr.dev/chainguard/static@sha256:68b8855b2ce85b1c649c0e6c69f93c214f4db75359e4fd07b1df951a4e2b0140

COPY --from=build /deepwalk/deepwalk /deepwalk

ENTRYPOINT ["/deepwalk"]
CMD ["--help"]
