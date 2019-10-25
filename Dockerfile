FROM gcr.io/distroless/base

LABEL maintainer="jamie@jsmth.co.uk"

COPY .env .
COPY overmyhouse .

CMD ["./overmyhouse"]
