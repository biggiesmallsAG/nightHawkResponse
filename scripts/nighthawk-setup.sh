#!/bin/bash
#
# nighthawk Response setup script
# author: roshan maskey <roshanmaskey@gmail.com>
#
#


# Application version specifications
ESVER="5.6.5"
GOVER="1.9.2"
NODEVER="7.9.0"


SCRIPT_DIR=`pwd`
INSTALL_DIR="/tmp/nighthawk"
NHROOT=/opt/nighthawk



__check_result() {
        if [ $1 -eq 0 ]; then
                echo "[ OK ]"
        else
                echo "[ FAILED ]"
        fi
}


__install_ubuntu_package() {
	echo -n "Installing $1"
	apt-get -y install "$1" > /dev/null 2>&1

	if [ $? -eq 0 ]; then
		echo "[ OK ]"
	else 
		echo "[ FAILED ]"
	fi
}

__install_golang_package() {
	echo -n "Installing golang package $1 "
	# check if package is already downloaded

	go get "$1" > /dev/null 2>&1

	if [ $? -eq 0 ]; then
		echo "[ OK ]"
	else
		echo "[ FAILED ]"
	fi
}


install_ubuntu_dependencies() {
	__install_ubuntu_package vim
	__install_ubuntu_package tree
	__install_ubuntu_package tcpdump
	__install_ubuntu_package wget
	__install_ubuntu_package curl
	__install_ubuntu_package net-tools
	__install_ubuntu_package openjdk-8-jre
	__install_ubuntu_package git
	__install_ubuntu_package nginx
	__install_ubuntu_package rabbitmq-server
	__install_ubuntu_package gcc
	__install_ubuntu_package g++
	__install_ubuntu_package sqlite3
}


install_golang() {
	# check current golang version
	if [ -f /usr/local/go/bin/go ]; then
		/usr/local/go/bin/go version | grep ${GOVER} > /dev/null 2>&1
		if [ $? -eq 0 ]; then
			echo "Installing golang ${GOVER} [ OK* ]"
			source /etc/profile.d/golang.sh
			return
		fi
	fi

	# remove existing golang pacakge
	# and install new version as specified $GOVER
	rm -rf /usr/local/go

	cd ${INSTALL_DIR}
	wget https://redirector.gvt1.com/edgedl/go/go${GOVER}.linux-amd64.tar.gz > /dev/null 2>&1 
	tar zxf go${GOVER}.linux-amd64.tar.gz
	mv go /usr/local
	cd ${SCRIPT_DIR}

	# Setting profile for golang
	golang="/etc/profile.d/golang.sh"

	if [ -f $golang ]; then
		sudo rm -f $golang
	fi

	echo "#!/bin/bash" >> $golang
	echo "#" >> $golang
	echo "# Golang profile - `date -u +%Y-%m-%dT%H:%M:%SZ`" >> $golang
	echo "export GOROOT=/usr/local/go" >> $golang
	echo "export PATH=\$PATH:\$GOROOT/bin" >> $golang

	source /etc/profile.d/golang.sh
}

install_elasticsearch() {
	dpkg --list | grep elasticsearch | grep ${ESVER} > /dev/null 2>&1
	if [ $? -eq 0 ]; then
		echo "Installing elasticsearch [ OK* ]"
		return
	fi

	# remove previous installation and re-install new package
	sudo dpkg --purge elasticsearch > /dev/null 2>&1

	echo -n "Installing elasticsearch "
	cd ${INSTALL_DIR}
	wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-${ESVER}.deb > /dev/null 2>&1
	dpkg --install elasticsearch-${ESVER}.deb > /dev/null 2>&1

	if [ $? -eq 0 ]; then
		echo "[ OK ]"
	else
		echo "[ FAILED ]"
	fi


	# Update elasticsearch configuration file 
	cd /etc/default
	if [ -f elasticsearch ]; then
		mv elasticsearch elasticsearch.bak
		sed 's/#DATA_DIR.*/DATA_DIR=\/opt\/nighthawk\/data\/elasticsearch/g' elasticsearch.bak > elasticsearch
	fi

	cd ${SCRIPT_DIR}
}

