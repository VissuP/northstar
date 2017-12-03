FROM node:7.9.0

RUN npm i -g @angular/cli@1.0.0-beta.28.3

RUN mkdir /app
WORKDIR /app
ADD entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
