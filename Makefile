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
	mkdir -p $(DEB_DIR)/var/www/html/wild-central
	mkdir -p $(DEB_DIR)/etc/nginx/sites-available
	
	cp $(BUILD_DIR)/$(BINARY_NAME) $(DEB_DIR)/usr/bin/
	cp wild-cloud-central.service $(DEB_DIR)/etc/systemd/system/
	cp config.yaml $(DEB_DIR)/etc/wild-cloud-central/config.yaml.example
	cp -r static/* $(DEB_DIR)/var/www/html/wild-central/
	cp wild-central-nginx.conf $(DEB_DIR)/etc/nginx/sites-available/wild-central
	
	# Create control file
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
	
	# Create postinst script for proper installation
	echo "#!/bin/bash" > $(DEB_DIR)/DEBIAN/postinst
	echo "set -e" >> $(DEB_DIR)/DEBIAN/postinst
	echo "# Create wildcloud user if it doesn't exist" >> $(DEB_DIR)/DEBIAN/postinst
	echo "if ! id wildcloud >/dev/null 2>&1; then" >> $(DEB_DIR)/DEBIAN/postinst
	echo "    useradd -r -s /bin/false wildcloud" >> $(DEB_DIR)/DEBIAN/postinst
	echo "fi" >> $(DEB_DIR)/DEBIAN/postinst
	echo "# Create required directories" >> $(DEB_DIR)/DEBIAN/postinst
	echo "mkdir -p /var/lib/wild-cloud-central /var/log/wild-cloud-central /var/www/html/talos /var/ftpd" >> $(DEB_DIR)/DEBIAN/postinst
	echo "chown wildcloud:wildcloud /var/lib/wild-cloud-central /var/log/wild-cloud-central" >> $(DEB_DIR)/DEBIAN/postinst
	echo "chown -R www-data:www-data /var/www/html" >> $(DEB_DIR)/DEBIAN/postinst
	echo "chmod 755 /var/ftpd" >> $(DEB_DIR)/DEBIAN/postinst
	echo "# Configure nginx" >> $(DEB_DIR)/DEBIAN/postinst
	echo "ln -sf /etc/nginx/sites-available/wild-central /etc/nginx/sites-enabled/" >> $(DEB_DIR)/DEBIAN/postinst
	echo "rm -f /etc/nginx/sites-enabled/default" >> $(DEB_DIR)/DEBIAN/postinst
	echo "# Enable systemd service" >> $(DEB_DIR)/DEBIAN/postinst
	echo "systemctl daemon-reload" >> $(DEB_DIR)/DEBIAN/postinst
	echo "systemctl enable wild-cloud-central" >> $(DEB_DIR)/DEBIAN/postinst
	echo "systemctl reload nginx || true" >> $(DEB_DIR)/DEBIAN/postinst
	chmod 755 $(DEB_DIR)/DEBIAN/postinst
	
	dpkg-deb --build $(DEB_DIR) $(BUILD_DIR)/wild-cloud-central_$(VERSION)_amd64.deb

dev:
	go run . &
	echo "Server started on http://localhost:8081"