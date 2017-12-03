# MODULE = Package
MODULE=$*
if [ -z $MODULE ]
  then
    # bamboo variable is assigned to MODULE
    MODULE=${bamboo_Module}
  else
    # Explicit MODULE is assigned
    MODULE=$MODULE
fi

# BRANCH = package branch
BRANCH=${bamboo_planRepository_branch}

# Using project tool
echo "Branching from this branch --- "
echo $BRANCH
bash project -b $BRANCH clone $MODULE

# pulling 3rd party dependencies
 #./project pull vendor-library

 # northstar-logger:LOGGER_TAG
 export LOGGER_TAG=0.3-1.2

 # Building individual module
 if [ -z $MODULE ]
    then
        make build | tee -a report.log
    else
        make -C go/src/github.com/verizonlabs/northstar/$MODULE build  | tee -a report.log

  fi
