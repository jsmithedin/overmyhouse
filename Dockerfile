FROM gcr.io/distroless/base

LABEL maintainer="jamie@jsmth.co.uk"

COPY overmyhouse .

CMD ["overmyhouse"]
