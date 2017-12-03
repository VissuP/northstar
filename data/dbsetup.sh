#!/bin/bash

retryCount=10
host=`echo $CASSANDRA_HOST | awk -F ',' '{print $1}'`
echo "Printing host $host and host list $CASSANDRA_PORT"
port=`echo $CASSANDRA_NATIVE_TRANSPORT_PORT`

echo "Running db setup on $host $port"

if [ "$CASSANDRA_AUTH_ENABLED" == "true" ]; then

        TEMP_TOKEN_URL=http://$GATEKEEPER_URL/token
        PERM_TOKEN_URL=http://$VAULT_URL/v1/sys/wrapping/unwrap
        CASSANDRA_SECRET_URL=http://$VAULT_URL/v1/secret/cassandra

        while [ "$TEMP_TOKEN" == "" ] || [ "$TEMP_TOKEN" == null ]; do
                echo "Fetching temporary token from Gatekeeper.. make sure the gatekeeper is up and running"
                TEMP_TOKEN=$(echo `curl -XPOST -H 'Content-type: application/json' $TEMP_TOKEN_URL -d '{"task_id": "'"$MESOS_TASK_ID"'"}'` | jq '.token')
                TEMP_TOKEN=`echo "${TEMP_TOKEN//\"}"`
                sleep 10
        done
        while [ "$PERM_TOKEN" == "" ] || [ "$PERM_TOKEN" == null ]; do
                echo "Fetching permanent token from Vault.. make sure the vault is up and running"
                PERM_TOKEN=$(echo `curl -XPOST -H "X-Vault-Token: $TEMP_TOKEN" $PERM_TOKEN_URL` | jq '.auth.client_token')
                PERM_TOKEN=`echo "${PERM_TOKEN//\"}"`
                sleep 10
        done
        while [ "$CASSANDRA_USERNAME" == "" ] || [ "$CASSANDRA_PASSWORD" == "" ] || [ "$CASSANDRA_USERNAME" == null ] || [ "$CASSANDRA_PASSWORD" == null ]; do
                echo "Fetching cassandra credentials from Vault.. make sure the appropriate secret has been created in vault"
        	CASSANDRA_USERNAME=$(echo `curl -XGET -H "X-Vault-Token: $PERM_TOKEN" $CASSANDRA_SECRET_URL` | jq '.data.username')
        	CASSANDRA_USERNAME=`echo "${CASSANDRA_USERNAME//\"}"`
       		CASSANDRA_PASSWORD=$(echo `curl -XGET -H "X-Vault-Token: $PERM_TOKEN" $CASSANDRA_SECRET_URL` | jq '.data.password')
        	CASSANDRA_PASSWORD=`echo "${CASSANDRA_PASSWORD//\"}"`
		sleep 10
	done

        export CASSANDRA_USERNAME=$CASSANDRA_USERNAME
        export CASSANDRA_PASSWORD=$CASSANDRA_PASSWORD

fi

cass_dir=/opt/cassandra
schema_dir=/opt/schemas

function waitForCassandraConnectivity() {
  if [ "$CASSANDRA_AUTH_ENABLED" == "true" ]; then
    until $cass_dir/bin/cqlsh --cqlversion $CASSANDRA_CQL_VERSION -u $CASSANDRA_USERNAME -p $CASSANDRA_PASSWORD $host $port  -e "describe keyspaces"
    do
        echo "Cassandra has not opened client port yet, will retry...."
        sleep 1
    done
  else
    until $cass_dir/bin/cqlsh --cqlversion $CASSANDRA_CQL_VERSION $host $port  -e "describe keyspaces"
    do
        echo "Cassandra has not opened client port yet, will retry...."
        sleep 1
    done
  fi

    echo "Cassandra client port opened now.."
}

function executeFileCommands() {
    count=$retryCount
    if [ "$CASSANDRA_AUTH_ENABLED" == "true" ]; then
    until [[ $count == 0 ]] || $cass_dir/bin/cqlsh --cqlversion $CASSANDRA_CQL_VERSION -u $CASSANDRA_USERNAME -p $CASSANDRA_PASSWORD $host $port -f $1
    do
        echo "Execution of $1 failed, will retry...."
        sleep 1
        count=$[$count - 1]
    done
  else
    until [[ $count == 0 ]] || $cass_dir/bin/cqlsh --cqlversion $CASSANDRA_CQL_VERSION $host $port  -f $1
    do
        echo "Execution of $1 failed, will retry...."
        sleep 1
        count=$[$count - 1]
    done
  fi

    if [[ $count == 0 ]] ; then
        echo "Tried execution of $1 $retryCount times. Moving ahead..."
    else
        echo "Execution of $1 success"
    fi
}

waitForCassandraConnectivity
executeFileCommands "$schema_dir/account.cql"
executeFileCommands "$schema_dir/account_upgrade.cql"
executeFileCommands "$schema_dir/cron.cql"
executeFileCommands "$schema_dir/cron_upgrade.cql"
executeFileCommands "$schema_dir/stream.cql"
executeFileCommands "$schema_dir/nssim.cql"
executeFileCommands "$schema_dir/templates_data.cql"

echo "Db setup is complete"
