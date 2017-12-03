FROM BASE_REGISTRY_URI/DAKOTA_LOGGER_IMAGE:LOGGER_TAG

COPY web/ /web/
RUN ls -la /web/*
ADD portal  /usr/local/bin/

CMD ["sh", "-c",  "{ exec portal 2>&1 1>&3 3>&- | /usr/local/bin/logger -st=tcp -ost=false; } 3>&1 1>&2 |  /usr/local/bin/logger -st=tcp"]

#Config-Auto-Generated ENV
ENV DKT_BUILD_ENV=dev
