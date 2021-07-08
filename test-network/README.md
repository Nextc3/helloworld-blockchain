#Rede de Teste com Scripts otimizados 

## Executando a rede de teste

Você pode usar o script `./network.sh` para criar uma rede de teste de malha simples. A rede de teste tem duas peer organizações  com um peer cada e um serviço de pedido de raft(etcd-raft) de nó único. Você também pode usar o script `./network.sh` para criar canais e implantar chaincode. Para obter mais informações, consulte [Usando a rede de teste do Fabric] (https://hyperledger-fabric.readthedocs.io/en/latest/test_network.html). A rede de teste está sendo introduzida no Fabric v2.0 como um substituto de longo prazo para o exemplo `first-network`.

Antes de implantar a rede de teste, você precisa seguir as instruções para [Instalar as amostras, binários e imagens do Docker] (https://hyperledger-fabric.readthedocs.io/en/latest/install.html) no Hyperledger Fabric documentação.

## Running the test network

You can use the `./network.sh` script to stand up a simple Fabric test network. The test network has two peer organizations with one peer each and a single node raft ordering service. You can also use the `./network.sh` script to create channels and deploy chaincode. For more information, see [Using the Fabric test network](https://hyperledger-fabric.readthedocs.io/en/latest/test_network.html). The test network is being introduced in Fabric v2.0 as the long term replacement for the `first-network` sample.

Before you can deploy the test network, you need to follow the instructions to [Install the Samples, Binaries and Docker Images](https://hyperledger-fabric.readthedocs.io/en/latest/install.html) in the Hyperledger Fabric documentation.
