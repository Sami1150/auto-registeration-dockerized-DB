#!/bin/bash
#Check if the TLS is enabled
TLS_PARAMETERS=""
if [ "$CORE_PEER_TLS_ENABLED" == "true" ]; then
   echo "*** Executing with TLS Enabled ***"
   TLS_PARAMETERS=" --tls true --cafile $ORDERER_CA_ROOTFILE"
fi
#echo $ORDERER_ADDRESS
peer channel create -c automobilechannel -f ./config/automobile-channel.tx --outputBlock ./config/automobilechannel.block -o localhost:7050