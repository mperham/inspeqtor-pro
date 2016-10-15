SHELL := /bin/bash
NAME=inspeqtor-pro
VERSION=1.0.1
# when fixing packaging bugs but not changing the binary, we increment this number
ITERATION=1

BASENAME=$(NAME)_$(VERSION)-$(ITERATION)

-include .local.sh

all: test

# gocc produces ill-formated code, clean it up with fmt
parsers: gocc fmt

# Generate parser and token package for cron jobs
gocc:
	cd $(shell pwd)/jobs && $(GOPATH)/bin/gocc format.bnf

# goimports produces slightly different formatted code from go fmt
fmt:
	find . -name "*.go" -exec goimports -w {} \;

lint:
	gometalinter ./...

prepare:
	go get github.com/stretchr/testify
	go get golang.org/x/crypto/nacl/box
	go get github.com/alecthomas/gometalinter
	go get -d github.com/mperham/inspeqtor
	go get github.com/jteeuwen/go-bindata
	gometalinter --install

license:
	@go run cmd/lic.go $(to)
	@rm -f license.bin

assets:
	@pushd ../inspeqtor >/dev/null && \
		cp ../inspeqtor-pro/templates/email/Job* templates/email && \
  	go-bindata -pkg inspeqtor -o templates.go templates/... && \
	 	rm templates/email/Job* && \
		popd >/dev/null
	@go generate github.com/mperham/inspeqtor-pro/channels
	@go generate github.com/mperham/inspeqtor-pro/expose

# rebuild Inspeqtor's email templates since we need to include the Job templates
test: assets
	@go test -parallel 4 \
	 	github.com/mperham/inspeqtor-pro \
	 	github.com/mperham/inspeqtor-pro/channels \
		github.com/mperham/inspeqtor-pro/expose \
		github.com/mperham/inspeqtor-pro/jobs \
		github.com/mperham/inspeqtor-pro/ownership \
		github.com/mperham/inspeqtor-pro/statsd \
		| grep -v "no test files"

cover:
	go test -cover -coverprofile cover.out github.com/mperham/inspeqtor-pro/expose
	go tool cover -html=cover.out

# brew update
# brew upgrade go --with-cc-common
build: test
	@GOOS=linux GOARCH=amd64 go build -o inspeqtor main.go self.go licensing.go

clean:
	rm -f main inspeqtor
	rm -rf packaging/output
	mkdir -p packaging/output/upstart
	mkdir -p packaging/output/systemd

git_update:
	cd ../inspeqtor && git pull
	cd ../inspeqtor-pro && git pull

real: assets
	GOMAXPROCS=4 go run -race main.go self.go licensing.go -l debug -s i.sock -c realtest

package: clean test build build_deb build_rpm
# 	build_rpm

purge_deb:
	ssh -t $(DEB_PRODUCTION) 'sudo apt-get purge -y $(NAME) && sudo rm -f /etc/inspeqtor' || true

purge_rpm:
	ssh -t $(RPM_PRODUCTION) 'sudo rpm -e $(NAME) && sudo rm -f /etc/inspeqtor' || true

