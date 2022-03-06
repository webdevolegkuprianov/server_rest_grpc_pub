FROM golang:latest
RUN apt-get -y update && apt-get -y upgrade
RUN mkdir /server_redis
ADD . /server_redis/ 
WORKDIR /server_redis
RUN make
CMD ["./apiserver"]
RUN find . ! -name 'apiserver' -type f -exec rm -f {} +
RUN rm -R -- */




