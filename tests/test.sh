#!/bin/sh
set -e

if [ ! -f test.sh ]; then
    echo "Make sure to run this script from the 'tests' directory"
    exit -1
fi
echo -n "Using terraform:"
if which terraform; then
    terraform version
else
    echo "terraform not found in path, terminating"
    exit -1
fi

tests_dir=`pwd`
registry_dir=$tests_dir/fake_registry
registry_bin_dir=$registry_dir/terraform-registry.env0.com/env0/env0/6.6.6/linux_amd64

mkdir -p $registry_bin_dir
cp -a ../terraform-provider-env0 $registry_bin_dir/terraform-provider-env0

banner() {
    echo "*************************************************"
    echo "** $@"
    echo "*************************************************"
}

run_test_dir() {
    set -e
    cd $1
    rm -fr .terraform .terraform.lock.hcl terraform.state terraform.rc
    cat > terraform.rc <<EOF
provider_installation {
  filesystem_mirror {
    path    = "$registry_dir"
    include = ["terraform-registry.env0.com/*/*"]
  }
}
EOF
    banner "Running terraform init in $1"
    TF_CLI_CONFIG_FILE=terraform.rc terraform init
    terraform fmt
    banner "Terraform init success in $1"
    TF_CLI_CONFIG_FILE=terraform.rc terraform apply -auto-approve
    banner "Test success in $1"
}

for name in $@; do
    run_test_dir $name
done
