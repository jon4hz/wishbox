FROM gcr.io/distroless/static
COPY wishbox /usr/local/bin/wishbox
WORKDIR /app
ENTRYPOINT [ "/usr/local/bin/wishbox" ]