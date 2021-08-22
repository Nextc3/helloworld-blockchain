#Consultar os ativos presentes no razão
peer chaincode query -C mychannel -n basic -c '{"Args":["GetAllAssets"]}'
