# prom2lyrid
Lyrid Service to cache Prometheus Exporters 

The preferred method to build and run is using docker container by calling:

<code>
docker build .
</code>
<br />
<code>
docker run --restart always -d -p 8081:8081 %image_tag% 
</code>