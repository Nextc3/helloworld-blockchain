#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

# This script brings up a Hyperledger Fabric network for testing smart contracts
# and applications. The test network consists of two organizations with one
# peer each, and a single node Raft ordering service. Users can also use this
# script to create a channel deploy a chaincode on the channel
#
# prepending $PWD/../bin to PATH to ensure we are picking up the correct binaries
# this may be commented out to resolve installed version of tools if desired
#pt-br
# Este script traz uma rede Hyperledger Fabric para testar contratos inteligentes
# e aplicativos. A rede de teste consiste em duas organizações com uma
# peer each, e um único nó serviço de pedidos de Raft. Os usuários também podem usar este
# script para criar um canal implantar um chaincode no canal
#
# prefixando $ PWD /../ bin para PATH para garantir que estamos pegando os binários corretos
# isso pode ser comentado para resolver a versão instalada das ferramentas, se desejado
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=${PWD}/configtx
export VERBOSE=false

. scripts/utils.sh

# Obtain CONTAINER_IDS and remove them
# This function is called when you bring a network down

## Obtenha CONTAINER_IDS e remova-os
# Esta função é chamada quando você desativa uma rede
# function clearContainers() {
#   infoln "Removing remaining containers"
#   infoln "Removendo recipientes restantes"
#   docker rm -f $(docker ps -aq --filter label=service=hyperledger-fabric) 2>/dev/null || true
#   docker rm -f $(docker ps -aq --filter name='dev-peer*') 2>/dev/null || true
# }
function clearContainers() {
  CONTAINER_IDS=$(docker ps -a | awk '($2 ~ /dev-peer.*/) {print $1}')
  if [ -z "$CONTAINER_IDS" -o "$CONTAINER_IDS" == " " ]; then
    infoln "No containers available for deletion"
  else
    docker rm -f $CONTAINER_IDS
  fi
}
# Delete any images that were generated as a part of this setup
# specifically the following images are often left behind:
# This function is called when you bring the network down

# Exclua todas as imagens que foram geradas como parte desta configuração
# especificamente, as seguintes imagens costumam ser deixadas para trás:
# Esta função é chamada quando você desativa a rede
# function removeUnwantedImages() {
#   infoln "Removing generated chaincode docker images"
#   infoln "Removendo imagens do docker do chaincode"
#   docker image rm -f $(docker images -aq --filter reference='dev-peer*') 2>/dev/null || true
# }
function removeUnwantedImages() {
  DOCKER_IMAGE_IDS=$(docker images | awk '($1 ~ /dev-peer.*/) {print $3}')
  if [ -z "$DOCKER_IMAGE_IDS" -o "$DOCKER_IMAGE_IDS" == " " ]; then
    infoln "No images available for deletion"
  else
    docker rmi -f $DOCKER_IMAGE_IDS
  fi
}
# Versions of fabric known not to work with the test network

# Versões de tecido que não funcionam com a rede de teste
NONWORKING_VERSIONS="^1\.0\. ^1\.1\. ^1\.2\. ^1\.3\. ^1\.4\."

# Do some basic sanity checking to make sure that the appropriate versions of fabric
# binaries/images are available. In the future, additional checking for the presence
# of go or other items could be added.

