FROM us-east1-docker.pkg.dev/shadowtacticalautomation/sct/shadow-base:v0.0.1

LABEL com.shadowtactical.appname=http-injest

COPY http-injest /http-injest

EXPOSE 8090
EXPOSE 8080

ENTRYPOINT ["/http-injest"]