install_kibana() {
	dpkg --list | grep kibana | grep ${ESVER} > /dev/null 2>&1
	if [ $? -eq 0 ]; then
		echo "Installing kibana [ OK* ]"
		return
	fi

	dpkg --purge kibana > /dev/null 2>&1

	echo -n "Installing kibana"
	cd ${INSTALL_DIR}
	wget https://artifacts.elastic.co/downloads/kibana/kibana-${ESVER}-amd64.deb > /dev/null 2>&1
	dpkg --install kibana-${ESVER}-amd64.deb
	if [ $? -eq 0 ]; then
		echo "[ OK ]"
	else
		echo "[ FAILED ]"
	fi

	cd ${SCRIPT_DIR}
}


install_golang_dependencies() {
	__install_golang_package "github.com/gorilla/mux"
	__install_golang_package "github.com/gorilla/websocket"
	__install_golang_package "gopkg.in/olivere/elastic.v5"
	__install_golang_package "github.com/streadway/amqp"
	__install_golang_package "github.com/shirou/gopsutil/cpu"
	__install_golang_package "github.com/shirou/gopsutil/mem"
	__install_golang_package "github.com/shirou/gopsutil/disk"
	__install_golang_package "gopkg.in/yaml.v2"
	__install_golang_package "github.com/mattn/go-sqlite3"
}

install_node() {
	echo "Installing node"

	# delete old node directory
	if [ -d $NHROOT/node ]; then
		rm -rf $NHROOT/node
	fi

	cd ${INSTALL_DIR}
	wget https://nodejs.org/dist/v${NODEVER}/node-v${NODEVER}-linux-x64.tar.gz > /dev/null 2>&1
	tar zxf node-v${NODEVER}-linux-x64.tar.gz
	mv node-v${NODEVER}-linux-x64 $NHROOT/node

	export PATH=$PATH:$NHROOT/node/bin
	cd ${SCRIPT_DIR}
}





#==================================================
# nighthawk response setttings
#==================================================

__create_directory() {
	if [ ! -d "$1" ]; then
		mkdir -p "$1"
	fi
}
create_nighthawk_directories() {

	# creating directories
	__create_directory $NHROOT
	__create_directory $NHROOT/bin
	__create_directory $NHROOT/etc
	__create_directory $NHROOT/etc/pki/cert		# store https server certificate here
	__create_directory $NHROOT/data
	__create_directory $NHROOT/data/elasticsearch
	__create_directory $NHROOT/data/triage		# Store json file as log per computername
	__create_directory $NHROOT/data/db 
	__create_directory $NHROOT/var
	__create_directory $NHROOT/media		# directory to store uploaded triage files 
	__create_directory $NHROOT/workspace		# temporary storage for processing data
	__create_directory $NHROOT/var/log		# store logs here
	__create_directory $NHROOT/wwwroot		# webroot

	# setting link for backward compatibility
	ln -s $NHROOT/media $NHROOT/var > /dev/null 2>&1
	ln -s $NHROOT/workspace $NHROOT/var > /dev/null 2>&1
	ln -s $NHROOT/data/db $NHROOT/var/db > /dev/null 2>&1
}


create_nighthawk_user() {
	echo "Creating nighthawk user and group"
	cat /etc/group | grep nighthawk > /dev/null 2>&1
	if [ $? -ne 0 ]; then
		groupadd --gid 3723 nighthawk 
	fi

	cat /etc/passwd | grep nighthawk > /dev/null 2>&1
	if [ $? -ne 0 ]; then
		useradd --uid 3723 --home-dir $NHROOT/var --gid 3723 --shell /bin/false --comment "nighthawk response user" nighthawk
	fi 
}

install_nighthawk_profile() {
	echo "Creating nighthawk profile"

	profile="/etc/profile.d/nighthawk.sh"

	if [ -f $profile ]; then
		rm -f $profile
	fi

	echo "#!/bin/bash" >> $profile
	echo "#" >> $profile
	echo "# nighthawk response profile - `date -u +%Y-%m-%dT%H:%M:%SZ`" >> $profile
	echo "export NHROOT=/opt/nighthawk" >> $profile
	echo "export PATH=\$PATH:\$NHROOT/bin:\$NHROOT/node/bin" >> $profile

	source /etc/profile.d/nighthawk.sh
}


