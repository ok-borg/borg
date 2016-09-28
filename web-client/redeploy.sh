docker build -t crufter/borg-web .
docker push crufter/borg-web
ansible-playbook ../daemon/ansible/deploy/web.yml -i ../daemon/ansible/hosts/live.hosts -vvvv
