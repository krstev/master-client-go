FROM ofayau/ejre:8-jre
ADD bin/goclient goclient
ADD templates/*.json templates/

ENTRYPOINT ["./goclient"]