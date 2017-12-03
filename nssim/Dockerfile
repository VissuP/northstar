FROM BASE_REGISTRY_URI/DAKOTA_LOGGER_IMAGE:LOGGER_TAG

ADD nssim  /usr/local/bin/
ADD config.json /nssim/config.json
ADD ./web/dist/ /web/dist/
RUN ls -la /web/*
CMD ["sh", "-c", "{ exec nssim 2>&1 1>&3 3>&- | /usr/local/bin/logger -st=tcp -ost=false; } 3>&1 1>&2 | /usr/local/bin/logger -st=tcp"]
