FROM golang:1.15.2-alpine
RUN apk add build-base
# Add Maintainer info
LABEL maintainer="Lalu Raynaldi Pratama Putra <laluraynaldi@gmail.com>"

# persistent / runtime deps
RUN apk add --no-cache \
        wkhtmltopdf \
        xvfb \
        ttf-dejavu ttf-droid ttf-freefont ttf-liberation \
    ;

RUN ln -s /usr/bin/wkhtmltopdf /usr/local/bin/wkhtmltopdf;
RUN chmod +x /usr/local/bin/wkhtmltopdf;

RUN mkdir /paysha
ADD . /paysha
WORKDIR /paysha
## Add this go mod download command to pull in any dependencies
RUN go mod download
## Our project will now successfully build with the necessary go libraries included.
RUN go build -o main .
## Our start command which kicks off
## our newly created binary executable
EXPOSE 8888 8888
CMD ["/paysha/main"]