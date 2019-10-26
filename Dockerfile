FROM gcr.io/distroless/base:latest

LABEL maintainer="jamie@jsmth.co.uk"

COPY .env .
COPY overmyhouse .

CMD ["./overmyhouse"]