# Faça algumas verificações básicas de sanidade para se certificar de que as versões adequadas de tecido
# binários / imagens estão disponíveis. No futuro, verificação adicional de presença
# de go ou outros itens podem ser adicionados.
function checkPrereqs() {
  ## Check if your have cloned the peer binaries and configuration files.

  ## Verifique se você clonou os binários de pares e arquivos de configuração
  peer version > /dev/null 2>&1

  if [[ $? -ne 0 || ! -d "../config" ]]; then
    errorln "Peer binary and configuration files not found.."
    errorln "Arquivos binários e de configuração de mesmo nível não encontrados."
    errorln
    errorln "Follow the instructions in the Fabric docs to install the Fabric Binaries:"
    errorln "Siga as instruções nos documentos do Fabric para instalar os binários do Fabric:"
    errorln "https://hyperledger-fabric.readthedocs.io/en/latest/install.html"
    exit 1
  fi
  # use the fabric tools container to see if the samples and binaries match your
  # use o contêiner de ferramentas de tecido para ver se as amostras e binários correspondem ao seu
  # docker images
  LOCAL_VERSION=$(peer version | sed -ne 's/ Version: //p')
  DOCKER_IMAGE_VERSION=$(docker run --rm hyperledger/fabric-tools:$IMAGETAG peer version | sed -ne 's/ Version: //p' | head -1)

  infoln "LOCAL_VERSION=$LOCAL_VERSION"
  infoln "DOCKER_IMAGE_VERSION=$DOCKER_IMAGE_VERSION"

  if [ "$LOCAL_VERSION" != "$DOCKER_IMAGE_VERSION" ]; then
    warnln "Local fabric binaries and docker images are out of  sync. This may cause problems."
    warnln "Os binários da malha local e as imagens do docker estão fora de sincronia. Isso pode causar problemas."
  fi

  for UNSUPPORTED_VERSION in $NONWORKING_VERSIONS; do
    infoln "$LOCAL_VERSION" | grep -q $UNSUPPORTED_VERSION
    if [ $? -eq 0 ]; then
      fatalln "Local Fabric binary version of $LOCAL_VERSION does not match the versions supported by the test network."
    fi

    infoln "$DOCKER_IMAGE_VERSION" | grep -q $UNSUPPORTED_VERSION
    if [ $? -eq 0 ]; then
      fatalln "Fabric Docker image version of $DOCKER_IMAGE_VERSION does not match the versions supported by the test network."
    fi
  done

  ## Check for fabric-ca
  if [ "$CRYPTO" == "Certificate Authorities" ]; then

    fabric-ca-client version > /dev/null 2>&1
    if [[ $? -ne 0 ]]; then
      errorln "fabric-ca-client binary not found.."
      errorln
      errorln "Follow the instructions in the Fabric docs to install the Fabric Binaries:"
      errorln "https://hyperledger-fabric.readthedocs.io/en/latest/install.html"
      exit 1
    fi
    CA_LOCAL_VERSION=$(fabric-ca-client version | sed -ne 's/ Version: //p')
    CA_DOCKER_IMAGE_VERSION=$(docker run --rm hyperledger/fabric-ca:$CA_IMAGETAG fabric-ca-client version | sed -ne 's/ Version: //p' | head -1)
    infoln "CA_LOCAL_VERSION=$CA_LOCAL_VERSION"
    infoln "CA_DOCKER_IMAGE_VERSION=$CA_DOCKER_IMAGE_VERSION"

    if [ "$CA_LOCAL_VERSION" != "$CA_DOCKER_IMAGE_VERSION" ]; then
      warnln "Local fabric-ca binaries and docker images are out of sync. This may cause problems."
    fi
  fi
}

# Before you can bring up a network, each organization needs to generate the crypto
# material that will define that organization on the network. Because Hyperledger
# Fabric is a permissioned blockchain, each node and user on the network needs to
# use certificates and keys to sign and verify its actions. In addition, each user
# needs to belong to an organization that is recognized as a member of the network.
# You can use the Cryptogen tool or Fabric CAs to generate the organization crypto
# material.

# By default, the sample network uses cryptogen. Cryptogen is a tool that is
# meant for development and testing that can quickly create the certificates and keys
# that can be consumed by a Fabric network. The cryptogen tool consumes a series
# of configuration files for each organization in the "organizations/cryptogen"
# directory. Cryptogen uses the files to generate the crypto  material for each
# org in the "organizations" directory.

