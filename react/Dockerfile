FROM node:18

WORKDIR /var/www/vite-project

COPY ./project/ /var/www/vite-project/
COPY ./startup.sh /startup.sh

RUN npm install
RUN chmod 755 /startup.sh

CMD ["/startup.sh"]