deploy_deb: clean build_deb purge_deb
	scp packaging/output/upstart/*.deb $(DEB_PRODUCTION):~
	ssh $(DEB_PRODUCTION) 'sudo rm -f /etc/inspeqtor && sudo dpkg -i $(NAME)_$(VERSION)-$(ITERATION)_amd64.deb && sudo ./fix && sudo restart inspeqtor || true'

deploy_rpm: clean build_rpm purge_rpm
	scp packaging/output/systemd/*.rpm $(RPM_PRODUCTION):~
	ssh -t $(RPM_PRODUCTION) 'sudo rm -f /etc/inspeqtor && sudo yum install -q -y $(NAME)-$(VERSION)-$(ITERATION).x86_64.rpm && sudo ./fix && sudo systemctl restart inspeqtor'

update_deb: clean build_deb
	scp packaging/output/upstart/*.deb $(DEB_PRODUCTION):~
	ssh $(DEB_PRODUCTION) 'sudo dpkg -i $(NAME)_$(VERSION)-$(ITERATION)_amd64.deb'

update_rpm: clean build_rpm
	scp packaging/output/systemd/*.rpm $(RPM_PRODUCTION):~
	ssh -t $(RPM_PRODUCTION) 'sudo yum install -q -y $(NAME)-$(VERSION)-$(ITERATION).x86_64.rpm'

deploy: deploy_deb deploy_rpm
purge: purge_deb purge_rpm

tag:
	git tag v$(VERSION)-$(ITERATION)
	git push --tags

# gem install -N fpm package_cloud
upload: package
	package_cloud push contribsys/inspeqtor-pro/ubuntu/xenial packaging/output/systemd/$(NAME)_$(VERSION)-$(ITERATION)_amd64.deb
	package_cloud push contribsys/inspeqtor-pro/ubuntu/precise packaging/output/upstart/$(NAME)_$(VERSION)-$(ITERATION)_amd64.deb
	package_cloud push contribsys/inspeqtor-pro/ubuntu/trusty packaging/output/upstart/$(NAME)_$(VERSION)-$(ITERATION)_amd64.deb
	package_cloud push contribsys/inspeqtor-pro/el/7 packaging/output/systemd/$(NAME)-$(VERSION)-$(ITERATION).x86_64.rpm
	package_cloud push contribsys/inspeqtor-pro/el/6 packaging/output/upstart/$(NAME)-$(VERSION)-$(ITERATION).x86_64.rpm

build_rpm: build_rpm_upstart build_rpm_systemd
build_deb: build_deb_upstart build_deb_systemd

build_rpm_upstart: build
	# gem install fpm
	# brew install rpm
	fpm -s dir -t rpm -n $(NAME) -v $(VERSION) -p packaging/output/upstart \
		--rpm-compression bzip2 --rpm-os linux \
		--replaces inspeqtor \
	 	--after-install ../inspeqtor/packaging/scripts/postinst.rpm.upstart \
	 	--before-remove ../inspeqtor/packaging/scripts/prerm.rpm.upstart \
		--after-remove ../inspeqtor/packaging/scripts/postrm.rpm.upstart \
		--url http://contribsys.com/inspeqtor \
		--description "Modern service monitoring" \
		-m "Contributed Systems LLC <oss@contribsys.com>" \
		--iteration $(ITERATION) --license "Commercial" \
		--vendor "Contributed Systems" -a amd64 \
		inspeqtor=/usr/bin/inspeqtor \
		../inspeqtor/packaging/root/=/

build_rpm_systemd: build
	# gem install fpm
	# brew install rpm
	fpm -s dir -t rpm -n $(NAME) -v $(VERSION) -p packaging/output/systemd \
		--rpm-compression bzip2 --rpm-os linux \
		--replaces inspeqtor \
	 	--after-install ../inspeqtor/packaging/scripts/postinst.rpm.systemd \
	 	--before-remove ../inspeqtor/packaging/scripts/prerm.rpm.systemd \
		--after-remove ../inspeqtor/packaging/scripts/postrm.rpm.systemd \
		--url http://contribsys.com/inspeqtor \
		--description "Modern service monitoring" \
		-m "Contributed Systems LLC <oss@contribsys.com>" \
		--iteration $(ITERATION) --license "Commercial" \
		--vendor "Contributed Systems" -a amd64 \
		inspeqtor=/usr/bin/inspeqtor \
		../inspeqtor/packaging/root/=/

build_deb_upstart: build
	# gem install fpm
	fpm -s dir -t deb -n $(NAME) -v $(VERSION) -p packaging/output/upstart \
		--deb-priority optional --category admin \
		--deb-compression bzip2 \
		--replaces inspeqtor \
	 	--after-install ../inspeqtor/packaging/scripts/postinst.deb.upstart \
	 	--before-remove ../inspeqtor/packaging/scripts/prerm.deb.upstart \
		--after-remove ../inspeqtor/packaging/scripts/postrm.deb.upstart \
		--url http://contribsys.com/inspeqtor \
		--description "Modern service monitoring" \
		-m "Contributed Systems LLC <oss@contribsys.com>" \
		--iteration $(ITERATION) --license "Commercial" \
		--vendor "Contributed Systems" -a amd64 \
		inspeqtor=/usr/bin/inspeqtor \
		../inspeqtor/packaging/root/=/

build_deb_systemd: build
	# gem install fpm
	fpm -s dir -t deb -n $(NAME) -v $(VERSION) -p packaging/output/systemd \
		--deb-priority optional --category admin \
		--deb-compression bzip2 \
		--replaces inspeqtor \
	 	--after-install ../inspeqtor/packaging/scripts/postinst.deb.systemd \
	 	--before-remove ../inspeqtor/packaging/scripts/prerm.deb.systemd \
		--after-remove ../inspeqtor/packaging/scripts/postrm.deb.systemd \
		--url http://contribsys.com/inspeqtor \
		--description "Modern service monitoring" \
		-m "Contributed Systems LLC <oss@contribsys.com>" \
		--iteration $(ITERATION) --license "Commercial" \
		--vendor "Contributed Systems" -a amd64 \
		inspeqtor=/usr/bin/inspeqtor \
		../inspeqtor/packaging/root/=/


.PHONY: all clean test build package upload