# You can also use Fabric CAs to generate the crypto material. CAs sign the certificates
# and keys that they generate to create a valid root of trust for each organization.
# The script uses Docker Compose to bring up three CAs, one for each peer organization
# and the ordering organization. The configuration file for creating the Fabric CA
# servers are in the "organizations/fabric-ca" directory. Within the same directory,
# the "registerEnroll.sh" script uses the Fabric CA client to create the identities,
# certificates, and MSP folders that are needed to create the test network in the
# "organizations/ordererOrganizations" directory.

# Create Organization crypto material using cryptogen or CAs
# Antes de criar uma rede, cada organização precisa gerar a criptografia
# material que definirá essa organização na rede. Porque Hyperledger
# Fabric é uma blockchain com permissão, cada nó e usuário na rede precisa
# use certificados e chaves para assinar e verificar suas ações. Além disso, cada usuário
# precisa pertencer a uma organização reconhecida como membro da rede.
# Você pode usar a ferramenta Cryptogen ou Fabric CAs para gerar a criptografia da organização
# material.

# Por padrão, a rede de amostra usa cryptogen. Cryptogen é uma ferramenta que é
# destinado a desenvolvimento e teste que podem criar rapidamente os certificados e chaves
# que pode ser consumido por uma rede Fabric. A ferramenta criptogênica consome uma série
# de arquivos de configuração para cada organização em "organization / cryptogen"
# diretório. O Cryptogen usa os arquivos para gerar o material criptográfico para cada
# org no diretório "organization".

# Você também pode usar Fabric CAs para gerar o material criptográfico. CAs assinam os certificados
# e as chaves que eles geram para criar uma raiz válida de confiança para cada organização.
# O script usa Docker Compose para trazer três CAs, um para cada organização de mesmo nível
# e a organização solicitante. O arquivo de configuração para criar o Fabric CA
# servidores estão no diretório "organization / fabric-ca". Dentro do mesmo diretório,
# o script "registerEnroll.sh" usa o cliente Fabric CA para criar as identidades,
# certificados e pastas MSP necessários para criar a rede de teste no
# diretório "organization / ordererOrganizations".

# Criar material criptográfico da organização usando criptogênicos ou CAs
function createOrgs() {
  if [ -d "organizations/peerOrganizations" ]; then
    rm -Rf organizations/peerOrganizations && rm -Rf organizations/ordererOrganizations
  fi

  # Create crypto material using cryptogen
  if [ "$CRYPTO" == "cryptogen" ]; then
    which cryptogen
    if [ "$?" -ne 0 ]; then
      fatalln "cryptogen tool não encontrado. saindo"
    fi
    infoln "Gerando certificados usando o cryptogen tool"

    infoln "Criando Identidades de Org1"

    set -x
    cryptogen generate --config=./organizations/cryptogen/crypto-config-org1.yaml --output="organizations"
    res=$?
    { set +x; } 2>/dev/null
    if [ $res -ne 0 ]; then
      fatalln "Falha em gerar os certificados..."
    fi

    infoln "Criando certificados para Org2"

    set -x
    cryptogen generate --config=./organizations/cryptogen/crypto-config-org2.yaml --output="organizations"
    res=$?
    { set +x; } 2>/dev/null
    if [ $res -ne 0 ]; then
      fatalln "Falha em criar os certificados..."
    fi

    infoln "Criando Orderer (Pedido) Org Identidades"

    set -x
    cryptogen generate --config=./organizations/cryptogen/crypto-config-orderer.yaml --output="organizations"
    res=$?
    { set +x; } 2>/dev/null
    if [ $res -ne 0 ]; then
      fatalln "Falha em gerar os certificados..."
    fi

  fi

  # Criando crypto material usando Fabric CA
  if [ "$CRYPTO" == "Certificate Authorities" ]; then
    infoln "Generating certificates using Fabric CA"
    IMAGE_TAG=${CA_IMAGETAG} docker-compose -f $COMPOSE_FILE_CA up -d 2>&1

    . organizations/fabric-ca/registerEnroll.sh

  while :
    do
      if [ ! -f "organizations/fabric-ca/org1/tls-cert.pem" ]; then
        sleep 1
      else
        break
      fi
    done

    infoln "Creating Org1 Identities"

    createOrg1

    infoln "Creating Org2 Identities"

    createOrg2

    infoln "Creating Orderer Org Identities"

    createOrderer

  fi

  infoln "Generating CCP files for Org1 and Org2"
  ./organizations/ccp-generate.sh
}

