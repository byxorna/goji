FROM golang:1.6-onbuild
MAINTAINER Gabe Conradi <gabe.conradi@gmail.com>
ENTRYPOINT ["app"]
CMD ["-h"]
