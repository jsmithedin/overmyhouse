FROM gcr.io/distroless/base:b54513ef989c81d68cb27d9c7958697e2fedd2c4

LABEL maintainer="jamie@jsmth.co.uk"

COPY .env .
COPY overmyhouse .

CMD ["./overmyhouse"]
