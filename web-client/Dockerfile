FROM nginx
COPY . /usr/share/nginx/html
RUN sed -i -e 's/ok-b\.org:9992/borg.crufter.com/g' /usr/share/nginx/html/app.js 