# Once you create the organization crypto material, you need to create the
# genesis block of the application channel.

# The configtxgen tool is used to create the genesis block. Configtxgen consumes a
# "configtx.yaml" file that contains the definitions for the sample network. The
# genesis block is defined using the "TwoOrgsApplicationGenesis" profile at the bottom
# of the file. This profile defines an application channel consisting of our two Peer Orgs.
# The peer and ordering organizations are defined in the "Profiles" section at the
# top of the file. As part of each organization profile, the file points to the
# location of the MSP directory for each member. This MSP is used to create the channel
# MSP that defines the root of trust for each organization. In essence, the channel
# MSP allows the nodes and users to be recognized as network members.
#
# If you receive the following warning, it can be safely ignored:
#
# [bccsp] GetDefault -> WARN 001 Before using BCCSP, please call InitFactories(). Falling back to bootBCCSP.
#
# You can ignore the logs regarding intermediate certs, we are not using them in
# this crypto implementation.

function createConsortium() {
  which configtxgen
  if [ "$?" -ne 0 ]; then
    fatalln "configtxgen tool not found."
  fi

  infoln "Generating Orderer Genesis block"

  # Note: For some unknown reason (at least for now) the block file can't be
  # named orderer.genesis.block or the orderer will fail to launch!
  set -x
  configtxgen -profile TwoOrgsOrdererGenesis -channelID system-channel -outputBlock ./system-genesis-block/genesis.block
  res=$?
  { set +x; } 2>/dev/null
  if [ $res -ne 0 ]; then
    fatalln "Failed to generate orderer genesis block..."
  fi
}


# After we create the org crypto material and the application channel genesis block,
# we can now bring up the peers and ordering service. By default, the base
# file for creating the network is "docker-compose-test-net.yaml" in the ``docker``
# folder. This file defines the environment variables and file mounts that
# point the crypto material and genesis block that were created in earlier.

# Bring up the peer and orderer nodes using docker compose.

# Depois de criar o material criptográfico da organização, você precisa criar o
# bloco de gênese do canal de aplicação.

# A ferramenta configtxgen é usada para criar o bloco genesis. Configtxgen consome um
# Arquivo "configtx.yaml" que contém as definições para a rede de amostra. O
# bloco genesis é definido usando o perfil "TwoOrgsApplicationGenesis" na parte inferior
# do arquivo. Este perfil define um canal de aplicação que consiste em nossas duas organizações pares.
# As organizações pares e ordenadoras são definidas na seção "Perfis" no
# início do arquivo. Como parte do perfil de cada organização, o arquivo aponta para o
# localização do diretório MSP para cada membro. Este MSP é usado para criar o canal
# MSP que define a raiz de confiança de cada organização. Em essência, o canal
# MSP permite que os nós e usuários sejam reconhecidos como membros da rede.
#
# Se você receber o seguinte aviso, ele pode ser ignorado com segurança:
#
# [bccsp] GetDefault -> WARN 001 Antes de usar o BCCSP, chame InitFactories (). Retornando ao bootBCCSP.
#
# Você pode ignorar os registros sobre certificados intermediários, não os estamos usando em
# esta implementação de criptografia.

