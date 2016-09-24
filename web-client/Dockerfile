FROM nginx
COPY . /usr/share/nginx/html
RUN sed -i -e 's/ok-b\.org:9992/crufter.borg.com/g' /usr/share/nginx/html/app.js 
