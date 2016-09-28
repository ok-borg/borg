docker build -t crufter/borg-web .
docker rm -f borg-web
cp ./nginx.conf /Users/crufter/nginx
docker run -p=80:80 -v /var/sitemap/sitemap.xml.gz:/usr/share/nginx/html/sitemap.xml.gz:ro -v /Users/crufter/nginx/nginx.conf:/etc/nginx/nginx.conf:ro --name=borg-web -d crufter/borg-web