# Depois de criarmos o material criptográfico org e o bloco de geração do canal de aplicação,
# agora podemos trazer os pares e o serviço de pedidos. Por padrão, a base
# arquivo para criar a rede é "docker-compose-test-net.yaml" no `` docker``
# pasta. Este arquivo define as variáveis ​​de ambiente e montagens de arquivo que
# apontar o material criptográfico e o bloco de gênese que foram criados anteriormente.

# Abra os nós de mesmo nível e de pedido usando o docker compose.
function networkUp() {
  checkPrereqs
  # generate artifacts if they don't exist
  #gerando artefatos se eles não existem
  if [ ! -d "organizations/peerOrganizations" ]; then
    createOrgs
	createConsortium
  fi

  COMPOSE_FILES="-f ${COMPOSE_FILE_BASE}"

  if [ "${DATABASE}" == "couchdb" ]; then
    COMPOSE_FILES="${COMPOSE_FILES} -f ${COMPOSE_FILE_COUCH}"
  fi

  IMAGE_TAG=$IMAGETAG docker-compose ${COMPOSE_FILES} up -d 2>&1


  docker ps -a
  if [ $? -ne 0 ]; then
    fatalln "Incapaz de iniciar a rede"
  fi
}

# call the script to create the channel, join the peers of org1 and org2,
# and then update the anchor peers for each organization
# chame o script para criar o canal, junte-se aos pares de org1 e org2,
# e, em seguida, atualize os pares âncora para cada organização
function createChannel() {
  # Bring up the network if it is not already up.
  # Abra a rede se ela ainda não estiver ativa.

  if [ ! -d "organizations/peerOrganizations" ]; then
    infoln "Bringing up network"
    infoln "Trazendo rede"
    networkUp
  fi

  # now run the script that creates a channel. This script uses configtxgen once
  # to create the channel creation transaction and the anchor peer updates.
  # agora execute o script que cria um canal. Este script usa configtxgen uma vez
  # para criar a transação de criação de canal e as atualizações de pares âncora.
  scripts/createChannel.sh $CHANNEL_NAME $CLI_DELAY $MAX_RETRY $VERBOSE
}


## Call the script to deploy a chaincode to the channel
## Chame o script para implantar um chaincode para o canal
function deployCC() {
  scripts/deployCC.sh $CHANNEL_NAME $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $CLI_DELAY $MAX_RETRY $VERBOSE

  if [ $? -ne 0 ]; then
    fatalln "Deploying chaincode failed"
    fatalln "Implantação do chaincode falhou"
  fi
}


