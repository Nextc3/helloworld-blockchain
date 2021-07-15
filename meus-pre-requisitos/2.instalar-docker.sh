#instalar docker versão stable da distribuição
sudo apt update
sudo apt install apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu bionic stable"
#pra 20.04 
#sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable"
sudo apt update
apt-cache policy docker-ce
sudo apt install docker-ce
sudo systemctl status docker
#resolver problemas de permissão do docker
sudo setfacl -m "g:docker:rw" /var/run/docker.sock
sudo addgroup --system docker
sudo adduser $USER docker
newgrp docker
