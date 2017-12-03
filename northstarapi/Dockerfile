FROM BASE_REGISTRY_URI/DAKOTA_LOGGER_IMAGE:LOGGER_TAG

ADD northstarapi  /usr/local/bin/
ADD scopes.json /northstarapi/scopes.json

CMD ["sh", "-c",  "{ exec northstarapi 2>&1 1>&3 3>&- | /usr/local/bin/logger -st=tcp -ost=false; } 3>&1 1>&2 |  /usr/local/bin/logger -st=tcp"]
