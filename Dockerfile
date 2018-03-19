FROM scratch
ADD iix.se-golang-backend /

ENV WEBROOT ""
ENV DBHOST ""
ENV DBUSER ""
ENV DBPASS ""
ENV DBNAME ""
ENV JWT ""

EXPOSE 80

CMD ["/iix.se-golang-backend"]
