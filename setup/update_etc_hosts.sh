#!/bin/bash
# Update /etc/hosts
source    ./manage_hosts.sh

HOSTNAME=peer1.excise.com
removehost $HOSTNAME            &> /dev/null
addhost $HOSTNAME
HOSTNAME=peer1.fbr.com
removehost $HOSTNAME            &> /dev/null
addhost $HOSTNAME
HOSTNAME=peer1.manufacturer.com
removehost $HOSTNAME            &> /dev/null
addhost $HOSTNAME
HOSTNAME=orderer.excise.com
removehost $HOSTNAME            &> /dev/null
addhost $HOSTNAME
HOSTNAME=postgresql
removehost $HOSTNAME            &> /dev/null
addhost $HOSTNAME
HOSTNAME=explorer
removehost $HOSTNAME            &> /dev/null
addhost $HOSTNAME
HOSTNAME=vagrant
removehost $HOSTNAME            &> /dev/null
addhost $HOSTNAME
HOSTNAME=excise-peer1.couchdb
removehost $HOSTNAME            &> /dev/null
addhost $HOSTNAME
HOSTNAME=fbr-peer1.couchdb
removehost $HOSTNAME            &> /dev/null
addhost $HOSTNAME
HOSTNAME=manufacturer-peer1.couchdb
removehost $HOSTNAME            &> /dev/null
addhost $HOSTNAME