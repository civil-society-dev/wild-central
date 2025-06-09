BINARY_NAME=wild-cloud-central
VERSION?=0.1.0
BUILD_DIR=build
DEB_DIR=debian-package

.PHONY: build clean test run install deb

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

clean:
	rm -rf $(BUILD_DIR) $(DEB_DIR)
	go clean

test:
	go test ./...

run:
	go run .

install: build
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/bin/
	sudo cp wild-cloud-central.service /etc/systemd/system/
	sudo mkdir -p /etc/wild-cloud-central
	sudo cp config.yaml /etc/wild-cloud-central/config.yaml.example
	sudo systemctl daemon-reload

deb: build
	mkdir -p $(DEB_DIR)/DEBIAN
	mkdir -p $(DEB_DIR)/usr/bin
	mkdir -p $(DEB_DIR)/etc/systemd/system
	mkdir -p $(DEB_DIR)/etc/wild-cloud-central
	
	cp $(BUILD_DIR)/$(BINARY_NAME) $(DEB_DIR)/usr/bin/
	cp wild-cloud-central.service $(DEB_DIR)/etc/systemd/system/
	cp config.yaml $(DEB_DIR)/etc/wild-cloud-central/config.yaml.example
	
	echo "Package: wild-cloud-central" > $(DEB_DIR)/DEBIAN/control
	echo "Version: $(VERSION)" >> $(DEB_DIR)/DEBIAN/control
	echo "Section: net" >> $(DEB_DIR)/DEBIAN/control
	echo "Priority: optional" >> $(DEB_DIR)/DEBIAN/control
	echo "Architecture: amd64" >> $(DEB_DIR)/DEBIAN/control
	echo "Depends: dnsmasq, nginx" >> $(DEB_DIR)/DEBIAN/control
	echo "Maintainer: Wild Cloud Team <admin@wildcloud.local>" >> $(DEB_DIR)/DEBIAN/control
	echo "Description: Wild Cloud Central Management Service" >> $(DEB_DIR)/DEBIAN/control
	echo " A web-based management service for wild-cloud infrastructure" >> $(DEB_DIR)/DEBIAN/control
	echo " providing DNS, DHCP, and PXE boot services configuration." >> $(DEB_DIR)/DEBIAN/control
	
	dpkg-deb --build $(DEB_DIR) $(BUILD_DIR)/wild-cloud-central_$(VERSION)_amd64.deb

dev:
	go run . &
	echo "Server started on http://localhost:8081"