# Tear down running network
## Derrube a rede em execução
function networkDown() {
  # stop org3 containers also in addition to org1 and org2, in case we were running sample to add org3
  # stop org3 containers também além de org1 e org2, caso estivéssemos executando o sample para adicionar org3
   docker-compose -f $COMPOSE_FILE_BASE -f $COMPOSE_FILE_COUCH -f $COMPOSE_FILE_CA down --volumes --remove-orphans
  docker-compose -f $COMPOSE_FILE_COUCH_ORG3 -f $COMPOSE_FILE_ORG3 down --volumes --remove-orphans
  # Don't remove the generated artifacts -- note, the ledgers are always removed
  ## Não remova os artefatos gerados - observe, os livros-razão são sempre removidos
  if [ "$MODE" != "restart" ]; then
    # Bring down the network, deleting the volumes
    #Cleanup the chaincode containers
    # Derrubar a rede, excluindo os volumes
    #Limpe os contêineres do chaincode
    clearContainers
    #Cleanup images
    #Limpando imagens
    removeUnwantedImages
     # remove orderer block and other channel configuration transactions and certs
    docker run --rm -v $(pwd):/data busybox sh -c 'cd /data && rm -rf system-genesis-block/*.block organizations/peerOrganizations organizations/ordererOrganizations'
    ## remove fabric ca artifacts
    docker run --rm -v $(pwd):/data busybox sh -c 'cd /data && rm -rf organizations/fabric-ca/org1/msp organizations/fabric-ca/org1/tls-cert.pem organizations/fabric-ca/org1/ca-cert.pem organizations/fabric-ca/org1/IssuerPublicKey organizations/fabric-ca/org1/IssuerRevocationPublicKey organizations/fabric-ca/org1/fabric-ca-server.db'
    docker run --rm -v $(pwd):/data busybox sh -c 'cd /data && rm -rf organizations/fabric-ca/org2/msp organizations/fabric-ca/org2/tls-cert.pem organizations/fabric-ca/org2/ca-cert.pem organizations/fabric-ca/org2/IssuerPublicKey organizations/fabric-ca/org2/IssuerRevocationPublicKey organizations/fabric-ca/org2/fabric-ca-server.db'
    docker run --rm -v $(pwd):/data busybox sh -c 'cd /data && rm -rf organizations/fabric-ca/ordererOrg/msp organizations/fabric-ca/ordererOrg/tls-cert.pem organizations/fabric-ca/ordererOrg/ca-cert.pem organizations/fabric-ca/ordererOrg/IssuerPublicKey organizations/fabric-ca/ordererOrg/IssuerRevocationPublicKey organizations/fabric-ca/ordererOrg/fabric-ca-server.db'
    docker run --rm -v $(pwd):/data busybox sh -c 'cd /data && rm -rf addOrg3/fabric-ca/org3/msp addOrg3/fabric-ca/org3/tls-cert.pem addOrg3/fabric-ca/org3/ca-cert.pem addOrg3/fabric-ca/org3/IssuerPublicKey addOrg3/fabric-ca/org3/IssuerRevocationPublicKey addOrg3/fabric-ca/org3/fabric-ca-server.db'
    # remove channel and script artifacts
    docker run --rm -v $(pwd):/data busybox sh -c 'cd /data && rm -rf channel-artifacts log.txt *.tar.gz'
  fi
}
OS_ARCH=$(echo "$(uname -s | tr '[:upper:]' '[:lower:]' | sed 's/mingw64_nt.*/windows/')-$(uname -m | sed 's/x86_64/amd64/g')" | awk '{print tolower($0)}')
# Using crpto vs CA. default is cryptogen
# Usando cryptogen vs CA. por padrão eh cryptogen
CRYPTO="cryptogen"
# timeout duration - the duration the CLI should wait for a response from
# another container before giving up
# duração do tempo limite - a duração que a CLI deve esperar por uma resposta
# outro contêiner antes de desistir
MAX_RETRY=5
# default for delay between commands
# padrão para atraso entre os comandos
CLI_DELAY=3
# channel name defaults to "mychannel"
# padrão de nome de canal 
CHANNEL_NAME="mychannel"
# chaincode name defaults to "NA"
# chaincode com nome padrão NA
CC_NAME="NA"
# chaincode path defaults to "NA"
#localização dele eh NA (que porra eh não sei)
CC_SRC_PATH="NA"
# endorsement policy defaults to "NA". This would allow chaincodes to use the majority default policy.
# a política de endosso é padronizada como "NA". Isso permitiria que os chaincodes usassem a política padrão da maioria
CC_END_POLICY="NA"
# collection configuration defaults to "NA"
CC_COLL_CONFIG="NA"
# chaincode init function defaults to "NA"
CC_INIT_FCN="NA"
# use this as the default docker-compose yaml definition
COMPOSE_FILE_BASE=docker/docker-compose-test-net.yaml
# docker-compose.yaml file if you are using couchdb
COMPOSE_FILE_COUCH=docker/docker-compose-couch.yaml
# certificate authorities compose file
COMPOSE_FILE_CA=docker/docker-compose-ca.yaml
# use this as the docker compose couch file for org3
COMPOSE_FILE_COUCH_ORG3=addOrg3/docker/docker-compose-couch-org3.yaml
# use this as the default docker-compose yaml definition for org3
COMPOSE_FILE_ORG3=addOrg3/docker/docker-compose-org3.yaml
#
# chaincode language defaults to "NA"
CC_SRC_LANGUAGE="NA"
# Chaincode version
CC_VERSION="1.0"
# Chaincode definition sequence
CC_SEQUENCE=1
# default image tag
IMAGETAG="latest"
# default ca image tag
CA_IMAGETAG="latest"
# default database
DATABASE="leveldb"



