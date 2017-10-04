#!/bin/bash

#
# nightHawk Response - Setup
#


# Setup log 
logfile="nighthawk_setup.log"
dist=""

# Global variables 
NHROOT=/opt/nighthawk
NHBIN=$NHROOT/bin
NHCONF=$NHROOT/etc
NHLIB=$NHROOT/lib
NHVAR=$NHROOT/var
NHTMP=$NHVAR/tmp 
NHLOG=$NHVAR/log




__get_os_version() {
	if [ -f /etc/lsb-release ]; then
		dist=`cat /etc/lsb-release | grep DISTRIB_ID | cut -d = -f2 | tr '[:upper:]' '[:lower:]'`
	fi

	if [ -f /etc/redhat-release ]; then
		dist="redhat"
	fi
}

create_workspace() {
	echo -n "[+] Creating folder structure " 
	mkdir -p $NHROOT >> $logfile 2>&1 
	mkdir -p $NHBIN  >> $logfile 2>&1
	mkdir -p $NHCONF >> $logfile 2>&1
	mkdir -p $NHLIB >> $logfile 2>&1
	mkdir -p $NHVAR >> $logfile 2>&1
	mkdir -p $NHTMP >> $logfile 2>&1
	mkdir -p $NHLOG >> $logfile 2>&1
	echo "[DONE]"
}


create_user() {
	echo -n "[+] Creating nighthakw user and group "
	groupadd --gid 3728 nighthawk >> $logfile 2>&1
	useradd --home-dir $NHVAR --gid 3728 --uid 3728 --shell /bin/false --comment "nightHawk Response" nighthawk >> $logfile 2>&1

	chown -R nighthawk:nighthawk $NHROOT >> $logfile 2>&1
	echo "[DONE]"
}


__install_elasticsearch() {
	echo -n "[+] Checking elasicsearch package "
	if [ "$dist" == "ubuntu" ]; then
		dpkg --status elasticsearch > /dev/null 2>&1
		if [ $? -eq 0 ]; then
			echo "[INSTALLED]"
		else 
			echo "[NOT INSTALLED]"
			
			echo -n "[+] Installing elasticsearch package "
			wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-5.3.1.deb > /dev/null 2>&1
			dpgk --install elasticsearch-5.3.1.deb > /dev/null 2>&1
			if [ $? -ne 0 ]; then 
				sudo apt-get install -f > /dev/null 2>&1 
				dpkg --install elasticsearch-5.3.1.deb > /dev/null 2>&1

				if [ $? -eq 0 ]; then
					echo "[DONE]"
				else
					echo "[FAILED]"
				fi
			fi
		fi
	fi
}

__install_kibana() {
	echo -n "[+] Checking kibana package " 
	if [ "$dist" == "ubuntu" ]; then
		dpkg --status kibana > /dev/null 2>&1
		if [ $? -eq 0 ]; then 
			echo "[INSTALLED]"
			return
		fi 

		echo "[NOT INSTALLED]"
		
		echo -n "[+] Installing kibana package "
		wget https://artifacts.elastic.co/downloads/kibana/kibana-5.3.1-amd64.deb > /dev/null 2>&1

		dpkg --install kibana-5.3.1-amd64.deb > /dev/null 2>&1
		if [ $? -ne 1 ]; then
			echo "[DONE]"
			return
		fi

		sudo apt-get install -y -f > /dev/null 2>&1
		dpkg --install kibana-5.3.1-amd64.deb > /dev/null 2>&1
		if [ $? -eq 0 ]; then
			echo "[DONE]"
		else
			echo "[FAILED]"
		fi
		
	fi

}



__ubuntu_installer() {
	echo -n "[+] Installing $1 package "
	apt-get install -y $1 > /dev/null 2>&1
	if [ $? -eq 0 ]; then
		echo "[DONE]"
	else
		echo "[FAILED]"
	fi
}

__npm_installer() {
	echo -n "[+] Installing npm $1 "
	sudo npm install -g $1 > /dev/null 2>&1 
	if [ $? -eq 0 ]; then
		echo "[DONE]"
	else
		echo "[FAILED]"
	fi
}

install_ubuntu_packages() {
	# Supporting Pre-requisite Packages
	__ubuntu_installer git 
	__ubuntu_installer vim 


	# nightHawk Response Supporting Packages 
	__ubuntu_installer rabbitmq-server 
	__ubuntu_installer amqp-tools 
	__ubuntu_installer python-amqp
	__ubuntu_installer nginx
	__ubuntu_installer npm  

	# Installing node packages
	__npm_installer @angular/cli
}

install_packages() {
	__install_elasticsearch
	__install_kibana
	install_ubuntu_packages
}



if [ `id -u` -ne 0 ]; then
	echo "This script must be run as root"
	exit 1
fi

__get_os_version
create_workspace
create_user 
install_packages