create_nighthawk_binaries() {
	echo "Creating nighthawk binaries"
	cd ..
	export GOPATH=$HOME/go:`pwd`
	/usr/local/go/bin/go build -ldflags '-s' -o $NHROOT/bin/nighthawk nighthawk.go
	/usr/local/go/bin/go build -ldflags '-s' -o $NHROOT/bin/nhapi nhapi.go
	/usr/local/go/bin/go build -ldflags '-s' -o $NHROOT/bin/nhlogger nhlogger.go
	cd ${SCRIPT_DIR}
}

install_nighthawk_config() {
	echo "Installing nighthawk configuration files"
	cp -f ${SCRIPT_DIR}/../config/*.json $NHROOT/etc
}


install_nighthawk_certificate() {
	echo -n "Installing self-signed server certificate "
	/usr/bin/openssl req -x509 -newkey rsa:2048 -keyout /opt/nighthawk/etc/pki/cert/server.key -nodes -days 720 -out /opt/nighthawk/etc/pki/cert/server.crt -subj "/C=AU/L=VIC/O=nighthawk response/CN=nighthawkresponse.local" > /dev/null 2>&1
	__check_result $?
}

install_nighthawk_web() {
	echo -n "Installing nighthawk web "

	# delete existing wwwroot files
	if [ -d $NHROOT/wwwroot ]; then
		rm -rf $NHROOT/wwwroot/*
	fi

	cp -r ${SCRIPT_DIR}/../nightHawkV4/* $NHROOT/wwwroot

	cd $NHROOT/wwwroot
        npm install materialize-css@^0.100.1 > /dev/null 2>&1
        npm install npm install typescript@'>=2.1.0 <2.4.0' > /dev/null 2>&1
        npm install > /dev/null 2>&1
	if [ $? -eq 0 ]; then
		echo "[ OK ]"
	else
		echo "[ FAILED ]"
	fi
	cd ${SCRIPT_DIR}
}


configure_nighthawk_index() {
	
	# Check if elasticsearch is already running using port
	netstat -an | grep :9200 > /dev/null 2>&1
	if [ $? -ne 0 ]; then
		echo "Starting elasticsearch" 
		systemctl start elasticsearch
		sleep 5
	fi

	cd ${SCRIPT_DIR}

	echo -n "Configuring elasticsearch investigations index "
	curl -s localhost:9200/investigation1 | grep investigations > /dev/null 2>&1
	if [ $? -eq 0 ]; then
		echo "[ OK* ]"
	else 
		curl -XPUT localhost:9200/investigation1 -H "application/json" -d "@ElasticMapping.json" > /dev/null 2>&1
		__check_result $?
	fi 

	# check nighthawk index
	echo -n "Configuring elasticsearch nighthawk index"
	curl -s localhost:9200/_cat/indices?v | grep nighthawk > /dev/null 2>&1
	if [ $? -eq 0 ]; then
		echo "[ OK* ]"
	else
		curl -XPUT localhost:9200/nighthawk -H "application/json" -d "@nighthawkdb.json" > /dev/null 2>&1
		__check_result $?
	fi

}



clean_up() {
	sudo rm -rf ${INSTALL_DIR}
}


configure_nighthawk_permission() {
	# set /opt/nighthawk user and group
	chown -R nighthawk:nighthawk $NHROOT
	chmod -R 775 $NHROOT 

	# adding local users to nighthawk group
	for u in `ls /home`; do usermod --append -G nighthawk $u; done

	# change ownership to elasticsearch
	chown -R elasticsearch:elasticsearch $NHROOT/data/elasticsearch
}

#================================
#  start_main
#================================

if [ `id -u` -ne 0 ];then
	echo "Installation script must be run as root"
	exit 1
fi

# create installation temporary directory
if [ -d ${INSTALL_DIR} ]; then
	rm -rf ${INSTALL_DIR}
fi
mkdir -p ${INSTALL_DIR}


# install OS dependencies
install_ubuntu_dependencies
install_golang
install_golang_dependencies

# nighthawk setup
create_nighthawk_directories
create_nighthawk_user

# install application dependencies
install_elasticsearch
install_kibana

# install and configure nighthawk application
install_nighthawk_profile
create_nighthawk_binaries
install_node
install_nighthawk_web
install_nighthawk_config
install_nighthawk_certificate

configure_nighthawk_index
configure_nighthawk_permission
