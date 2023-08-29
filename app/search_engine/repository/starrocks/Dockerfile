FROM centos:centos7

ENV StarRocks_version=3.1.2
ENV StarRocks_home=/data/deploy

# Replace the "<url_to_download_specific_ver_of_starrocks>" below with the download path of the StarRocks that you want to deploy.
ENV StarRocks_url=https://releases.mirrorship.cn/starrocks/StarRocks-3.1.2.tar.gz

# Install StarRocks.
RUN yum -y install wget
RUN mkdir -p $StarRocks_home
RUN wget -SO $StarRocks_home/StarRocks-${StarRocks_version}.tar.gz  $StarRocks_url
RUN cd $StarRocks_home && tar zxf StarRocks-${StarRocks_version}.tar.gz

# Install Java JDK.
RUN yum -y install java-1.8.0-openjdk-devel.x86_64
RUN rpm -ql java-1.8.0-openjdk-devel.x86_64 | grep bin$

# Create directory for FE meta and BE storage in StarRocks.
RUN mkdir -p $StarRocks_home/StarRocks-${StarRocks_version}/fe/meta
RUN mkdir -p $StarRocks_home/StarRocks-${StarRocks_version}/be/storage

# Install relevant tools.
RUN yum -y install mysql net-tools telnet

# Run Setup script.
COPY run_script.sh $StarRocks_home/run_script.sh
RUN chmod +x $StarRocks_home/run_script.sh
CMD $StarRocks_home/run_script.sh