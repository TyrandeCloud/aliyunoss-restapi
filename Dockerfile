FROM ubuntu:20.04

RUN runtimeDeps=" \
            tzdata \
        " \
  && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y $runtimeDeps

#ENV TZ=Asia/Shanghai
#RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
#RUN dpkg-reconfigure -f noninteractive tzdata

RUN apt autoremove && apt clean

WORKDIR /server
COPY conf conf
COPY etc etc
COPY data/json json
COPY oss-proxy .

CMD ["/server/oss-proxy"]