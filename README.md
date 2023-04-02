<!-- first time  -->
Clone project
Run vagrant
cd setup
./init-vexpress.sh

Now log out and log back in.

<!-- Run Project  -->
cd autoDocker
./init-setup.sh

. bins/set-context.sh excise

<!-- setting the admin context and now install the chaincode -->
Follow Steps in ChaincodeInstallingSteps file

<!-- 2nd time  -->

cd autoDocker 
./launch.sh

<!-- to remove all containers and crypto material -->

./clean.sh all


