# Dockerfile
FROM debian
RUN mkdir /opt/splitshift 
WORKDIR /opt/splitshift
COPY splitshift /opt/splitshift/splitshift
COPY error.gtpl /opt/splitshift/
COPY upload.gtpl /opt/splitshift/
EXPOSE 8000
CMD ["/opt/splitshift/splitshift"]
