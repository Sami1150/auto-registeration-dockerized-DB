#!/bin/bash
#Cleans up the specified volume

# $1 = volume to be deleted

# AutomobileRegisterationSystem_data-orderer.excise.com
# AutomobileRegisterationSystem_data-peer1.excise.com
# AutomobileRegisterationSystem_data-peer1.manufacturer.com
# AutomobileRegisterationSystem_data-peer1.fbr.com


# use command to verify if volume is deleted
# docker volume list

docker volume rm  $1

echo "Done."
