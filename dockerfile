FROM registry.access.redhat.com/ubi8/go-toolset:1.16.12-10 as builder
COPY --chown=1001:1001 ./ /src
WORKDIR /src
RUN go build

FROM registry.access.redhat.com/ubi8/ubi-micro:8.5-836
COPY --chown=1001:1001 --from=builder /src/authfox /authfox
COPY --chown=1001:1001 --from=builder /src/swagger /swagger
USER 1001
ENTRYPOINT [ "/authfox" ]
