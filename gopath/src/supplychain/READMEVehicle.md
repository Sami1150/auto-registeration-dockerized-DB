Demonstrates the use of history
===============================
History.go chaincode manages Cars as assets on the chain.
<!-- Each Car managed on the chain is identified by a key = VIN (Vehicle Identification Number) -->
Test data will be setup in the init() function

TransferOwnership  = Transaction to Transfer the ownership of the car
GetVehicleHistory  = History of all ownership transfer transactions can be retrieved

Setup chaincode vendoring
=========================
Switch to the chaincode folder. 
cd $GOPATH/src/token/history
./govendor.sh

If vendoring is not done the "instantiate" command will fail in net mode
This script may take upto 10+ min sometime :) PLEASE be patient
(ignore warning/error)

Setup the chaincode
===================
. set-env.sh acme
<!-- set-chain-env.sh -n history  -v 1.0   -p  AutomobileBS  -c '{"Args":["init","Pkr","10000000", "Pakistani Rupee!!!","MAdil"]}'
set-chain-env.sh -n transfer  -v 1.0   -p  AutomobileBS  -c '{"Args":["init","Pkr","10000000", "Pakistani Rupee!!!","MAdil"]}' -->
set-chain-env.sh -n supplychain  -v 1.0   -p  AutomobileBS  -c '{"Args":["init","Pkr","10000000", "Pakistani Rupee!!!","MAdil"]}'

installing chaincode
==================
 chain.sh install -p
chain.sh instantiate
Query
=====
Query the balance for 'Madil'
 set-chain-env.sh         -q   '{"Args":["balanceOf","3520299610969"]}'
 chain.sh query

Invoke
======
Transfer 100 tokens from 'MAdil' to 'Dar'
  set-chain-env.sh         -i   '{"Args":["transfer", "3520299610969", "1234", "100","Dar"]}'
  chain.sh  invoke

   set-chain-env.sh         -q   '{"Args":["balanceOf","1234"]}'
 chain.sh query

 set-chain-env.sh         -q   '{"Args":["balanceOf","6666"]}'
 chain.sh query
 
 set-chain-env.sh         -i   '{"Args":["Manufacture", "3520299610969", "23B8", "A7655","Honda","2022","Car","Civic","7777","10000","12000","2/11/2022"]}'
  chain.sh  invoke

   set-chain-env.sh         -i   '{"Args":["Manufacture", "3520299610969", "23Be8", "A76e55","Honda","2022","Car","City","7777","10000","12000","2/11/2022"]}'
  chain.sh  invoke

  set-chain-env.sh  -q '{"Args": ["GetVehiclesByCNIC", "7777"]}'
chain.sh query
Query
=====
Check the balance for 'Madil' & 'Dar'
 set-chain-env.sh         -q   '{"Args":["balanceOf","3520299610969"]}'
 chain.sh query
 set-chain-env.sh         -q   '{"Args":["balanceOf","6666"]}'
 chain.sh query

Query 
====
get the vehicle by chassis no ,engine no,company name
by chassis no
<!-- set-chain-env.sh         -q   '{"Args":["getStateRangeOnKey","23B8"]}'
 chain.sh query
 by engine no
 set-chain-env.sh         -q   '{"Args":["getStateRangeOnKey","A7655"]}'
 chain.sh query

 by company name
  set-chain-env.sh         -q   '{"Args":["getStateRangeOnKey","Honda"]}'
 chain.sh query -->
  by combining two or more
   set-chain-env.sh         -q   '{"Args":["getStateRangeOnKey","A7655","23B8","Honda"]}'
 chain.sh query




<!-- set-chain-env.sh  -q '{"Args": ["GetVehicleByVin", "100"]}'
chain.sh query -->

Transfer Ownership
 from ,to ,car composite key which is chassisNo~engineNo~companyName

chain.sh invoke
set-chain-env.sh  -i '{"Args": ["TransferOwnership", "7777","1234","A7655","23B8","Honda","2019-02-01"]}'
chain.sh invoke



get vehicle history by providing three things
chassisNO~EngineNo~Company name
set-chain-env.sh  -q '{"Args": ["GetVehicleHistory", "A7655","23B8","Honda"]}'
chain.sh query




<!-- Assets:
======
VIN,Make,Model,Year,Owner
100,toyota,corolla,2001,J Smith
200,honda,civic,199,G Roger
300,audi,a5,1999,S Ripple
400,bmw,x5,2013,M Jane
500,toyota,camry,2018,J Hoover -->

KeyModification
===============
https://godoc.org/github.com/hyperledger/fabric/protos/ledger/queryresult

Additional function
===================
set-chain-env.sh  -q '{"Args": ["GetVehiclesByYear", "2012"]}'
chain.sh query

Index JSON:
{
    "index": {
       "fields": [
          "year"
       ]
    },
    "name": "index-on-year",
    "ddoc": "index-on-year",
    "type": "json"
 }