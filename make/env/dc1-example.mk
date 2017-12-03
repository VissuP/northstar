include $(MAKEPATH)/env/common.mk
include $(MAKEPATH)/env/dc1.mk

# Object
OBJECT_BLOB_STORAGE_USER_ID = example
OBJECT_BLOB_STORAGE_USER_SECRET = example

# Kafka
KAFKA_BROKER_PORT = 31816

# Cassandra
CASSANDRA_SSL_STORAGE_PORT = 31817
CASSANDRA_STORAGE_PORT = 31818
CASSANDRA_API_PORT = 31819
CASSANDRA_JMX_PORT = 31820
CASSANDRA_RPC_PORT = 31821
CASSANDRA_NATIVE_TRANSPORT_PORT = 31822
CASSANDRA_EXECUTOR_API_PORT = 31823

# TS settings
TS_ENV = staging

#Dakota -- Note that these ports change on environment.
UTPROVISION_PORT = 11292
DKT_ENV = staging

# Loggging
ENABLE_DEBUG = true
DKT_LOGGER_IS_KAFKA_ENABLED = false
DKT_LOGGER_DUMP_MSG_STDOUT = true

# Memory settings
NORTHSTARAPI_EXECUTION_MEMORY_DEFAULT = 100

# Modules
ENABLE_HTTP = true
ENABLE_NSQL = true
ENABLE_NSFTP = true
ENABLE_NSSFTP = true
ENABLE_NSOBJECT = true
ENABLE_NSSTREAM = true
ENABLE_NSKV = true

#NSSIM_TEST_GROUPS
NSSIM_TEST_GROUPS = "notebook_crud, notebook_execution, keyvalue, object, template, transformation, nsql_native, nsql_spark,  generic_execution"

# EKK
USE_EKK_STACK = true

# Spark
USE_DPE_SPARK=true
