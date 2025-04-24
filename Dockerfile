FROM gcr.io/distroless/static-debian12
USER nonroot:nonroot
COPY --chown=nonroot:nonroot mcp-nomad /
ENTRYPOINT ["/mcp-nomad"]