# Parse commandline args

## Parse mode
if [[ $# -lt 1 ]] ; then
  printHelp
  exit 0
else
  MODE=$1
  shift
fi

# se foi passado criação do canal
if [[ $# -ge 1 ]] ; then
  key="$1"
  if [[ "$key" == "createChannel" ]]; then
      export MODE="createChannel"
      shift
  fi
fi

# parse flags

while [[ $# -ge 1 ]] ; do
  key="$1"
  case $key in
  -h )
    printHelp $MODE
    exit 0
    ;;
  -c )
    CHANNEL_NAME="$2"
    shift
    ;;
  -ca )
    CRYPTO="Certificate Authorities"
    ;;
  -r )
    MAX_RETRY="$2"
    shift
    ;;
  -d )
    CLI_DELAY="$2"
    shift
    ;;
  -s )
    DATABASE="$2"
    shift
    ;;
  -ccl )
    CC_SRC_LANGUAGE="$2"
    shift
    ;;
  -ccn )
    CC_NAME="$2"
    shift
    ;;
  -ccv )
    CC_VERSION="$2"
    shift
    ;;
  -ccs )
    CC_SEQUENCE="$2"
    shift
    ;;
  -ccp )
    CC_SRC_PATH="$2"
    shift
    ;;
  -ccep )
    CC_END_POLICY="$2"
    shift
    ;;
  -cccg )
    CC_COLL_CONFIG="$2"
    shift
    ;;
  -cci )
    CC_INIT_FCN="$2"
    shift
    ;;
  -i )
    IMAGETAG="$2"
    shift
    ;;
  -cai )
    CA_IMAGETAG="$2"
    shift
    ;;

  -verbose )
    VERBOSE=true
    shift
    ;;
  * )
    errorln "Unknown flag: $key"
    printHelp
    exit 1
    ;;
  esac
  shift
done

# Are we generating crypto material with this command?
if [ ! -d "organizations/peerOrganizations" ]; then
  CRYPTO_MODE="com criptografia de '${CRYPTO}'"
else
  CRYPTO_MODE=""
fi

# Determinar o modo de operação e imprimir o que solicitamos
if [ "$MODE" == "up" ]; then
  infoln "Iniciando nós com CLI timeout de '${MAX_RETRY}' tries and CLI delay of '${CLI_DELAY}' segundos e usando banco de dados '${DATABASE}' ${CRYPTO_MODE}"
elif [ "$MODE" == "createChannel" ]; then
  infoln "Criando canal '${CHANNEL_NAME}'."
  infoln "Se a rede não estiver ativa, iniciando nós com CLI timeout de '${MAX_RETRY}' tries and CLI delay of '${CLI_DELAY}' segundos e usando banco de dados '${DATABASE} ${CRYPTO_MODE}"
elif [ "$MODE" == "down" ]; then
  infoln "Parando a rede"
elif [ "$MODE" == "restart" ]; then
  infoln "Restartando a rede"
elif [ "$MODE" == "deployCC" ]; then
  infoln "fazendo deploy chaincode no canal '${CHANNEL_NAME}'"
else
  printHelp
  exit 1
fi

if [ "${MODE}" == "up" ]; then
  networkUp
elif [ "${MODE}" == "createChannel" ]; then
  createChannel
elif [ "${MODE}" == "deployCC" ]; then
  deployCC
elif [ "${MODE}" == "down" ]; then
  networkDown
else
  printHelp
  exit 1
fi
