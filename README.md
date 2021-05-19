# prom2lyrid

Full Readme - Coming soon

Reads the local Prometheus metric exporters, caches and propagates them (as a text) to the proxy by calling its HTTP endpoint.

The preferred method to build and run is using docker container by calling:

<code>
docker build .
</code>
<br />
<code>
docker run --restart always -d -p 8081:8081 %image_tag% 
</code>
