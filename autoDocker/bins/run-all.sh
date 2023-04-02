#!/bin/bash

############################
# Setup anchor peer for excise
# Set the context
# $1 = tls in case TLS need to be enabled
. bins/set-context.sh excise $1


# Submit the channel create tx
bins/submit-channel-create.sh

# Give time for the channel tx to propagate
sleep 3s

# Join acme peer to the channel
bins/join-channel.sh

sleep 3s

# Update anchor peer in channel
bins/anchor-update.sh

############################
# Setup anchor peer for fbr
# Set the context
# $1 = tls in case TLS need to be enabled
. bins/set-context.sh fbr $1

# Join the fbr peer
bins/join-channel.sh

sleep 3s

# Update anchor peer in channel
bins/anchor-update.sh

############################
# Setup anchor peer for manufacturer
# Set the context
# $1 = tls in case TLS need to be enabled
. bins/set-context.sh manufacturer $1

# Join the manufacturer peer
bins/join-channel.sh

sleep 3s

# Update anchor peer in channel
bins/anchor-update.sh