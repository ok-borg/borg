docker build -t crufter/borg-test-web .
docker push crufter/borg-test-web
ansible-playbook ../daemon/ansible/deploy/test-web.yml -i ../daemon/ansible/hosts/live.hosts -vvvv
