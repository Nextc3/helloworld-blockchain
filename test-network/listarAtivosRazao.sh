#Consultar os ativos presentes no razão
source variaveis_de_ambienteOrg1.sh
peer chaincode query -C mychannel -n helloworld -c '{"Args":["QueryAllOis"]}'