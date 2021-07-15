apt install py-pip
apt install python3-dev
apt install libffi-dev
apt install openssl-dev
apt install gcc
apt install libc-dev
apt install rust
apt install cargo
apt install make.

sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
#sudo curl \ -L https://raw.githubusercontent.com/docker/compose/1.29.2/contrib/completion/bash/docker-compose \ -o /etc/bash_completion.d/docker-compose
#inicializar o docker
sudo systemctl start docker
#sempre que iniciar o sistema
sudo systemctl enable docker