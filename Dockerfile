FROM ubuntu:18.04

ADD drone-release /usr/local/bin/drone-release
RUN apt-get update -y 
RUN apt-get -y --no-install-recommends install subversion ca-certificates language-pack-zh-hans
RUN dpkg-reconfigure -f noninteractive locales
RUN rm -fr \
    /var/cache/* \
    /var/lib/apt/lists/* \
    /root/.cache/ \
    /tmp/* 

CMD ["drone-release"]