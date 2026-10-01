FROM docker.io/tiangolo/nginx-rtmp:latest
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 1935
