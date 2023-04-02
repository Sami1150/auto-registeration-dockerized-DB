===============================================
Manually Setup Environment Using binaries on VM - QuickWay
===============================================

cd autoDocker
./init-setup.sh

===============================================
Manually Setup Environment Using binaries on VM
===============================================
Part-1   Setup the network artefacts
====================================

# 1 Generate the crypto matrial

cd autoDocker
cd config
cryptogen generate --config=crypto-config.yaml



# 2  Generate the network artefacts
export FABRIC_CFG_PATH=$PWD

configtxgen -outputBlock  ./orderer/automobilegenesis.block -channelID ordererchannel  -profile AutomobileOrdererGenesis

configtxgen -outputCreateChannelTx  automobile-channel.tx -channelID automobilechannel  -profile AutomobileChannel

cd ..

=====================================
Part-2  Setup the excise, fbr & manufacturer peers
=====================================
# 1 Launch the environment

docker-compose -f ./config/docker-compose-base.yaml -f ./config/docker-compose.couchdb.yaml up -d

# 2 As Excise create | join | update channel
. bins/set-context.sh excise
./bins/submit-channel-create.sh
./bins/join-channel.sh
./bins/anchor-update.sh

peer channel list

# 2 As fbr join | update channel
. bins/set-context.sh fbr
./bins/join-channel.sh
./bins/anchor-update.sh

peer channel list

# 3 As manufacturer join | update channel
. bins/set-context.sh manufacturer
./bins/join-channel.sh
./bins/anchor-update.sh

peer channel list


===============================================
Shutdown all docker containers
===============================================
cd autoDocker 
./shutdown.sh

===============================================
To relaunch docker containers with existing crypto material
===============================================
cd autoDocker
./launch.sh

===============================================
To delete all containers but not the crypto material
===============================================
cd autoDocker
./clean.sh

===============================================
To delete all containers along with the crypto material
===============================================
cd autoDocker
./clean.sh all

===============================================
To delete ledger volume
===============================================
cd autoDocker
./clean-volume.sh #VOLUME_NAME






